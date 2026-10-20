package controllers

import (
	"API_MID_AGROCAMPO/models"
	"API_MID_AGROCAMPO/services"
	"encoding/json"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// PreciosSubastaController operations for PrecioSubasta
type PreciosSubastaController struct {
	beego.Controller
}

// validarPrecio revisa los datos antes de mandarlos al CRUD
func validarPrecio(p models.PrecioSubasta, esNuevo bool) string {
	if esNuevo && p.SubastaId <= 0 {
		return "Falta subasta_id"
	}
	if esNuevo && p.CategoriaGanadoId <= 0 {
		return "Falta categoria_ganado_id"
	}
	if p.PrecioPorKg <= 0 {
		return "El precio por kilo debe ser mayor que cero"
	}
	return ""
}

// GetAll trae los precios del CRUD
// para ver los de una subasta: ?query=subasta_id:1
func (c *PreciosSubastaController) GetAll() {
	respuesta, err := services.Pedir("GET", "/tr_precio_subasta_ganado"+parametrosUrl(&c.Controller), nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	if !respuesta.Success {
		devolverRespuestaCrud(&c.Controller, respuesta)
		return
	}

	precios := []models.PrecioSubasta{}
	json.Unmarshal(respuesta.Data, &precios)
	responder(&c.Controller, 200, true, "Precios consultados", precios)
}

// GetOne trae un precio por id
func (c *PreciosSubastaController) GetOne() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("GET", "/tr_precio_subasta_ganado/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}

// Post valida y crea un precio
func (c *PreciosSubastaController) Post() {
	var precio models.PrecioSubasta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &precio)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarPrecio(precio, true)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	respuesta, err := services.Pedir("POST", "/tr_precio_subasta_ganado", precio)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}

// Put valida y modifica el precio (solo cambia el precio y la unidad)
func (c *PreciosSubastaController) Put() {
	id := c.Ctx.Input.Param(":id")

	var precio models.PrecioSubasta
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &precio)
	if err != nil {
		responder(&c.Controller, 400, false, "El JSON enviado no es valido", err.Error())
		return
	}

	mensaje := validarPrecio(precio, false)
	if mensaje != "" {
		responder(&c.Controller, 400, false, mensaje, nil)
		return
	}

	numero, _ := strconv.Atoi(id)
	precio.Id = numero

	respuesta, err := services.Pedir("PUT", "/tr_precio_subasta_ganado/"+id, precio)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}

// Delete elimina un precio (el CRUD solo lo desactiva)
func (c *PreciosSubastaController) Delete() {
	id := c.Ctx.Input.Param(":id")

	respuesta, err := services.Pedir("DELETE", "/tr_precio_subasta_ganado/"+id, nil)
	if err != nil {
		fallaCrud(&c.Controller, err)
		return
	}
	devolverRespuestaCrud(&c.Controller, respuesta)
}
