package controllers

import (
	"API_MID_AGROCAMPO/models"
	"bytes"
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

// GetAllVeterinarias lista las clínicas activas registradas en la CRUD.
func (c *VeterinariaController) GetAllVeterinarias() {
	urlVeterinaria, err := beego.AppConfig.String("url_veterinaria")
	if err != nil {
		c.responderError(http.StatusInternalServerError, "Falta configurar la URL de la CRUD de veterinarias")
		return
	}

	parametros := url.Values{}
	parametros.Set("limit", "100")
	var respuesta models.RespuestaVeterinariasCRUD
	if err := consultarCRUD(c.Ctx.Request.Context(), urlVeterinaria, parametros, &respuesta); err != nil {
		c.responderError(http.StatusBadGateway, "No fue posible consultar las veterinarias en la CRUD")
		return
	}
	if !respuesta.Success {
		c.responderError(http.StatusBadGateway, "La CRUD no pudo consultar las veterinarias")
		return
	}

	veterinarias := make([]models.VeterinariaCRUD, 0, len(respuesta.Data))
	for _, veterinaria := range respuesta.Data {
		if veterinaria.Activo {
			servicios, err := c.consultarServicios(veterinaria.ID)
			if err != nil {
				c.responderError(http.StatusBadGateway, "No fue posible consultar los servicios de las veterinarias")
				return
			}
			veterinaria.Servicios = servicios
			veterinarias = append(veterinarias, veterinaria)
		}
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    veterinarias,
	}
	c.ServeJSON()
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

// PostCrearVeterinaria crea la veterinaria y luego sus productos en las CRUD.
func (c *VeterinariaController) PostCrearVeterinaria() {
	var solicitud models.CrearVeterinariaRequest

	// Lee el JSON enviado por el frontend.
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(&solicitud); err != nil {
		c.responderError(http.StatusBadRequest, "El JSON enviado no es válido")
		return
	}

	// Valida los campos obligatorios de la veterinaria.
	if strings.TrimSpace(solicitud.Nombre) == "" ||
		strings.TrimSpace(solicitud.Direccion) == "" ||
		strings.TrimSpace(solicitud.Ciudad) == "" ||
		strings.TrimSpace(solicitud.Telefono) == "" {
		c.responderError(
			http.StatusBadRequest,
			"Nombre, dirección, ciudad y teléfono son obligatorios",
		)
		return
	}

	// Valida los productos antes de empezar a guardar datos.
	for i, producto := range solicitud.Productos {
		if strings.TrimSpace(producto.Nombre) == "" {
			c.responderError(
				http.StatusBadRequest,
				fmt.Sprintf("El producto en la posición %d debe tener nombre", i),
			)
			return
		}
		if producto.Precio < 0 {
			c.responderError(
				http.StatusBadRequest,
				fmt.Sprintf("El precio del producto en la posición %d no puede ser negativo", i),
			)
			return
		}
	}

	// Lee las direcciones de las CRUD desde conf/app.conf.
	urlVeterinaria, err := beego.AppConfig.String("url_veterinaria")
	if err != nil {
		c.responderError(
			http.StatusInternalServerError,
			"Falta configurar url_veterinaria",
		)
		return
	}

	urlProductos, err := beego.AppConfig.String("url_productos_veterinaria")
	if err != nil {
		c.responderError(
			http.StatusInternalServerError,
			"Falta configurar url_productos_veterinaria",
		)
		return
	}

	// Prepara la veterinaria para enviarla a su CRUD.
	veterinariaCRUD := models.CrearVeterinariaCRUD{
		PropietarioID:   solicitud.PropietarioID,
		Nombre:          solicitud.Nombre,
		Descripcion:     solicitud.Descripcion,
		Direccion:       solicitud.Direccion,
		Ciudad:          solicitud.Ciudad,
		Departamento:    solicitud.Departamento,
		Telefono:        solicitud.Telefono,
		CorreoPublico:   solicitud.CorreoPublico,
		HorarioAtencion: solicitud.HorarioAtencion,
		Latitud:         solicitud.Latitud,
		Longitud:        solicitud.Longitud,
		Aprobado:        false,
		Activo:          true,
	}

	var respuestaVeterinaria models.RespuestaVeterinariaCRUD
	_, err = enviarPOSTCRUD(
		c.Ctx.Request.Context(),
		urlVeterinaria,
		veterinariaCRUD,
		&respuestaVeterinaria,
	)
	if err != nil {
		c.responderError(
			http.StatusBadGateway,
			"No fue posible crear la veterinaria en la CRUD",
		)
		return
	}

	if !respuestaVeterinaria.Success || respuestaVeterinaria.Data.ID <= 0 {
		c.responderError(
			http.StatusBadGateway,
			"La CRUD no devolvió el ID de la veterinaria creada",
		)
		return
	}

	idVeterinaria := respuestaVeterinaria.Data.ID
	productosCreados := make([]models.ProductoCRUD, 0, len(solicitud.Productos))

	// Crea cada producto usando el ID que devolvió la CRUD de veterinarias.
	for i, producto := range solicitud.Productos {
		productoCRUD := models.CrearProductoVeterinariaCRUD{
			Veterinaria: models.ReferenciaVeterinariaCRUD{
				ID: idVeterinaria,
			},
			Nombre:      producto.Nombre,
			Descripcion: producto.Descripcion,
			ImagenURL:   producto.ImagenURL,
			Precio:      producto.Precio,
			Activo:      true,
		}

		var respuestaProducto models.RespuestaCrearProductoCRUD
		_, err := enviarPOSTCRUD(
			c.Ctx.Request.Context(),
			urlProductos,
			productoCRUD,
			&respuestaProducto,
		)
		if err != nil || !respuestaProducto.Success {
			// La veterinaria y los productos anteriores ya pueden haberse guardado.
			c.Ctx.Output.SetStatus(http.StatusBadGateway)
			c.Data["json"] = map[string]interface{}{
				"success":          false,
				"message":          "La veterinaria se creó, pero falló la creación de uno de sus productos",
				"producto_fallido": i,
				"data": map[string]interface{}{
					"veterinaria": respuestaVeterinaria.Data,
					"productos":   productosCreados,
				},
			}
			c.ServeJSON()
			return
		}

		productosCreados = append(productosCreados, respuestaProducto.Data)
	}

	// Devuelve al frontend todo lo que se creó.
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = models.CrearVeterinariaResponse{
		Success: true,
		Message: "Veterinaria y productos creados correctamente",
		Data: models.CrearVeterinariaDataResponse{
			Veterinaria: respuestaVeterinaria.Data,
			Productos:   productosCreados,
		},
	}
	c.ServeJSON()
}

// enviarPOSTCRUD envía JSON a una CRUD y decodifica su respuesta.
func enviarPOSTCRUD(
	ctx context.Context,
	endpoint string,
	datos interface{},
	destino interface{},
) (int, error) {
	cuerpoJSON, err := json.Marshal(datos)
	if err != nil {
		return 0, err
	}

	peticion, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(cuerpoJSON),
	)
	if err != nil {
		return 0, err
	}

	peticion.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{Timeout: 10 * time.Second}
	respuesta, err := cliente.Do(peticion)
	if err != nil {
		return 0, err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode < http.StatusOK ||
		respuesta.StatusCode >= http.StatusMultipleChoices {
		return respuesta.StatusCode, fmt.Errorf(
			"la CRUD respondió con HTTP %d",
			respuesta.StatusCode,
		)
	}

	if err := json.NewDecoder(respuesta.Body).Decode(destino); err != nil {
		return respuesta.StatusCode, err
	}

	return respuesta.StatusCode, nil
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
