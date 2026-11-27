package controllers

import (
	"API_MID_AGROCAMPO/models"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

// El CRUD actualiza filas completas. Se conserva el registro original y se
// modifican solo los campos editables para no borrar relaciones ni metadatos.
func (c *VeterinariaController) DeleteVeterinaria() {
	id, ok := c.idVeterinaria()
	if !ok {
		c.responderError(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}
	endpoint, _ := beego.AppConfig.String("url_veterinaria")
	endpoint = strings.TrimRight(endpoint, "/") + "/" + strconv.Itoa(id)
	if _, ok := c.obtenerRegistro(endpoint, id); !ok {
		return
	}
	if c.modificarRegistro(http.MethodDelete, endpoint, nil) {
		c.responderDatos(http.StatusOK, nil)
	}
}

func (c *VeterinariaController) PutVeterinaria() {
	if status, mensaje := autorizarEdicion(c.Ctx.Request); status != 0 {
		c.responderError(status, mensaje)
		return
	}
	id, ok := c.idVeterinaria()
	if !ok {
		c.responderError(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}
	var cambios map[string]json.RawMessage
	if json.Unmarshal(c.Ctx.Input.RequestBody, &cambios) != nil || len(cambios) == 0 {
		c.responderError(http.StatusBadRequest, "Envía los campos de la veterinaria en formato JSON")
		return
	}
	endpoint, _ := beego.AppConfig.String("url_veterinaria")
	endpoint = strings.TrimRight(endpoint, "/") + "/" + strconv.Itoa(id)
	registro, ok := c.obtenerRegistro(endpoint, id)
	if !ok {
		return
	}
	for _, campo := range []string{"nombre", "descripcion", "direccion", "ciudad", "departamento", "telefono", "correo_publico", "horario_atencion"} {
		if valor, existe := cambios[campo]; existe {
			var texto string
			if string(valor) == "null" || json.Unmarshal(valor, &texto) != nil {
				c.responderError(http.StatusBadRequest, "El campo "+campo+" debe ser texto")
				return
			}
			registro[campo] = valor
		}
	}
	for _, campo := range []string{"nombre", "direccion", "ciudad", "telefono"} {
		var texto string
		json.Unmarshal(registro[campo], &texto)
		if strings.TrimSpace(texto) == "" {
			c.responderError(http.StatusBadRequest, "Nombre, dirección, ciudad y teléfono son obligatorios")
			return
		}
	}
	if c.modificarRegistro(http.MethodPut, endpoint, registro) {
		c.responderDatos(http.StatusOK, registro)
	}
}

func (c *VeterinariaController) PostProducto() { c.guardarProducto(false) }

func (c *VeterinariaController) PutProducto() {
	if status, mensaje := autorizarEdicion(c.Ctx.Request); status != 0 {
		c.responderError(status, mensaje)
		return
	}
	c.guardarProducto(true)
}

func (c *VeterinariaController) guardarProducto(editar bool) {
	id, ok := c.idVeterinaria()
	if !ok {
		c.responderError(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}
	var producto models.CrearProductoRequest
	if json.Unmarshal(c.Ctx.Input.RequestBody, &producto) != nil || strings.TrimSpace(producto.Nombre) == "" || producto.Precio < 0 {
		c.responderError(http.StatusBadRequest, "Envía un producto con nombre y precio mayor o igual a cero")
		return
	}
	endpoint, _ := beego.AppConfig.String("url_productos_veterinaria")
	endpoint = strings.TrimRight(endpoint, "/")
	metodo := http.MethodPost
	var datos interface{} = models.CrearProductoVeterinariaCRUD{
		Veterinaria: models.ReferenciaVeterinariaCRUD{ID: id},
		Nombre:      strings.TrimSpace(producto.Nombre), Descripcion: producto.Descripcion,
		ImagenURL: producto.ImagenURL, Precio: producto.Precio, Activo: true,
	}
	if editar {
		registro, urlProducto, ok := c.obtenerProducto(endpoint, id)
		if !ok {
			return
		}
		for campo, valor := range map[string]interface{}{"nombre": strings.TrimSpace(producto.Nombre), "descripcion": producto.Descripcion, "imagen_url": producto.ImagenURL, "precio": producto.Precio} {
			registro[campo], _ = json.Marshal(valor)
		}
		endpoint, datos, metodo = urlProducto, registro, http.MethodPut
	} else {
		urlVeterinaria, _ := beego.AppConfig.String("url_veterinaria")
		if _, ok := c.obtenerRegistro(strings.TrimRight(urlVeterinaria, "/")+"/"+strconv.Itoa(id), id); !ok {
			return
		}
	}
	var respuesta struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
	}
	if _, err := enviarJSONCRUD(c.Ctx.Request.Context(), metodo, endpoint, datos, &respuesta); err != nil || !respuesta.Success {
		c.responderError(http.StatusBadGateway, "La CRUD no pudo guardar el producto")
		return
	}
	if editar {
		c.responderDatos(http.StatusOK, datos)
	} else {
		var creado models.ProductoCRUD
		if json.Unmarshal(respuesta.Data, &creado) != nil || creado.ID <= 0 {
			c.responderError(http.StatusBadGateway, "La CRUD no confirmó el ID del producto creado")
			return
		}
		c.responderDatos(http.StatusCreated, creado)
	}
}

func (c *VeterinariaController) DeleteProducto() {
	id, ok := c.idVeterinaria()
	if !ok {
		c.responderError(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}
	endpoint, _ := beego.AppConfig.String("url_productos_veterinaria")
	_, endpoint, ok = c.obtenerProducto(strings.TrimRight(endpoint, "/"), id)
	if ok && c.modificarRegistro(http.MethodDelete, endpoint, nil) {
		c.responderDatos(http.StatusOK, nil)
	}
}

func (c *VeterinariaController) obtenerProducto(endpoint string, idVeterinaria int) (map[string]json.RawMessage, string, bool) {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":productoId"))
	if err != nil || id <= 0 {
		c.responderError(http.StatusBadRequest, "El ID del producto no es válido")
		return nil, "", false
	}
	endpoint += "/" + strconv.Itoa(id)
	registro, ok := c.obtenerRegistro(endpoint, id)
	if !ok {
		return nil, "", false
	}
	var veterinaria models.ReferenciaVeterinariaCRUD
	if json.Unmarshal(registro["veterinaria"], &veterinaria) != nil || veterinaria.ID != idVeterinaria {
		c.responderError(http.StatusNotFound, "El producto no pertenece a esta veterinaria")
		return nil, "", false
	}
	return registro, endpoint, true
}

func (c *VeterinariaController) obtenerRegistro(endpoint string, id int) (map[string]json.RawMessage, bool) {
	var respuesta struct {
		Success bool                       `json:"success"`
		Data    map[string]json.RawMessage `json:"data"`
	}
	if err := consultarCRUD(c.Ctx.Request.Context(), endpoint, nil, &respuesta); err != nil {
		c.responderError(http.StatusBadGateway, "No fue posible consultar el registro en la CRUD")
		return nil, false
	}
	var recibido int
	json.Unmarshal(respuesta.Data["id"], &recibido)
	if !respuesta.Success || recibido != id {
		c.responderError(http.StatusNotFound, "No se encontró el registro solicitado")
		return nil, false
	}
	return respuesta.Data, true
}

func (c *VeterinariaController) modificarRegistro(metodo, endpoint string, datos interface{}) bool {
	var respuesta struct {
		Success bool `json:"success"`
	}
	if _, err := enviarJSONCRUD(c.Ctx.Request.Context(), metodo, endpoint, datos, &respuesta); err != nil || !respuesta.Success {
		c.responderError(http.StatusBadGateway, "La CRUD no confirmó la operación")
		return false
	}
	return true
}

func (c *VeterinariaController) responderDatos(status int, datos interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": true, "data": datos}
	c.ServeJSON()
}
