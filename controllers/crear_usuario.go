package controllers

import (
	"api_mid/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// CrearUsuarioController operacion para CrearUsuario
type CrearUsuarioController struct {
	beego.Controller
}

// responderError centraliza la respuesta de error
func (c *CrearUsuarioController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

// CrearUsuario ...
// @Title Crear usuario
// @Description Recibe el json del usuario y lo crea en el CRUD de usuarios
// @Param body body models.CrearUsuarioRequest true "Informacion del usuario"
// @Success 201 {object} models.CrearUsuarioResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @router /crear [post]
func (c *CrearUsuarioController) CrearUsuario() {
	var solicitud models.CrearUsuarioRequest

	err := json.Unmarshal(c.Ctx.Input.RequestBody, &solicitud)
	if err != nil {
		c.responderError(http.StatusBadRequest, "El Json enviado no tiene el formato valido")
		return
	}

	if solicitud.NombreCompleto == "" {
		c.responderError(http.StatusBadRequest, "El campo nombre completo es obligatorio")
		return
	}

	if solicitud.Correo == "" {
		c.responderError(http.StatusBadRequest, "El campo correo es obligatorio")
		return
	}

	if solicitud.IdRol <= 0 {
		c.responderError(http.StatusBadRequest, "El campo id rol es obligatorio")
		return
	}

	urlUsuarios, _ := beego.AppConfig.String("url_usuarios")
	if urlUsuarios == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_usuarios no esta configurada en el app.conf")
		return
	}

	usuarioCrud := models.UsuarioCrudRequest{
		NombreCompleto:       solicitud.NombreCompleto,
		Correo:               solicitud.Correo,
		Telefono:             solicitud.Telefono,
		VerificacionDosPasos: solicitud.VerificacionDosPasos,
		Avatar:               solicitud.Avatar,
		Activo:               solicitud.Activo,
	}

	jsonUsuario, err := json.Marshal(usuarioCrud)
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No fue posible convertir el json de usuarios")
		return
	}

	respuestaUsuarios, err := http.Post(urlUsuarios, "application/json", bytes.NewBuffer(jsonUsuario))
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No se pudo realizar el post de usuario")
		return
	}
	defer respuestaUsuarios.Body.Close()

	if respuestaUsuarios.StatusCode < http.StatusOK || respuestaUsuarios.StatusCode >= http.StatusMultipleChoices {
		c.responderError(respuestaUsuarios.StatusCode, "Api CRUD de usuario no creo el usuario")
		return
	}

	cuerpoUsuarioCreado, err := io.ReadAll(respuestaUsuarios.Body)
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No fue posible leer la respuesta del api usuarios")
		return
	}

	var usuarioCreado models.UsuarioCrudResponse
	err = json.Unmarshal(cuerpoUsuarioCreado, &usuarioCreado)
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No fue posible procesar la respuesta del api usuarios")
		return
	}

	resultado := models.CrearUsuarioResponse{
		Mensaje:   "Usuario creado correctamente",
		IDUsuario: usuarioCreado.Id,
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = resultado
	c.ServeJSON()
}