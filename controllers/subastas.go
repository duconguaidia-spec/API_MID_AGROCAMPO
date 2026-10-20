package controllers

import (
	"API_MID_AGROCAMPO/models"
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// SubastaController operations for Subasta
type SubastasController struct {
	beego.Controller
}

// Post ...
// @Title Post
// @Description crear la subasta
// @router / [post]
func (c *SubastasController) Post() {
	var dato models.Subastas
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}

	_, err = models.AddSubastas(&dato)
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo crear la subasta", err.Error())
		return
	}
	responder(&c.Controller, 201, true, "Subasta creada", dato)
}

// GetOne ...
// @Title GetOne
// @Description consultar la subasta por id
// @router /:id [get]
func (c *SubastasController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	dato, err := models.GetSubastasById(id)
	if err != nil {
		responder(&c.Controller, 404, false, "No existe la subasta", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Subasta consultada", dato)
}

// GetAll ...
// @Title GetAll
// @Description consultar las subastas
// @router / [get]
func (c *SubastasController) GetAll() {
	query, fields, sortby, order, offset, limit, err := leerParametrosLista(&c.Controller)
	if err != nil {
		responder(&c.Controller, 400, false, "Los parametros de la consulta no son validos", err.Error())
		return
	}

	lista, err := models.GetAllSubastas(query, fields, sortby, order, offset, limit)
	if err != nil {
		responder(&c.Controller, 400, false, "Error consultando las subastas", err.Error())
		return
	}

	// si no hay registros se devuelve una lista vacia, no un error
	if lista == nil {
		lista = []interface{}{}
	}
	responder(&c.Controller, 200, true, "Subastas consultadas", lista)
}

// Put ...
// @Title Put
// @Description modificar la subasta
// @router /:id [put]
func (c *SubastasController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	dato := models.Subastas{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}
	dato.Id = id

	err = models.UpdateSubastasById(&dato)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe la subasta", nil)
		return
	}
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo modificar la subasta", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Subasta modificada", dato)
}

// Delete ...
// @Title Delete
// @Description eliminar la subasta (solo lo desactiva)
// @router /:id [delete]
func (c *SubastasController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	err = models.DeleteSubastas(id)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe la subasta", nil)
		return
	}
	if err != nil {
		responder(&c.Controller, 400, false, "No se pudo eliminar la subasta", err.Error())
		return
	}
	responder(&c.Controller, 200, true, "Subasta eliminada", nil)
}