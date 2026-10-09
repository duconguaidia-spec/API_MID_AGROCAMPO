package controllers

import (
	"api_mid/models"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	beego "github.com/beego/beego/v2/server/web"
)

// Recuperar_contrasenaController operations for Recuperar_contrasena
type Recuperar_contrasenaController struct {
	beego.Controller
}

func (c *Recuperar_contrasenaController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

// Post ...
// @Title Create
// @Description create Recuperar_contrasena
// @Param	body		body 	models.Recuperar_contrasena	true		"body for Recuperar_contrasena content"
// @Success 201 {object} models.Recuperar_contrasena
// @Failure 403 body is empty
// @router / [post]
func (c *Recuperar_contrasenaController) Post() {
	var solicitud models.RecuperarContrasenaRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud)
	if err != nil {
		c.responderError(http.StatusBadRequest, "El Json enviado no tiene el formato valido")
		return
	}

	if solicitud.Correo == "" {
		c.responderError(http.StatusBadRequest, "El campo correo es obligatorio")
		return
	}

	urlUsuarios, _ := beego.AppConfig.String("url_usuarios")
	if urlUsuarios == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_usuarios no esta configurada en el app.conf")
		return
	}

	urlConsulta := urlUsuarios + "?query=Correo:" + url.QueryEscape(solicitud.Correo) + "&limit=1"

	respuestaUsuarios, err := http.Get(urlConsulta)
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No se pudo consultar el api de usuarios")
		return
	}
	defer respuestaUsuarios.Body.Close()

	if respuestaUsuarios.StatusCode < http.StatusOK || respuestaUsuarios.StatusCode >= http.StatusMultipleChoices {
		c.responderError(respuestaUsuarios.StatusCode, "Api CRUD de usuario no pudo consultar el usuario")
		return
	}

	
}