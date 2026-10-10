package controllers

import (
	"encoding/json"
	"strconv"

	"agrocampo_crud_gestion/models"


	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

)

// Estadisticas_del_ganadoController operations for Estadisticas_del_ganado
type Estadisticas_del_ganadoController struct {
	beego.Controller
}


// Post ...
// @Title Create
// @Description create Estadisticas_del_ganado
// @Param	body		body 	models.Estadisticas_del_ganado	true		"body for Estadisticas_del_ganado content"
// @Success 201 {object} models.Estadisticas_del_ganado
// @Failure 403 body is empty
// @router / [post]
func (c *Estadisticas_del_ganadoController) Post() {
	var dato models.Estadisticadelganado
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}

	dato.Id =0
	dato.Activo = true

	o := orm.NewOrm()
	_,err = o.Insert(&dato)
	if err != nil{
		responder(&c.Controller, 500, false, "No se puedo crear el registro" , err.Error())
		return
	}
	responder(&c.Controller,201, true, "Registro creado", dato)

}



// GetOne ...
// @Title GetOne
// @Description get Estadisticas_del_ganado by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Estadisticas_del_ganado
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Estadisticas_del_ganadoController) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get Estadisticas_del_ganado
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Estadisticas_del_ganado
// @Failure 403
// @router / [get]
func (c *Estadisticas_del_ganadoController) GetAll() {

}

// Put ...
// @Title Put
// @Description update the Estadisticas_del_ganado
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Estadisticas_del_ganado	true		"body for Estadisticas_del_ganado content"
// @Success 200 {object} models.Estadisticas_del_ganado
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Estadisticas_del_ganadoController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Estadisticas_del_ganado
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Estadisticas_del_ganadoController) Delete() {

}
