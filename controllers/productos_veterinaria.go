package controllers

import (
	"api_mid_agrocampo/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// ProductosVeterinariaController consulta los productos de una veterinaria.
type ProductosVeterinariaController struct {
	beego.Controller
}

// GetProductosPorVeterinaria busca los productos usando el ID de la ruta.
func (c *ProductosVeterinariaController) GetProductosPorVeterinaria() {
	// Lee el ID de la veterinaria que viene en la dirección.
	idTexto := c.Ctx.Input.Param(":id")
	veterinariaID, err := strconv.Atoi(idTexto)
	if err != nil || veterinariaID <= 0 {
		c.CustomAbort(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}

	// Obtiene la dirección de la CRUD desde la configuración del MID.
	urlCRUD, err := beego.AppConfig.String("url_productos_veterinaria")
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "Falta configurar la URL de la CRUD de productos")
		return
	}

	// Agrega a la petición el filtro por veterinaria.
	endpoint, err := url.Parse(urlCRUD)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "La URL de la CRUD no es válida")
		return
	}

	parametros := endpoint.Query()
	parametros.Set("query", fmt.Sprintf("IdVeterinaria__Id:%d", veterinariaID))
	parametros.Set("limit", "100")
	endpoint.RawQuery = parametros.Encode()

	// Envía la petición GET a la CRUD.
	peticion, err := http.NewRequestWithContext(
		c.Ctx.Request.Context(),
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		c.CustomAbort(http.StatusInternalServerError, "No se pudo preparar la consulta")
		return
	}

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		c.CustomAbort(http.StatusBadGateway, "No fue posible conectar con la CRUD de productos")
		return
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		c.CustomAbort(http.StatusBadGateway, "La CRUD de productos respondió con un error")
		return
	}

	// Convierte la respuesta JSON de la CRUD a las estructuras de models.
	var resultadoCRUD models.RespuestaProductosCRUD
	if err := json.NewDecoder(respuesta.Body).Decode(&resultadoCRUD); err != nil {
		c.CustomAbort(http.StatusBadGateway, "No se pudo leer la respuesta de la CRUD")
		return
	}

	if !resultadoCRUD.Success {
		c.CustomAbort(http.StatusBadGateway, "La CRUD no pudo consultar los productos")
		return
	}

	// Prepara los datos para el frontend y deja solo los productos activos.
	productos := make([]models.ProductoVeterinaria, 0, len(resultadoCRUD.Data))
	for _, producto := range resultadoCRUD.Data {
		if !producto.Activo {
			continue
		}

		productos = append(productos, models.ProductoVeterinaria{
			ID:            producto.ID,
			VeterinariaID: producto.Veterinaria.ID,
			Nombre:        producto.Nombre,
			Descripcion:   producto.Descripcion,
			ImagenURL:     producto.ImagenURL,
			Precio:        producto.Precio,
			Activo:        producto.Activo,
		})
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    productos,
	}
	c.ServeJSON()
}