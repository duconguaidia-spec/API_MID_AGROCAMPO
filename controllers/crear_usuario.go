package controllers

import (
	"API_MID_AGROCAMPO/models"
	service "API_MID_AGROCAMPO/services"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/mail"

	"github.com/beego/beego/v2/core/logs"
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

	// El CRUD responde 200 aunque la operacion falle, por eso se revisa el campo success
	var respuestaCrud models.RespuestaCrud
	err = json.Unmarshal(resultado, &respuestaCrud)
	if err != nil || !respuestaCrud.Success {
		return nil, http.StatusBadGateway, errors.New("el api CRUD no completo la operacion: " + respuestaCrud.Message)
	}

	return resultado, http.StatusOK, nil
}

// enviarBienvenida envia por Gmail el mensaje de bienvenida al usuario recien creado
func (c *CrearUsuarioController) enviarBienvenida(correo string, nombre string) error {
	servicioGmail, err := service.NewGmailService()
	if err != nil {
		return err
	}

	// El asunto se codifica para que la tilde y los signos se vean bien en el correo
	asunto := mime.QEncoding.Encode("UTF-8", "¡Bienvenido a AgroCampo!")
	cuerpo := fmt.Sprintf(
		"Hola %s,\n\n¡Bienvenido a AgroCampo! Tu cuenta fue creada con éxito y ya puedes iniciar sesión con este correo.\n\nNos alegra tenerte con nosotros.",
		nombre,
	)

	_, err = servicioGmail.SendEmail(correo, asunto, cuerpo)
	return err
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

	// El correo se usa como destinatario del mensaje, por eso debe ser una direccion valida y limpia
	direccion, err := mail.ParseAddress(solicitud.Correo)
	if err != nil || direccion.Address != solicitud.Correo {
		c.responderError(http.StatusBadRequest, "El campo correo no tiene un formato valido")
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
	if err != nil || usuarioCreado.Data.Id == 0 {
		c.responderError(http.StatusInternalServerError, "No fue posible procesar la respuesta del api usuarios")
		return
	}

	contrasenaCrud := models.ContrasenaCrudRequest{
		IdUsuario:      models.RefCrud{Id: usuarioCreado.Data.Id},
		ContrasenaHash: string(hash),
		Activa:         true,
	}

	_, status, err = c.crear(urlContrasenas, contrasenaCrud)
	if err != nil {
		c.responderError(status, "El usuario fue creado, pero no fue posible crear la contrasena")
		return
	}

	// El usuario ya esta creado: si el correo de bienvenida falla solo se avisa en la respuesta
	correoEnviado := true
	err = c.enviarBienvenida(solicitud.Correo, solicitud.Nombre)
	if err != nil {
		correoEnviado = false
		logs.Warn("No fue posible enviar el correo de bienvenida a %s: %v", solicitud.Correo, err)
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = models.CrearUsuarioResponse{
		Mensaje:                 "Usuario creado correctamente",
		IDUsuario:               usuarioCreado.Data.Id,
		CorreoBienvenidaEnviado: correoEnviado,
	}
	c.ServeJSON()
}