package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type Producto struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	ImagenURL   string  `json:"imagen_url"`
	Precio      float64 `json:"precio"`
	Activo      bool    `json:"activo"`
	Veterinaria struct {
		ID int `json:"id"`
	} `json:"veterinaria"`
}

type RespuestaCRUD struct {
	Success bool       `json:"success"`
	Data    []Producto `json:"data"`
}

type ProductosVeterinariaController struct {
	beego.Controller
}

func (c *ProductosVeterinariaController) GetProductos() {
	idTexto := c.Ctx.Input.Param(":id")
	veterinariaID, err := strconv.Atoi(idTexto)
	if err != nil || veterinariaID <= 0 {
		c.CustomAbort(http.StatusBadRequest, "ID de veterinaria no válido")
		return
	}

	urlCRUD, err := beego.AppConfig.String("url_productos_veterinaria")
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "Falta configurar la URL de la CRUD")
		return
	}

	endpoint, err := url.Parse(urlCRUD)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "La URL de la CRUD no es válida")
		return
	}

	parametros := endpoint.Query()
	parametros.Set("query", fmt.Sprintf("IdVeterinaria__Id:%d", veterinariaID))
	endpoint.RawQuery = parametros.Encode()

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Get(endpoint.String())
	if err != nil {
		c.CustomAbort(http.StatusBadGateway, "No fue posible conectar con la CRUD")
		return
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		c.CustomAbort(http.StatusBadGateway, "La CRUD respondió con un error")
		return
	}

	var resultado RespuestaCRUD
	if err := json.NewDecoder(respuesta.Body).Decode(&resultado); err != nil {
		c.CustomAbort(http.StatusBadGateway, "No se pudo leer la respuesta de la CRUD")
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    resultado.Data,
	}
	c.ServeJSON()
}