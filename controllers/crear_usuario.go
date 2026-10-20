package controllers

import (
	"API_MID_AGROCAMPO/models"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	beego "github.com/beego/beego/v2/server/web"
	"golang.org/x/crypto/bcrypt"
)

// CrearUsuarioController metodo para CrearUsuario
type CrearUsuarioController struct {
	beego.Controller
}

func (c *CrearUsuarioController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

func (c *CrearUsuarioController) crear(url string, datos interface{}) ([]byte, int, error) {
	cuerpo, err := json.Marshal(datos)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	respuesta, err := http.Post(url, "application/json", bytes.NewBuffer(cuerpo))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer respuesta.Body.Close()

	resultado, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if respuesta.StatusCode < http.StatusOK || respuesta.StatusCode >= http.StatusMultipleChoices {
		return nil, respuesta.StatusCode, errors.New("respuesta no exitosa del api")
	}

	return resultado, http.StatusOK, nil
}

// CrearUsuario ...
// @Title Crear usuario
// @Description Recibe el json del usuario y lo crea en el CRUD de usuarios y contrasenas
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

	if solicitud.Nombre == "" {
		c.responderError(http.StatusBadRequest, "El campo nombre es obligatorio")
		return
	}

	if solicitud.Apellido == "" {
		c.responderError(http.StatusBadRequest, "El campo apellido es obligatorio")
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

	if solicitud.Contrasena == "" {
		c.responderError(http.StatusBadRequest, "El campo contrasena es obligatorio")
		return
	}

	urlUsuarios, _ := beego.AppConfig.String("url_usuarios")
	if urlUsuarios == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_usuarios no esta configurada en el app.conf")
		return
	}

	urlContrasenas, _ := beego.AppConfig.String("url_contrasenas")
	if urlContrasenas == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_contrasenas no esta configurada en el app.conf")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(solicitud.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		c.responderError(http.StatusInternalServerError, "No fue posible procesar la contrasena")
		return
	}

	usuarioCrud := models.UsuarioCrudRequest{
		Nombre:     solicitud.Nombre,
		Apellido:   solicitud.Apellido,
		Correo:     solicitud.Correo,
		Telefono:   solicitud.Telefono,
		IdRol:      models.RefCrud{Id: solicitud.IdRol},
		Avatar:     solicitud.Avatar,
		Verificado: false,
		Activo:     true,
	}

	cuerpoUsuario, status, err := c.crear(urlUsuarios, usuarioCrud)
	if err != nil {
		c.responderError(status, "Api CRUD de usuario no creo el usuario")
		return
	}

	var usuarioCreado models.UsuarioCrudResponse
	err = json.Unmarshal(cuerpoUsuario, &usuarioCreado)
	if err != nil || usuarioCreado.Id == 0 {
		c.responderError(http.StatusInternalServerError, "No fue posible procesar la respuesta del api usuarios")
		return
	}

	contrasenaCrud := models.ContrasenaCrudRequest{
		IdUsuario:      models.RefCrud{Id: usuarioCreado.Id},
		ContrasenaHash: string(hash),
		Activa:         true,
	}

	_, status, err = c.crear(urlContrasenas, contrasenaCrud)
	if err != nil {
		c.responderError(status, "El usuario fue creado, pero no fue posible crear la contrasena")
		return
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = models.CrearUsuarioResponse{
		Mensaje:   "Usuario creado correctamente",
		IDUsuario: usuarioCreado.Id,
	}
	c.ServeJSON()
}