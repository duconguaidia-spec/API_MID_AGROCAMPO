package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type ModulosController struct{ beego.Controller }

type recursoModulo struct{ config, ruta string }

var recursosModulo = map[string]recursoModulo{
	"noticias": {"crud_noticias", "noticias"}, "categorias-noticias": {"crud_noticias", "categorias"},
	"categorias-contenido": {"crud_contenido", "categoria"}, "temas": {"crud_contenido", "temaforo"},
	"respuestas": {"crud_contenido", "respuestaforo"}, "videos": {"crud_contenido", "videoeducativo"},
	"preguntas": {"crud_soporte", "preguntafrecuente"}, "tickets": {"crud_soporte", "tickets"},
	"temas-guardados":   {"crud_contenido", "guardados"},
	"categorias-ganado": {"crud_agropecuario", "categoria_ganado"}, "subastas": {"crud_agropecuario", "subasta"},
	"precios-ganado": {"crud_agropecuario", "tr_precio_subasta_ganado"}, "mediciones": {"crud_agropecuario", "sector"},
	"reportes": {"crud_reportes", "reporteactividad"},
}

var clienteModulos = &http.Client{Timeout: 10 * time.Second}

func (c *ModulosController) responder(status int, datos interface{}, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": status < 400, "data": datos, "message": mensaje}
	c.ServeJSON()
}

