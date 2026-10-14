package controllers_InicioSesion

import (
	"API_MID_AGROCAMPO/models_InicioSesion"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// HomeController operations for Home
type HomeController struct {
	beego.Controller
}

func (c *HomeController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

func (c *HomeController) obtener(url string, destino interface{}) (int, error) {
	respuesta, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode < http.StatusOK || respuesta.StatusCode >= http.StatusMultipleChoices {
		return respuesta.StatusCode, errors.New("respuesta no exitosa del api")
	}

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal(cuerpo, destino)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (c *HomeController) contar(url string) *int {
	var lista []json.RawMessage
	_, err := c.obtener(url+"?limit=-1", &lista)
	if err != nil {
		return nil
	}

	total := len(lista)
	return &total
}

func (c *HomeController) calcularEficiencia(url string) *int {
	var solicitudes []models_InicioSesion.SolicitudCrud
	_, err := c.obtener(url+"?limit=-1", &solicitudes)
	if err != nil {
		return nil
	}

	eficiencia := 0
	if len(solicitudes) > 0 {
		resueltas := 0
		for _, solicitud := range solicitudes {
			if strings.EqualFold(solicitud.Estado, "Resuelta") {
				resueltas++
			}
		}
		eficiencia = resueltas * 100 / len(solicitudes)
	}

	return &eficiencia
}

func validarToken(token string, secreto string) (int, error) {
	partes := strings.Split(token, ".")
	if len(partes) != 3 {
		return 0, errors.New("formato de token invalido")
	}

	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write([]byte(partes[0] + "." + partes[1]))

	firma, err := base64.RawURLEncoding.DecodeString(partes[2])
	if err != nil || !hmac.Equal(firma, mac.Sum(nil)) {
		return 0, errors.New("firma invalida")
	}

	cuerpo, err := base64.RawURLEncoding.DecodeString(partes[1])
	if err != nil {
		return 0, errors.New("payload invalido")
	}

	var claims models_InicioSesion.TokenClaims
	err = json.Unmarshal(cuerpo, &claims)
	if err != nil {
		return 0, errors.New("payload invalido")
	}

	if claims.Exp != 0 && time.Now().Unix() > claims.Exp {
		return 0, errors.New("token expirado")
	}

	if claims.Id <= 0 {
		return 0, errors.New("token sin id de usuario")
	}

	return claims.Id, nil
}

// GetAll ...
// @Title GetAll
// @Description get Home
// @Param	Authorization	header	string	true	"Bearer token del login"
// @Success 200 {object} models_InicioSesion.HomeResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @router / [get]
func (c *HomeController) GetAll() {
	secreto, _ := beego.AppConfig.String("jwt_secret")
	if secreto == "" {
		c.responderError(http.StatusInternalServerError, "La variable jwt_secret no esta configurada en el app.conf")
		return
	}

	token := strings.TrimPrefix(c.Ctx.Input.Header("Authorization"), "Bearer ")
	idUsuario, err := validarToken(token, secreto)
	if err != nil {
		c.responderError(http.StatusUnauthorized, "Token invalido o expirado")
		return
	}

	urls := map[string]string{}
	nombres := []string{"url_usuarios", "url_perfiles", "url_foros", "url_productores", "url_veterinarias", "url_solicitudes", "url_noticias"}
	for _, nombre := range nombres {
		valor, _ := beego.AppConfig.String(nombre)
		if valor == "" {
			c.responderError(http.StatusInternalServerError, "La variable "+nombre+" no esta configurada en el app.conf")
			return
		}
		urls[nombre] = valor
	}

	id := strconv.Itoa(idUsuario)

	var usuario models_InicioSesion.UsuarioPerfilCrud
	status, err := c.obtener(urls["url_usuarios"]+"/"+id, &usuario)
	if err != nil {
		c.responderError(status, "No fue posible consultar el usuario")
		return
	}

	if !usuario.Activo {
		c.responderError(http.StatusForbidden, "El usuario se encuentra inactivo")
		return
	}

	var perfil models_InicioSesion.PerfilCrud
	_, _ = c.obtener(urls["url_perfiles"]+"/"+id, &perfil)

	noticias := []models_InicioSesion.NoticiaCrud{}
	_, err = c.obtener(urls["url_noticias"]+"?sortby=Fecha&order=desc&limit=5", &noticias)
	if err != nil || noticias == nil {
		noticias = []models_InicioSesion.NoticiaCrud{}
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = models_InicioSesion.HomeResponse{
		Perfil: models_InicioSesion.PerfilHome{
			IDUsuario:      usuario.Id,
			NombreCompleto: usuario.NombreCompleto,
			Correo:         usuario.Correo,
			Telefono:       usuario.Telefono,
			IdRol:          usuario.IdRol,
			Avatar:         usuario.Avatar,
			Ubicacion:      perfil.Ubicacion,
		},
		Contadores: models_InicioSesion.ContadoresHome{
			Foros:        c.contar(urls["url_foros"]),
			Productores:  c.contar(urls["url_productores"]),
			Veterinarias: c.contar(urls["url_veterinarias"]),
			Eficiencia:   c.calcularEficiencia(urls["url_solicitudes"]),
		},
		Noticias: noticias,
	}
	c.ServeJSON()
}