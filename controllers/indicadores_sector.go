package controllers

import (
	"API_MID_AGROCAMPO/models"
	"API_MID_AGROCAMPO/services"
	"encoding/json"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// IndicadoresSectorController operations for IndicadorSector
type IndicadoresSectorController struct {
	beego.Controller
}

// validarIndicador revisa los datos antes de mandarlos al CRUD
func validarIndicador(i models.IndicadorSector, esNuevo bool) string {
	if esNuevo && strings.TrimSpace(i.Nombre) == "" {
		return "Falta el nombre del indicador"
	}
	if esNuevo && strings.TrimSpace(i.Unidad) == "" {
		return "Falta la unidad del indicador"
	}
	if len(i.Nombre) > 180 {
		return "El nombre no puede tener mas de 180 letras"
	}
	if len(i.Unidad) > 60 {
		return "La unidad no puede tener mas de 60 letras"
	}
	return ""
}

// GetAll trae los indicadores del CRUD
func (c *IndicadoresSectorController) GetAll() {
	respuesta, err := services.Pedir("GET", "/indicador_sector"+parametrosUrl(&c.Controller), nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	indicadores := []models.IndicadorSector{}
	json.Unmarshal(respuesta.Data, &indicadores)
	responder(&c.Controller, 200, true, "Indicadores consultados", indicadores)
}

// GetOne trae un indicador con sus ultimas 12 mediciones
func (c *IndicadoresSectorController) GetOne() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("GET", "/indicador_sector/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}
	indicador := models.IndicadorSector{}
	json.Unmarshal(respuesta.Data, &indicador)

	// las mediciones de ese indicador, de la mas nueva a la mas vieja
	ruta := "/medicion_sector?query=indicador_id:" + id + "&sortby=fecha&order=desc&limit=12"
	respuestaMediciones, err := services.Pedir("GET", ruta, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	mediciones := []models.MedicionSector{}
	json.Unmarshal(respuestaMediciones.Data, &mediciones)
	for i := range mediciones {
		mediciones[i].Fecha = fechaDeCrud(mediciones[i].Fecha)
	}

	// la primera de la lista es la mas reciente
	datos := map[string]interface{}{
		"indicador":    indicador,
		"mediciones":   mediciones,
		"tiene_datos":  len(mediciones) > 0,
		"ultimo_valor": 0.0,
	}
	if len(mediciones) > 0 {
		datos["ultimo_valor"] = mediciones[0].Valor
	}
	responder(&c.Controller, 200, true, "Indicador consultado", datos)
}

// Post valida y crea un indicador
func (c *IndicadoresSectorController) Post() {
	var indicador models.IndicadorSector
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &indicador)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarIndicador(indicador, true)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	respuesta, err := services.Pedir("POST", "/indicador_sector", indicador)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}

// Put valida y modifica un indicador
func (c *IndicadoresSectorController) Put() {
	id := c.Ctx.Input.Param(":id")

	var indicador models.IndicadorSector
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &indicador)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarIndicador(indicador, false)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	numero, _ := strconv.Atoi(id)
	indicador.Id = numero

	respuesta, err := services.Pedir("PUT", "/indicador_sector/"+id, indicador)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}

// Delete elimina un indicador y sus mediciones (el CRUD solo los desactiva)
func (c *IndicadoresSectorController) Delete() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("DELETE", "/indicador_sector/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}
