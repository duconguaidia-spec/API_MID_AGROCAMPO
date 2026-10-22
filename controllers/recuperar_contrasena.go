package controllers

import (
	"API_MID_AGROCAMPO/models"
	"API_MID_AGROCAMPO/services"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// Minutos que dura vigente el código enviado al correo
const minutosExpiracionCodigo = 15

type Recuperar_contrasenaController struct {
	beego.Controller
}

// Solicitar genera el código de 6 dígitos, guarda su hash en el CRUD y lo envía por correo
func (c *Recuperar_contrasenaController) Solicitar() {
	baseURL := beego.AppConfig.DefaultString("url_crud_usuarios", "")

	// 1. Leer y validar el correo recibido
	var solicitud models.SolicitudRecuperacion
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err != nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "El cuerpo de la solicitud no tiene un formato válido")
		return
	}
	correo := strings.ToLower(strings.TrimSpace(solicitud.Correo))
	if _, err := mail.ParseAddress(correo); err != nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "El correo no tiene un formato válido")
		return
	}

	// 2. Buscar el usuario en el CRUD
	usuario, err := buscarUsuarioPorCorreo(baseURL, correo)
	if err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible consultar el usuario: "+err.Error())
		return
	}

	// Si no existe (o está inactivo) se responde igual, para no revelar qué correos están registrados
	activo, hayActivo := usuario["activo"].(bool)
	if usuario == nil || (hayActivo && !activo) {
		c.responderOKRecuperacion("Si el correo está registrado, recibirás un código de recuperación")
		return
	}
	idUsuario, ok := idComoEntero(usuario["id"])
	if !ok {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "El CRUD devolvió un usuario sin id")
		return
	}
	nombre, _ := usuario["nombre"].(string)

	// 3. Generar el código y guardar solo su hash en TokenRecuperacion
	codigo, err := generarCodigo()
	if err != nil {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "No fue posible generar el código")
		return
	}
	token := map[string]interface{}{
		"token_hash": hashearCodigo(codigo),
		"usuario":    map[string]interface{}{"id": idUsuario},
		"expira_en":  time.Now().Add(minutosExpiracionCodigo * time.Minute).Format(time.RFC3339),
	}
	if _, err := llamarCRUD(http.MethodPost, baseURL+"/TokenRecuperacion", token); err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible guardar el código: "+err.Error())
		return
	}

	// 4. Enviar el código por Gmail (el asunto se codifica para que las tildes se vean bien)
	servicioGmail, err := services.NewGmailService()
	if err != nil {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "No fue posible iniciar el servicio de Gmail: "+err.Error())
		return
	}
	asunto := mime.QEncoding.Encode("UTF-8", "Código de recuperación de contraseña")
	cuerpo := fmt.Sprintf(
		"Hola %s,\n\nTu código de recuperación es: %s\n\nExpira en %d minutos. Si no lo solicitaste, ignora este correo.",
		nombre, codigo, minutosExpiracionCodigo,
	)
	if _, err := servicioGmail.SendEmail(correo, asunto, cuerpo); err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible enviar el correo: "+err.Error())
		return
	}

	c.responderOKRecuperacion("Si el correo está registrado, recibirás un código de recuperación")
}

