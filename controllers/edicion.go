package controllers

import (
	"errors"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

var clienteEdicionInterna = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func autorizarEdicion(r *http.Request) (int, string) {
	usuario, identificado := sesionActual(r)
	if !identificado {
		return http.StatusUnauthorized, "Inicia sesión para editar"
	}
	if !usuario.gestionaContenido() {
		return http.StatusForbidden, "Solo administradores y veterinarios pueden editar"
	}
	return 0, ""
}

func agregarClaveEdicionInterna(r *http.Request) error {
	clave := beego.AppConfig.DefaultString("clave_crud_usuarios", "")
	if len(clave) < 32 {
		return errors.New("falta configurar la clave interna de AgroCampo")
	}
	r.Header.Set("X-Agrocampo-Interno", clave)
	return nil
}