func (c *ModulosController) Operar() {
	nombre := c.Ctx.Input.Param(":recurso")
	recurso, existe := recursosModulo[nombre]
	if !existe {
		c.responder(404, nil, "Recurso no encontrado")
		return
	}
	metodo := c.Ctx.Request.Method
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		n, err := strconv.Atoi(id)
		if err != nil || n <= 0 {
			c.responder(400, nil, "ID inválido")
			return
		}
	}
	usuario, identificado := sesionActual(c.Ctx.Request)
	if nombre == "temas-guardados" {
		if !identificado {
			c.responder(401, nil, "Inicia sesión para consultar tus temas guardados")
			return
		}
		if id != "" {
			c.responder(405, nil, "Operación no disponible")
			return
		}
		if metodo == http.MethodGet {
			recurso.ruta += "/" + strconv.Itoa(usuario.ID)
		}
	}
	if nombre == "tickets" && metodo != http.MethodPost {
		c.responder(405, nil, "Operación no disponible")
		return
	}
	if nombre == "reportes" && !identificado {
		c.responder(401, nil, "Inicia sesión para consultar o solicitar reportes")
		return
	}
	var cuerpo []byte
	if metodo != http.MethodGet {
		if !identificado {
			c.responder(401, nil, "Inicia sesión para continuar")
			return
		}
		var datos map[string]interface{}
		if json.Unmarshal(c.Ctx.Input.RequestBody, &datos) != nil || datos == nil {
			c.responder(400, nil, "Envía un objeto JSON válido")
			return
		}
		switch nombre {
		case "temas-guardados":
			datos = map[string]interface{}{"usuario_id": usuario.ID, "tema_id": datos["tema_id"], "guardado": datos["guardado"]}
		case "temas", "respuestas":
			if metodo != http.MethodPost || id != "" {
				c.responder(405, nil, "Operación no disponible")
				return
			}
			datos["autor_id"] = usuario.ID
			datos["activo"] = true
			if nombre == "temas" {
				datos["estado"] = "abierto"
				datos["vistas"] = 0
			}
		case "noticias", "videos":
			if metodo != http.MethodPost || id != "" || !usuario.gestionaContenido() {
				c.responder(403, nil, "No tienes permisos para publicar contenido")
				return
			}
			datos["autor_id"] = usuario.ID
			if nombre == "noticias" {
				datos["estado"] = "publicada"
				datos["acceso_limitado"] = false
			} else {
				datos["activo"] = true
			}
		case "tickets":
			datos = map[string]interface{}{"usuario_id": usuario.ID, "asunto": "Incidencia desde AgroCampo", "descripcion": datos["descripcion"]}
		case "reportes":
			if metodo != http.MethodPost || id != "" {
				c.responder(405, nil, "Operación no disponible")
				return
			}
			formato, _ := datos["formato"].(string)
			if formato != "pdf" && formato != "xlsx" {
				c.responder(400, nil, "Selecciona PDF o Excel")
				return
			}
			filas, err := cotizacionesParaReporte(c.Ctx.Request)
			if err != nil {
				c.responder(502, nil, "No fue posible obtener las cotizaciones del reporte")
				return
			}
			if len(filas) == 0 {
				c.responder(409, nil, "Todavía no hay cotizaciones para exportar")
				return
			}
			instantanea, _ := json.Marshal(datosReporte{Cotizaciones: filas})
			datos = map[string]interface{}{"tipo": "ganado", "formato": formato, "estado": "completado", "filtros": string(instantanea), "activo": true, "solicitado_por": usuario.ID, "completado_en": time.Now().UTC()}
		case "precios-ganado":
			if metodo != http.MethodPost || id != "" || usuario.Rol != "administrador" {
				c.responder(403, nil, "No tienes permisos para actualizar precios")
				return
			}
			recurso.ruta = "precio-subasta"
		default:
			c.responder(405, nil, "Operación no disponible")
			return
		}
		delete(datos, "id")
		cuerpo, _ = json.Marshal(datos)
	}
	query := "limit=-1"
	if nombre == "reportes" && metodo == http.MethodGet {
		if id != "" {
			c.responder(405, nil, "Consulta el listado de tus reportes")
			return
		}
		query += "&query=IdUsuarioSolicitante:" + strconv.Itoa(usuario.ID)
	}
	if nombre == "noticias" {
		query += "&query=Estado:publicada,AccesoLimitado:false"
	}
	if id != "" {
		recurso.ruta += "/" + id
	} else if nombre != "temas-guardados" || metodo != http.MethodGet {
		recurso.ruta += "/"
	}
	respuesta, status, err := peticionModulo(c.Ctx.Request, recurso.config, recurso.ruta+"?"+query, metodo, cuerpo)
	if err != nil {
		c.responder(502, nil, "No fue posible comunicarse con el servicio")
		return
	}
	if status >= 400 || !respuesta.Success {
		if status >= 500 {
			status = 502
		} else if status < 400 {
			status = 400
		}
		mensaje := respuesta.Message
		if status == 502 || mensaje == "" {
			mensaje = "El servicio no pudo completar la operación"
		}
		c.responder(status, nil, mensaje)
		return
	}
	if metodo == http.MethodGet && (len(respuesta.Data) == 0 || string(respuesta.Data) == "null") {
		if id != "" {
			c.responder(404, nil, "Registro no encontrado")
			return
		}
		respuesta.Data = json.RawMessage(`[]`)
	}
	if nombre == "noticias" && id != "" {
		var noticia struct {
			Estado         string
			AccesoLimitado bool `json:"acceso_limitado"`
		}
		if json.Unmarshal(respuesta.Data, &noticia) != nil || noticia.Estado != "publicada" || noticia.AccesoLimitado {
			c.responder(404, nil, "Noticia no encontrada")
			return
		}
	}
	if nombre == "reportes" && metodo == http.MethodGet {
		var reportes []reporteGanado
		if json.Unmarshal(respuesta.Data, &reportes) != nil {
			c.responder(502, nil, "No fue posible leer los reportes")
			return
		}
		publicos := make([]map[string]interface{}, 0, len(reportes))
		for _, reporte := range reportes {
			if reporte.Usuario != usuario.ID || !reporte.Activo {
				continue
			}
			datos := map[string]interface{}{"id": reporte.ID, "formato": reporte.Formato, "estado": reporte.Estado, "fecha_creacion": reporte.Fecha}
			if reporte.Estado == "completado" && reporte.instantaneaValida() {
				datos["archivo_url"] = "/api/v1/reportes/" + strconv.Itoa(reporte.ID) + "/archivo"
			}
			publicos = append(publicos, datos)
		}
		c.responder(200, publicos, "")
		return
	}
	if informacion := respuesta.Data; len(informacion) > 0 {
		c.responder(status, informacion, "")
	} else {
		c.responder(status, nil, "")
	}
}

type respuestaModulo struct {
	Success bool            `json:"success"`
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func peticionModulo(origen *http.Request, config, ruta, metodo string, datos []byte) (respuestaModulo, int, error) {
	var resultado respuestaModulo
	base := strings.TrimRight(beego.AppConfig.DefaultString(config, ""), "/")
	req, err := http.NewRequestWithContext(origen.Context(), metodo, base+"/"+ruta, bytes.NewReader(datos))
	if err != nil {
		return resultado, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if config == "crud_usuarios" {
		req.Header.Set("X-Agrocampo-Interno", beego.AppConfig.DefaultString("clave_crud_usuarios", ""))
	}
	resp, err := clienteModulos.Do(req)
	if err != nil {
		return resultado, 0, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(io.LimitReader(resp.Body, 16*1024*1024)).Decode(&resultado)
	estado := resp.StatusCode
	if estado < 400 && resultado.Status >= 200 && resultado.Status <= 599 && (!resultado.Success || resultado.Status < 400) {
		estado = resultado.Status
	}
	return resultado, estado, err
}