// Restablecer valida el código recibido y guarda la nueva contraseña en el CRUD
func (c *Recuperar_contrasenaController) Restablecer() {
	baseURL := beego.AppConfig.DefaultString("url_crud_usuarios", "")

	// 1. Leer y validar los datos recibidos
	var solicitud models.SolicitudRestablecer
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud); err != nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "El cuerpo de la solicitud no tiene un formato válido")
		return
	}
	correo := strings.ToLower(strings.TrimSpace(solicitud.Correo))
	codigo := strings.TrimSpace(solicitud.Codigo)
	if _, err := mail.ParseAddress(correo); err != nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "El correo no tiene un formato válido")
		return
	}
	if len(codigo) != 6 {
		c.responderErrorRecuperacion(http.StatusBadRequest, "El código debe tener 6 dígitos")
		return
	}
	if len(solicitud.NuevaContrasena) < 8 {
		c.responderErrorRecuperacion(http.StatusBadRequest, "La nueva contraseña debe tener mínimo 8 caracteres")
		return
	}

	// 2. Buscar el usuario
	usuario, err := buscarUsuarioPorCorreo(baseURL, correo)
	if err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible consultar el usuario: "+err.Error())
		return
	}
	if usuario == nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "Código inválido o expirado")
		return
	}
	idUsuario, ok := idComoEntero(usuario["id"])
	if !ok {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "El CRUD devolvió un usuario sin id")
		return
	}

	// 3. Traer el token más reciente del usuario
	rutaToken := baseURL + "/TokenRecuperacion?limit=1&sortby=Id&order=desc&query=" +
		url.QueryEscape("IdUsuario:"+strconv.Itoa(idUsuario))
	resultadoToken, err := llamarCRUD(http.MethodGet, rutaToken, nil)
	if err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible consultar el código: "+err.Error())
		return
	}
	token := primerRegistro(resultadoToken)
	if token == nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "Código inválido o expirado")
		return
	}

	// 4. Validar: que no esté usado, que no haya expirado y que el hash coincida
	if token["usado_en"] != nil {
		c.responderErrorRecuperacion(http.StatusBadRequest, "Código inválido o expirado")
		return
	}
	expiraTexto, _ := token["expira_en"].(string)
	expira, err := time.Parse(time.RFC3339, expiraTexto)
	if err != nil || time.Now().After(expira) {
		c.responderErrorRecuperacion(http.StatusBadRequest, "Código inválido o expirado")
		return
	}
	hashGuardado, _ := token["token_hash"].(string)
	if subtle.ConstantTimeCompare([]byte(hashGuardado), []byte(hashearCodigo(codigo))) != 1 {
		c.responderErrorRecuperacion(http.StatusBadRequest, "Código inválido o expirado")
		return
	}

	// 5. Marcar el token como usado antes de cambiar la contraseña, para que no se pueda reutilizar.
	// Se envía el objeto completo porque el PUT del CRUD reemplaza toda la fila
	idToken, ok := idComoEntero(token["id"])
	if !ok {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "El CRUD devolvió un token sin id")
		return
	}
	token["usado_en"] = time.Now().Format(time.RFC3339)
	if _, err := llamarCRUD(http.MethodPut, baseURL+"/TokenRecuperacion/"+strconv.Itoa(idToken), token); err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible marcar el código como usado: "+err.Error())
		return
	}

	// 6. Hashear la nueva contraseña
	hash, err := bcrypt.GenerateFromPassword([]byte(solicitud.NuevaContrasena), bcrypt.DefaultCost)
	if err != nil {
		c.responderErrorRecuperacion(http.StatusInternalServerError, "No fue posible procesar la contraseña")
		return
	}

	// 7. Actualizar la contraseña existente (GET + PUT completo) o crearla si el usuario no tenía
	rutaContrasena := baseURL + "/Contrasena?limit=1&sortby=Id&order=desc&query=" +
		url.QueryEscape("IdUsuario:"+strconv.Itoa(idUsuario))
	resultadoContrasena, err := llamarCRUD(http.MethodGet, rutaContrasena, nil)
	if err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible consultar la contraseña: "+err.Error())
		return
	}
	existente := primerRegistro(resultadoContrasena)
	if existente != nil {
		idContrasena, ok := idComoEntero(existente["Id"])
		if !ok {
			c.responderErrorRecuperacion(http.StatusInternalServerError, "El CRUD devolvió una contraseña sin id")
			return
		}
		existente["ContrasenaHash"] = string(hash)
		existente["Activa"] = true
		_, err = llamarCRUD(http.MethodPut, baseURL+"/Contrasena/"+strconv.Itoa(idContrasena), existente)
	} else {
		nueva := map[string]interface{}{
			"IdUsuario":      map[string]interface{}{"id": idUsuario},
			"ContrasenaHash": string(hash),
			"Activa":         true,
		}
		_, err = llamarCRUD(http.MethodPost, baseURL+"/Contrasena", nueva)
	}
	if err != nil {
		c.responderErrorRecuperacion(http.StatusBadGateway, "No fue posible guardar la nueva contraseña: "+err.Error())
		return
	}

	c.responderOKRecuperacion("Contraseña actualizada correctamente")
}

// responderErrorRecuperacion responde con el formato de error del MID
func (c *Recuperar_contrasenaController) responderErrorRecuperacion(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": false, "status": status, "message": mensaje}
	c.ServeJSON()
}

// responderOKRecuperacion responde con el formato de éxito del MID
func (c *Recuperar_contrasenaController) responderOKRecuperacion(mensaje string) {
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]interface{}{"success": true, "status": 200, "message": mensaje}
	c.ServeJSON()
}

// llamarCRUD envía una petición al CRUD y devuelve su respuesta decodificada.
// Devuelve error si el CRUD responde success=false
func llamarCRUD(metodo string, ruta string, cuerpo interface{}) (map[string]interface{}, error) {
	var lector io.Reader
	if cuerpo != nil {
		datos, err := json.Marshal(cuerpo)
		if err != nil {
			return nil, err
		}
		lector = bytes.NewReader(datos)
	}

	peticion, err := http.NewRequest(metodo, ruta, lector)
	if err != nil {
		return nil, err
	}
	peticion.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return nil, err
	}
	defer respuesta.Body.Close()

	var resultado map[string]interface{}
	if err := json.NewDecoder(respuesta.Body).Decode(&resultado); err != nil {
		return nil, fmt.Errorf("respuesta inválida del CRUD (HTTP %d): %w", respuesta.StatusCode, err)
	}
	if exito, _ := resultado["success"].(bool); !exito {
		return nil, fmt.Errorf("el CRUD respondió con error: %v", resultado["message"])
	}
	return resultado, nil
}

// primerRegistro devuelve el primer elemento de "data" en una lista del CRUD, o nil si está vacía
func primerRegistro(resultado map[string]interface{}) map[string]interface{} {
	lista, ok := resultado["data"].([]interface{})
	if !ok || len(lista) == 0 {
		return nil
	}
	registro, _ := lista[0].(map[string]interface{})
	return registro
}

// buscarUsuarioPorCorreo consulta /Usuario filtrando por correo; devuelve nil si no existe
func buscarUsuarioPorCorreo(baseURL string, correo string) (map[string]interface{}, error) {
	ruta := baseURL + "/Usuario?limit=1&query=" + url.QueryEscape("correo:"+correo)
	resultado, err := llamarCRUD(http.MethodGet, ruta, nil)
	if err != nil {
		return nil, err
	}
	return primerRegistro(resultado), nil
}

// idComoEntero convierte un id leído del JSON (float64) a int
func idComoEntero(valor interface{}) (int, bool) {
	numero, ok := valor.(float64)
	return int(numero), ok
}

// generarCodigo crea un código numérico de 6 dígitos con un generador seguro
func generarCodigo() (string, error) {
	numero, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", numero.Int64()), nil
}

// hashearCodigo devuelve el SHA-256 del código en hexadecimal
func hashearCodigo(codigo string) string {
	suma := sha256.Sum256([]byte(codigo))
	return hex.EncodeToString(suma[:])
}