package controllers

import (
	"API_MID_AGROCAMPO/services"
	"errors"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// devolverRespuestaCrud devuelve tal cual lo que contesto el CRUD
func devolverRespuestaCrud(c *beego.Controller, r services.RespuestaCrud) {
	codigo := r.Status
	if codigo == 0 {
		codigo = 500
	}
	responder(c, codigo, r.Success, r.Message, r.Data)
}

// fallaCrud se usa cuando el MID no pudo hablar con el CRUD
// (casi siempre es porque el CRUD no esta corriendo)
func fallaCrud(c *beego.Controller, err error) {
	responder(c, 502, false, "No se pudo comunicar con el CRUD", err.Error())
}

// parametrosUrl copia lo que venga despues del ? en la URL, para pasarselo al CRUD
// por ejemplo: ?query=indicador_id:1&limit=20
func parametrosUrl(c *beego.Controller) string {
	if c.Ctx.Request.URL.RawQuery == "" {
		return ""
	}
	return "?" + c.Ctx.Request.URL.RawQuery
}

// fechaParaCrud pasa "2026-10-25" a "2026-10-25T00:00:00Z", que es lo que entiende el CRUD
func fechaParaCrud(fecha string) (string, error) {
	t, err := time.Parse("2006-01-02", fecha)
	if err != nil {
		return "", errors.New("la fecha debe tener el formato AAAA-MM-DD")
	}
	return t.Format(time.RFC3339), nil
}

// fechaDeCrud hace lo contrario: de "2026-10-25T00:00:00Z" deja solo "2026-10-25"
func fechaDeCrud(fecha string) string {
	if len(fecha) >= 10 {
		return fecha[:10]
	}
	return fecha
}
