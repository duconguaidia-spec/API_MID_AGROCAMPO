package controllers

import (
	"api_mid/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

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

// GetAll ...
// @Title GetAll
// @Description get Home
// @Param	id	query	string	true	"Id del usuario (temporal, luego sale del token)"
// @Success 200 {object} models.HomeResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @router / [get]
func (c *HomeController) GetAll() {
	id, err := strconv.Atoi(c.GetString("id"))
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
	_, _ = c.obtener(urlPerfiles+"/"+strconv.Itoa(id), &perfil)

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = models.HomeResponse{
		Perfil: models.PerfilHome{
			IDUsuario:      usuario.Id,
			NombreCompleto: usuario.NombreCompleto,
			Correo:         usuario.Correo,
			Telefono:       usuario.Telefono,
			IdRol:          usuario.IdRol,
			Avatar:         usuario.Avatar,
			Ubicacion:      perfil.Ubicacion,
		},
	}
	c.ServeJSON()
}