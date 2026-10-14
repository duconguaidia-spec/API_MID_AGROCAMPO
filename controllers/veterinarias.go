package controllers

import (
	"api_mid_AgroCampo/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// VeterinariaController prepara los datos que necesita la página de detalle.
type VeterinariaController struct {
	beego.Controller
}

// GetDetalleVeterinaria reúne la clínica, sus servicios y sus productos.
func (c *VeterinariaController) GetDetalleVeterinaria() {
	id, ok := c.idVeterinaria()
	if !ok {
		c.responderError(http.StatusBadRequest, "El ID de la veterinaria no es válido")
		return
	}

	// Primero se consulta la veterinaria. Si no existe, no tiene sentido pedir sus relaciones.
	urlVeterinaria, err := beego.AppConfig.String("url_veterinaria")
	if err != nil {
		c.responderError(http.StatusInternalServerError, "Falta configurar la URL de la CRUD de veterinarias")
		return
	}

	var respuestaVeterinaria models.RespuestaVeterinariaCRUD
	if err := consultarCRUD(c.Ctx.Request.Context(), urlVeterinaria+"/"+strconv.Itoa(id), nil, &respuestaVeterinaria); err != nil {
		c.responderError(http.StatusBadGateway, "No fue posible consultar la veterinaria en la CRUD")
		return
	}
	if !respuestaVeterinaria.Success || respuestaVeterinaria.Data.ID != id {
		c.responderError(http.StatusNotFound, "No se encontró la veterinaria solicitada")
		return
	}

	// Estas consultas se filtran por el mismo ID para traer solo datos de esta clínica.
	productos, err := c.consultarProductos(id)
	if err != nil {
		c.responderError(http.StatusBadGateway, "No fue posible consultar los productos de la veterinaria")
		return
	}

	servicios, err := c.consultarServicios(id)
	if err != nil {
		c.responderError(http.StatusBadGateway, "No fue posible consultar los servicios de la veterinaria")
		return
	}

	// El MID entrega un solo objeto, aunque internamente haya consultado varias rutas CRUD.
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data": models.DetalleVeterinaria{
			Veterinaria: respuestaVeterinaria.Data,
			Servicios:   servicios,
			Productos:   productos,
		},
	}
	c.ServeJSON()
}

func (c *VeterinariaController) idVeterinaria() (int, bool) {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	return id, err == nil && id > 0
}

func (c *VeterinariaController) consultarProductos(id int) ([]models.ProductoVeterinaria, error) {
	urlCRUD, err := beego.AppConfig.String("url_productos_veterinaria")
	if err != nil {
		return nil, err
	}

	parametros := url.Values{}
	parametros.Set("query", fmt.Sprintf("IdVeterinaria__Id:%d", id))
	parametros.Set("limit", "100")

	var respuesta models.RespuestaProductosCRUD
	if err := consultarCRUD(c.Ctx.Request.Context(), urlCRUD, parametros, &respuesta); err != nil {
		return nil, err
	}
	if !respuesta.Success {
		return nil, fmt.Errorf("la CRUD indicó que la consulta de productos falló")
	}

	productos := make([]models.ProductoVeterinaria, 0, len(respuesta.Data))
	for _, producto := range respuesta.Data {
		if !producto.Activo {
			continue
		}
		productos = append(productos, models.ProductoVeterinaria{
			ID:          producto.ID,
			Nombre:      producto.Nombre,
			Descripcion: producto.Descripcion,
			ImagenURL:   producto.ImagenURL,
			Precio:      producto.Precio,
		})
	}
	return productos, nil
}

func (c *VeterinariaController) consultarServicios(id int) ([]models.ServicioVeterinaria, error) {
	// Se consultan los servicios propios de la clínica.
	urlServicios, err := beego.AppConfig.String("url_servicios_veterinaria")
	if err != nil {
		return nil, err
	}

	parametros := url.Values{}
	parametros.Set("query", fmt.Sprintf("IdVeterinaria__Id:%d", id))
	parametros.Set("limit", "100")

	var respuestaServicios models.RespuestaServiciosCRUD
	if err := consultarCRUD(c.Ctx.Request.Context(), urlServicios, parametros, &respuestaServicios); err != nil {
		return nil, err
	}
	if !respuestaServicios.Success {
		return nil, fmt.Errorf("la CRUD indicó que la consulta de servicios falló")
	}

	servicios := make([]models.ServicioVeterinaria, 0, len(respuestaServicios.Data))
	for _, servicio := range respuestaServicios.Data {
		if servicio.Activo {
			servicios = append(servicios, models.ServicioVeterinaria{
				ID:          servicio.ID,
				Nombre:      servicio.Nombre,
				Descripcion: servicio.Descripcion,
				Precio:      servicio.Precio,
			})
		}
	}

	// También se consultan los servicios estándar asociados mediante la tabla intermedia.
	urlRelaciones, err := beego.AppConfig.String("url_relaciones_servicios_veterinaria")
	if err != nil {
		return nil, err
	}

	var respuestaRelaciones models.RespuestaRelacionesServiciosCRUD
	if err := consultarCRUD(c.Ctx.Request.Context(), urlRelaciones, parametros, &respuestaRelaciones); err != nil {
		return nil, err
	}
	if !respuestaRelaciones.Success {
		return nil, fmt.Errorf("la CRUD indicó que la consulta de servicios asociados falló")
	}

	nombres := make(map[string]bool, len(servicios))
	for _, servicio := range servicios {
		nombres[strings.ToLower(servicio.Nombre)] = true
	}
	for _, relacion := range respuestaRelaciones.Data {
		nombre := strings.TrimSpace(relacion.ServicioGeneral.Nombre)
		if relacion.Activo && relacion.ServicioGeneral.Activo && nombre != "" && !nombres[strings.ToLower(nombre)] {
			servicios = append(servicios, models.ServicioVeterinaria{Nombre: nombre})
			nombres[strings.ToLower(nombre)] = true
		}
	}

	return servicios, nil
}

// consultarCRUD hace una petición GET y convierte su JSON en la estructura indicada.
func consultarCRUD(ctx context.Context, endpoint string, parametros url.Values, destino interface{}) error {
	direccion, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if parametros != nil {
		consulta := direccion.Query()
		for clave, valores := range parametros {
			for _, valor := range valores {
				consulta.Add(clave, valor)
			}
		}
		direccion.RawQuery = consulta.Encode()
	}

	peticion, err := http.NewRequestWithContext(ctx, http.MethodGet, direccion.String(), nil)
	if err != nil {
		return err
	}

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode < http.StatusOK || respuesta.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("la CRUD respondió con estado HTTP %d", respuesta.StatusCode)
	}
	return json.NewDecoder(respuesta.Body).Decode(destino)
}

func (c *VeterinariaController) responderError(estado int, mensaje string) {
	c.Ctx.Output.SetStatus(estado)
	c.Data["json"] = map[string]interface{}{
		"success": false,
		"message": mensaje,
	}
	c.ServeJSON()
}
