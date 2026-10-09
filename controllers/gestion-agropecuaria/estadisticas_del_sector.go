package controllers

import (

	"encoding/json"
	"strconv"
	
	"agrocampo_crud_gestion/models"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
)

// Estadisticas_del_sectorController operations for Estadisticas_del_sector
type Estadisticas_del_sectorController struct {
	beego.Controller
}

// Post ...
// @Title Create
// @Description create Estadisticas_del_sector
// @Param	body		body 	models.Estadisticas_del_sector	true		"body for Estadisticas_del_sector content"
// @Success 201 {object} models.Estadisticas_del_sector
// @Failure 403 body is empty
// @router / [post]
// URLMapping ...
func (c *Estadisticas_del_sectorController)Post() {

	var dato models.Estadisticas_del_sector
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


//getone trae un registro por id
// GetOne ...
// @Title GetOne
// @Description get Estadisticas_del_sector by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Estadisticas_del_sector
// @Failure 403 :id is empty
// @router /:id [get]

func (c *Estadisticas_del_sectorController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil{
		responder(&c.Controller, 400, false, "El id no es valido", nil)
		return
	}

	o := orm.NewOrm()
	dato := models.Estadisticas_del_sector{Id: id}
	err = o.Read(&dato)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe el registro", nil)
		return
	}

	if err != nil {
		responder(&c.Controller, 404, false, "Error consultado en el registro",err.Error())
		return
	}
	if !dato.Activo {
		responder (&c.Controller,404, false, "No existe el registro", nil)
		return
	
	}
		responder(&c.Controller, 200, true,"Registro consultado", dato)

	}


// GetAll ...
// @Title GetAll
// @Description get Estadisticas_del_sector
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Estadisticas_del_sector
// @Failure 403
// @router / [get]
func (c *Estadisticas_del_sectorController) GetAll() {
	o := orm.NewOrm()
	datos :=[]models.Estadisticas_del_sector{}

	consulta := o.QueryTable(new(models.Estadistica_del_sector)).Filter("activo", true)

	fecha := c.GetString("fecha")
	if fecha != ""{
		consulta =consulta.Filter("fecha", fecha)
	}
	categoria := c.GetString("categoria")
	if categoria != ""{
		consulta  = consulta.Filter("categoria", categoria)
	}
	_, err := consulta.OrderBy("-fecha", "categoria").All(&datos)
	if err != nil{
		responder(&c.Controller, 500, false, "Error consultando las registros", err.Error())
		return
	}

	responder(&c.Controller, 200, true, "Registros buscados",datos)
	
}

// Put ...
// @Title Put
// @Description update the Estadisticas_del_sector
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Estadisticas_del_sector	true		"body for Estadisticas_del_sector content"
// @Success 200 {object} models.Estadisticas_del_sector
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Estadisticas_del_sectorController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil{
		responder(&c.Controller, 400,false, "El id no es valido",nil)
		return
	}

	o := orm.NewOrm()
	dato := models.Estadisticas_del_sector{Id:id}
	err = o.Read(&dato)
	if err != nil {
		responder(&c.Controller, 404, false, "No existe el registro",nil)
		return
	}

	err = json.Unmarshal(c.Ctx.Input.RequestBody, &dato)
	if err !=nil{
		responder(&c.Controller, 400, false, "El JSON enviado no es valido",err.Error())
		return
	}
	dato.Id = id
	dato.Activo = true

	_, err= o.Update(&dato)
	if err !=nil{
		responder(&c.Controller, 500, false, "No se pudo modificar el registro",err.Error())
		return
	}
	responder(&c.Controller,200, true, "Registro modificxdp", dato)
}

// Delete ...
// @Title Delete
// @Description delete the Estadisticas_del_sector
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Estadisticas_del_sectorController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err !=nil{
		responder(&c.Controller, 400, false,"El id no es valido", nil)
		return
	}

	o := orm.NewOrm()
	dato := models.Estadisticas_del_sector{Id:id}
	err = o.Read(&dato)
	if err != nil {
		responder(&c.Controller, 404, false, "No existe el registro",nil)
		return
	}

	dato.Activo = false
	_, err= o.Update(&dato, "Activo")
	if err !=nil{
		responder(&c.Controller, 500, false, "No se pudo eliminar el registro",err.Error())
		return
	}
	responder(&c.Controller,200, true, "Registro eliminado ",nil)
}

