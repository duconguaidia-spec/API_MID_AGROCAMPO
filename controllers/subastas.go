package controllers

import (
	"API_MID_AGROCAMPO/models"
	"API_MID_AGROCAMPO/services"
	"encoding/json"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// SubastasController operations for Subastas
type SubastasController struct {
	beego.Controller
}

// validarSubasta revisa los datos antes de mandarlos al CRUD
// devuelve un texto con el problema, o "" si todo esta bien
func validarSubasta(s models.Subasta, esNueva bool) string {
	if esNueva && s.RegistradaPor <= 0 {
		return "Falta registrada_por (el id del usuario que registra la subasta)"
	}
	if esNueva && strings.TrimSpace(s.FechaSubasta) == "" {
		return "Falta la fecha de la subasta"
	}
	if len(s.Nombre) > 180 {
		return "El nombre no puede tener mas de 180 letras"
	}
	if len(s.Ubicacion) > 180 {
		return "La ubicacion no puede tener mas de 180 letras"
	}
	return ""
}

// GetAll trae las subastas del CRUD
// se puede usar ?limit=20&offset=0&sortby=fecha_subasta&order=desc
func (c *SubastasController) GetAll() {
	respuesta, err := services.Pedir("GET", "/subasta"+parametrosUrl(&c.Controller), nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	subastas := []models.Subasta{}
	json.Unmarshal(respuesta.Data, &subastas)
	for i := range subastas {
		subastas[i].FechaSubasta = fechaDeCrud(subastas[i].FechaSubasta)
	}
	responder(&c.Controller, 200, true, "Subastas consultadas", subastas)
}

// GetOne trae una subasta con todos sus precios y el promedio
func (c *SubastasController) GetOne() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("GET", "/subasta/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}
	subasta := models.Subasta{}
	json.Unmarshal(respuesta.Data, &subasta)
	subasta.FechaSubasta = fechaDeCrud(subasta.FechaSubasta)

	// ahora los precios de esa subasta
	respuestaPrecios, err := services.Pedir("GET", "/tr_precio_subasta_ganado?query=subasta_id:"+id+"&limit=100", nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	precios := []models.PrecioSubasta{}
	json.Unmarshal(respuestaPrecios.Data, &precios)

	// promedio de los precios
	promedio := 0.0
	if len(precios) > 0 {
		suma := 0.0
		for _, precio := range precios {
			suma = suma + precio.PrecioPorKg
		}
		promedio = suma / float64(len(precios))
	}

	datos := map[string]interface{}{
		"subasta":  subasta,
		"precios":  precios,
		"promedio": promedio,
	}
	responder(&c.Controller, 200, true, "Subasta consultada", datos)
}

// Post valida y crea una subasta
func (c *SubastasController) Post() {
	var subasta models.Subasta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subasta)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarSubasta(subasta, true)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	// el CRUD necesita la fecha con hora, el cliente la manda solo como AAAA-MM-DD
	subasta.FechaSubasta, err = fechaParaCrud(subasta.FechaSubasta)
	if err != nil {
		responder(&c.Controller, 400, false, err.Error(), nil)
		return
	}

	respuesta, err := services.Pedir("POST", "/subasta", subasta)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	creada := models.Subasta{}
	json.Unmarshal(respuesta.Data, &creada)
	creada.FechaSubasta = fechaDeCrud(creada.FechaSubasta)
	responder(&c.Controller, 201, true, "Subasta creada", creada)
}

// Put valida y modifica una subasta
func (c *SubastasController) Put() {
	id := c.Ctx.Input.Param(":id")

	var subasta models.Subasta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &subasta)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarSubasta(subasta, false)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	// en la modificacion la fecha es opcional
	if subasta.FechaSubasta != "" {
		subasta.FechaSubasta, err = fechaParaCrud(subasta.FechaSubasta)
		if err != nil {
			responder(&c.Controller, 400, false, err.Error(), nil)
			return
		}
	}

	numero, _ := strconv.Atoi(id)
	subasta.Id = numero

	respuesta, err := services.Pedir("PUT", "/subasta/"+id, subasta)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	modificada := models.Subasta{}
	json.Unmarshal(respuesta.Data, &modificada)
	modificada.FechaSubasta = fechaDeCrud(modificada.FechaSubasta)
	responder(&c.Controller, 200, true, "Subasta modificada", modificada)
}

// Delete elimina una subasta (el CRUD solo la desactiva)
func (c *SubastasController) Delete() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("DELETE", "/subasta/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}
