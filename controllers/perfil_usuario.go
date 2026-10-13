package controllers

import (
	"API_MID_AGROCAMPO/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// Perfil_usuarioController operations for Perfil_usuario
type Perfil_usuarioController struct {
	beego.Controller
}

func (c *Perfil_usuarioController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

func (c *Perfil_usuarioController) obtener(url string, destino interface{}) (int, error) {
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

// GetOne ...
// @Title GetOne
// @Description get Perfil_usuario by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Perfil_usuario
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Perfil_usuarioController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		c.responderError(http.StatusBadRequest, "El id del usuario no es valido")
		return
	}

	urlUsuarios, _ := beego.AppConfig.String("url_usuarios")
	if urlUsuarios == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_usuarios no esta configurada en el app.conf")
		return
	}

	urlPerfiles, _ := beego.AppConfig.String("url_perfiles")
	if urlPerfiles == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_perfiles no esta configurada en el app.conf")
		return
	}

	var usuario models.UsuarioPerfilCrud
	status, err := c.obtener(urlUsuarios+"/"+strconv.Itoa(id), &usuario)
	if err != nil {
		c.responderError(status, "No fue posible consultar el usuario")
		return
	}

	if !usuario.Activo {
		c.responderError(http.StatusForbidden, "El usuario se encuentra inactivo")
		return
	}

	var perfil models.PerfilCrud
	status, err = c.obtener(urlPerfiles+"/"+strconv.Itoa(id), &perfil)
	if err != nil {
		c.responderError(status, "No fue posible consultar el perfil del usuario")
		return
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = models.PerfilUsuarioResponse{
		IDUsuario:            usuario.Id,
		NombreCompleto:       usuario.NombreCompleto,
		Correo:               usuario.Correo,
		Telefono:             usuario.Telefono,
		IdRol:                usuario.IdRol,
		VerificacionDosPasos: usuario.VerificacionDosPasos,
		Avatar:               usuario.Avatar,
		Ubicacion:            perfil.Ubicacion,
		Activo:               usuario.Activo,
	}
	c.ServeJSON()
}