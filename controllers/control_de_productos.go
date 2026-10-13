package controllers

import (
	"API_MID_AGROCAMPO/controllers"
	"encoding/json"
	"strconv"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

)

// Control_de_productosController operations for Control_de_productos
type Control_de_productosController struct {
	beego.Controller
}



// Post ...
// @Title Create
// @Description create Control_de_productos
// @Param	body		body 	models.Control_de_productos	true		"body for Control_de_productos content"
// @Success 201 {object} models.Control_de_productos
// @Failure 403 body is empty
// @router / [post]
func (c *Control_de_productosController) Post() {
	var control_de_productos models.Control_de_productos
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &control_de_productos)
	if err != nil {
		responder(&c.Controller, 400, false, "El json enviado no es valido", err.Error())
		return
	}

	control_de_productos.Id = 0
	control_de_productos.Activo = true

	o := orm.NewOrm()
	_, err = o.Insert(&control_de_productos)
	if err != nil {
		responder(&c.Controller, 500, false, "No se pudo crear el producto", err.Error())
		return
	}
	responder(&c.Controller, 201, true, "Producto creado", control_de_productos)

}

// GetOne ...
// @Title GetOne
// @Description get Control_de_productos by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Control_de_productos
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Control_de_productosController) GetOne() {
		id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		responder(&c.Controller, 400, false, "El id no es valido ", err.Error())
		return
	}

	o := orm.NewOrm()
	Control_de_productos:= models.Control_de_productos{Id: id}
	err = o.Read(&Control_de_productos)
	if err == orm.ErrNoRows {
		responder(&c.Controller, 404, false, "No existe el producto", err.Error())
		return
	}
	if err != nil {
		responder(&c.Controller, 500, false, "Error consultado en el producto", nil)
		return
	}
	if !Control_de_productos.Activo {
		responder(&c.Controller, 404, false, "No existe el producto", nil)
		return
	}

	responder(&c.Controller, 200, true, "Producto consultado", Control_de_productos)

}

// GetAll ...
// @Title GetAll
// @Description get Control_de_productos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Control_de_productos
// @Failure 403
// @router / [get]
func (c *Control_de_productosController) GetAll() {
	o := orm.NewOrm()
	datos := []models.Control_de_productos{}

	consulta := o.QueryTable(new(models.Control_de_productos)).Filter("activo", true)

	fecha := c.GetString("fecha")
	if fecha != "" {
		consulta = consulta.Filter("fecha", fecha)
	}
	categoria := c.GetString("categoria")
	if categoria != "" {
		consulta = consulta.Filter("categoria", categoria)
	}
	_, err := consulta.OrderBy("-fecha", "categoria").All(&datos)
	if err != nil {
		responder(&c.Controller, 500, false, "Error consultando las registros", err.Error())
		return
	}

	responder(&c.Controller, 200, true, "Registros buscados", datos)

}


// Put ...
// @Title Put
// @Description update the Control_de_productos
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Control_de_productos	true		"body for Control_de_productos content"
// @Success 200 {object} models.Control_de_productos
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Control_de_productosController) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the Control_de_productos
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Control_de_productosController) Delete() {

}
