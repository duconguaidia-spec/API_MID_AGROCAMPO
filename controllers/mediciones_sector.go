package controllers

import (
	"API_MID_AGROCAMPO/models"
	"API_MID_AGROCAMPO/services"
	"encoding/json"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// MedicionesSectorController operations for MedicionSector
type MedicionesSectorController struct {
	beego.Controller
}

// validarMedicion revisa los datos antes de mandarlos al CRUD
func validarMedicion(m models.MedicionSector, esNueva bool) string {
	if esNueva && m.IndicadorId <= 0 {
		return "Falta indicador_id"
	}
	if esNueva && strings.TrimSpace(m.Fecha) == "" {
		return "Falta la fecha de la medicion"
	}
	if len(m.Region) > 120 {
		return "La region no puede tener mas de 120 letras"
	}
	return ""
}

// GetAll trae las mediciones del CRUD
// para ver las de un indicador: ?query=indicador_id:1&sortby=fecha&order=desc
func (c *MedicionesSectorController) GetAll() {
	respuesta, err := services.Pedir("GET", "/medicion_sector"+parametrosUrl(&c.Controller), nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	mediciones := []models.MedicionSector{}
	json.Unmarshal(respuesta.Data, &mediciones)
	for i := range mediciones {
		mediciones[i].Fecha = fechaDeCrud(mediciones[i].Fecha)
	}
	responder(&c.Controller, 200, true, "Mediciones consultadas", mediciones)
}

// GetOne trae una medicion por id
func (c *MedicionesSectorController) GetOne() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("GET", "/medicion_sector/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	medicion := models.MedicionSector{}
	json.Unmarshal(respuesta.Data, &medicion)
	medicion.Fecha = fechaDeCrud(medicion.Fecha)
	responder(&c.Controller, 200, true, "Medicion consultada", medicion)
}

// Post valida y crea una medicion
func (c *MedicionesSectorController) Post() {
	var medicion models.MedicionSector
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &medicion)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarMedicion(medicion, true)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	// el CRUD necesita la fecha con hora, el cliente la manda solo como AAAA-MM-DD
	medicion.Fecha, err = fechaParaCrud(medicion.Fecha)
	if err != nil {
		responder(&c.Controller, 400, false, err.Error(), nil)
		return
	}

	respuesta, err := services.Pedir("POST", "/medicion_sector", medicion)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	creada := models.MedicionSector{}
	json.Unmarshal(respuesta.Data, &creada)
	creada.Fecha = fechaDeCrud(creada.Fecha)
	responder(&c.Controller, 201, true, "Medicion creada", creada)
}

// Put valida y modifica la medicion (solo cambia el valor y la fuente)
func (c *MedicionesSectorController) Put() {
	id := c.Ctx.Input.Param(":id")

	var medicion models.MedicionSector
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &medicion)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	numero, _ := strconv.Atoi(id)
	medicion.Id = numero

	// la fecha no se puede cambiar, asi que no se manda
	medicion.Fecha = ""

	respuesta, err := services.Pedir("PUT", "/medicion_sector/"+id, medicion)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	modificada := models.MedicionSector{}
	json.Unmarshal(respuesta.Data, &modificada)
	modificada.Fecha = fechaDeCrud(modificada.Fecha)
	responder(&c.Controller, 200, true, "Medicion modificada", modificada)
}

// Delete elimina una medicion (el CRUD solo la desactiva)
func (c *MedicionesSectorController) Delete() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("DELETE", "/medicion_sector/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}
