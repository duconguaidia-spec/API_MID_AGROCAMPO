package controllers

import (
	"API_MID_AGROCAMPO/services"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type CuentaController struct{ ModulosController }

type cuentaSesion struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Correo   string `json:"correo"`
	Rol      string `json:"rol"`
}

func (u cuentaSesion) gestionaContenido() bool {
	return u.Rol == "administrador" || u.Rol == "veterinario"
}

type sesion struct {
	usuario cuentaSesion
	vence   time.Time
}

var sesiones = struct {
	sync.Mutex
	datos map[string]sesion
}{datos: map[string]sesion{}}

const nombreSesion = "agrocampo_sesion"

func sesionActual(r *http.Request) (cuentaSesion, bool) {
	cookie, err := r.Cookie(nombreSesion)
	if err != nil {
		return cuentaSesion{}, false
	}
	sesiones.Lock()
	defer sesiones.Unlock()
	s, ok := sesiones.datos[cookie.Value]
	return s.usuario, ok && time.Now().Before(s.vence)
}

func (c *CuentaController) Operar() {
	r := c.Ctx.Request
	accion := c.Ctx.Input.Param(":accion")
	usuario, ok := sesionActual(r)
	if accion == "sesion" && r.Method == http.MethodGet {
		if !ok {
			c.responder(401, nil, "Inicia sesión")
			return
		}
		c.responder(200, usuario, "")
		return
	}
	if accion == "sesion" && r.Method == http.MethodDelete {
		if cookie, err := r.Cookie(nombreSesion); err == nil {
			sesiones.Lock()
			delete(sesiones.datos, cookie.Value)
			sesiones.Unlock()
		}
		c.cookie("", -1)
		c.responder(200, nil, "")
		return
	}
	datos := map[string]interface{}{}
	if r.Method != http.MethodGet && (json.Unmarshal(c.Ctx.Input.RequestBody, &datos) != nil || datos == nil) {
		c.responder(400, nil, "Datos inválidos")
		return
	}
	switch accion {
	case "registro", "login", "recuperar", "restablecer":
		if r.Method != http.MethodPost {
			c.responder(405, nil, "Operación no disponible")
			return
		}
		delete(datos, "id")
	case "perfil":
		if !ok {
			c.responder(401, nil, "Inicia sesión para acceder al perfil")
			return
		}
		if r.Method == http.MethodGet {
			datos = map[string]interface{}{}
		} else if r.Method == http.MethodPut {
			accion = "guardar-perfil"
		} else {
			c.responder(405, nil, "Operación no disponible")
			return
		}
		datos["id"] = usuario.ID
		delete(datos, "contrasena")
		delete(datos, "token")
	default:
		c.responder(404, nil, "Operación no disponible")
		return
	}
	var token string
	if accion == "recuperar" {
		correo, _ := datos["correo"].(string)
		correo = strings.ToLower(strings.TrimSpace(correo))
		validado, err := mail.ParseAddress(correo)
		if err != nil || validado.Address != correo {
			c.responder(400, nil, "Ingresa un correo válido")
			return
		}
		token = tokenAleatorio()
		datos = map[string]interface{}{"correo": correo, "token": hashToken(token)}
		accion = "crear-token"
	} else if accion == "restablecer" {
		valor, _ := datos["token"].(string)
		if _, err := hex.DecodeString(valor); err != nil || len(valor) != 64 {
			c.responder(400, nil, "Enlace inválido o expirado")
			return
		}
		datos["token"] = hashToken(valor)
	}
	cuerpo, _ := json.Marshal(datos)
	if len(beego.AppConfig.DefaultString("clave_crud_usuarios", "")) < 32 {
		c.responder(502, nil, "Falta configurar la conexión interna con Usuarios")
		return
	}
	resp, status, err := peticionModulo(r, "crud_usuarios", "cuenta/"+accion, http.MethodPost, cuerpo)
	if status == http.StatusForbidden {
		logs.Warn("Usuarios rechazó la autenticación interna. Comprueba que la MID y Usuarios lean la misma clave")
		c.responder(502, nil, "La MID y Usuarios no pudieron autenticarse entre sí")
		return
	}
	if err != nil || status >= 500 {
		c.responder(502, nil, "El servicio de usuarios no está disponible")
		return
	}
	if status >= 400 || !resp.Success {
		if status < 400 {
			status = 400
		}
		c.responder(status, nil, resp.Message)
		return
	}
	switch accion {
	case "login":
		if json.Unmarshal(resp.Data, &usuario) != nil || usuario.ID <= 0 {
			c.responder(502, nil, "La cuenta no pudo verificarse")
			return
		}
		clave := tokenAleatorio()
		if actual, err := r.Cookie(nombreSesion); err == nil {
			sesiones.Lock()
			delete(sesiones.datos, actual.Value)
			sesiones.Unlock()
		}
		sesiones.Lock()
		for k, s := range sesiones.datos {
			if time.Now().After(s.vence) {
				delete(sesiones.datos, k)
			}
		}
		sesiones.datos[clave] = sesion{usuario, time.Now().Add(8 * time.Hour)}
		sesiones.Unlock()
		c.cookie(clave, 8*60*60)
		c.responder(200, usuario, "")
		return
	case "crear-token":
		var cuenta struct{ Encontrada bool }
		if json.Unmarshal(resp.Data, &cuenta) != nil {
			c.responder(502, nil, "No se pudo preparar la recuperación")
			return
		}
		if cuenta.Encontrada {
			if _, err := os.Stat("token/token.json"); err != nil {
				c.responder(502, nil, "El correo de recuperación no está configurado")
				return
			}
			servicio, err := services.NewGmailService()
			if err == nil {
				base := strings.TrimRight(beego.AppConfig.DefaultString("frontend_publico", "http://localhost:4200"), "/")
				_, err = servicio.SendEmail(datos["correo"].(string), "Recupera tu cuenta de AgroCampo", "Abre este enlace en los próximos 15 minutos:\n"+base+"/#/restablecer-contrasena/"+token)
			}
			if err != nil {
				c.responder(502, nil, "No se pudo enviar el correo de recuperación")
				return
			}
		}
		c.responder(200, nil, "Si el correo está registrado, recibirás un enlace de recuperación")
		return
	case "restablecer":
		var cuenta cuentaSesion
		json.Unmarshal(resp.Data, &cuenta)
		sesiones.Lock()
		for clave, s := range sesiones.datos {
			if s.usuario.ID == cuenta.ID {
				delete(sesiones.datos, clave)
			}
		}
		sesiones.Unlock()
		c.responder(200, nil, "Contraseña actualizada")
		return
	case "guardar-perfil":
		if actual, err := r.Cookie(nombreSesion); err == nil {
			sesiones.Lock()
			if s, existe := sesiones.datos[actual.Value]; existe {
				s.usuario.Nombre, _ = datos["nombre"].(string)
				s.usuario.Apellido, _ = datos["apellido"].(string)
				sesiones.datos[actual.Value] = s
			}
			sesiones.Unlock()
		}
		c.responder(200, nil, "")
		return
	}
	c.responder(status, resp.Data, resp.Message)
}

func (c *CuentaController) cookie(valor string, duracion int) {
	http.SetCookie(c.Ctx.ResponseWriter, &http.Cookie{Name: nombreSesion, Value: valor, Path: "/api", HttpOnly: true, SameSite: http.SameSiteStrictMode, Secure: c.Ctx.Request.TLS != nil || beego.AppConfig.DefaultBool("cookie_segura", false), MaxAge: duracion})
}

func tokenAleatorio() string { var b [32]byte; rand.Read(b[:]); return hex.EncodeToString(b[:]) }
func hashToken(valor string) string {
	suma := sha256.Sum256([]byte(valor))
	return hex.EncodeToString(suma[:])
}
