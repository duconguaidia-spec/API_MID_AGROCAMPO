package controllers

import (
	"API_MID_AGROCAMPO/models"
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// TrPrecioSubastaGanadoController operations for TrPrecioSubastaGanado
type TrPrecioSubastaGanadoController struct {
	beego.Controller
}

// Post ...
// @Title Post
// @Description crear el precio
// @router / [post]
func (c *TrPrecioSubastaGanadoController) Post() {
	var dato models.TrPrecioSubastaGanado
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}

	_, err = models.AddTrPrecioSubastaGanado(&dato)
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo crear el precio", err.Error())
		return
	}
	responder(&c.Controller, 201, true, "Precio creado", dato)
}

// GetOne ...
// @Title GetOne
// @Description consultar el precio por id
// @router /:id [get]
func (c *TrPrecioSubastaGanadoController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	dato, err := models.GetTrPrecioSubastaGanadoById(id)
	if err != nil {
		responder(&c.Controller, 404, false, "No existe el precio", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Precio consultado", dato)
}

// GetAll ...
// @Title GetAll
// @Description consultar los precios
// @router / [get]
func (c *TrPrecioSubastaGanadoController) GetAll() {
	query, fields, sortby, order, offset, limit, err := leerParametrosLista(&c.Controller)
	if err != nil {
		responder(&c.Controller, 400, false, "Los parametros de la consulta no son validos", err.Error())
		return
	}

	lista, err := models.GetAllTrPrecioSubastaGanado(query, fields, sortby, order, offset, limit)
	if err != nil {
		responder(&c.Controller, 400, false, "Error consultando los precios", err.Error())
		return
	}

	// si no hay registros se devuelve una lista vacia, no un error
	if lista == nil {
		lista = []interface{}{}
	}
	responder(&c.Controller, 200, true, "Precios consultados", lista)
}

// Put ...
// @Title Put
// @Description modificar el precio
// @router /:id [put]
func (c *TrPrecioSubastaGanadoController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	dato := models.TrPrecioSubastaGanado{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}
	dato.Id = id

	err = models.UpdateTrPrecioSubastaGanadoById(&dato)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe el precio", nil)
		return
	}
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo modificar el precio", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Precio modificado", dato)
}

// Delete ...
// @Title Delete
// @Description eliminar el precio (solo lo desactiva)
// @router /:id [delete]
func (c *TrPrecioSubastaGanadoController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	err = models.DeleteTrPrecioSubastaGanado(id)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe el precio", nil)
		return
	}
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo eliminar el precio", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Precio eliminado", nil)
}