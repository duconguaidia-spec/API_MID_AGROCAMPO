package controllers

import (
	"API_MID_AGROCAMPO/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

// HomeController operations for Home
type HomeController struct {
	beego.Controller
}

func (c *HomeController) responderError(status int, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]string{"Error": mensaje}
	c.ServeJSON()
}

func (c *HomeController) obtener(url string, destino interface{}) (int, error) {
	peticion, err := http.NewRequestWithContext(c.Ctx.Request.Context(), http.MethodGet, url, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	respuesta, err := clienteModulos.Do(peticion)
	if err != nil {
		logs.Warn("Home: no fue posible consultar %s: %v", url, err)
		return http.StatusInternalServerError, err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode < http.StatusOK || respuesta.StatusCode >= http.StatusMultipleChoices {
		logs.Warn("Home: %s respondio con el estado %d", url, respuesta.StatusCode)
		return respuesta.StatusCode, errors.New("respuesta no exitosa del api")
	}

	cuerpo, err := io.ReadAll(respuesta.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Algunos CRUD devuelven la lista dentro de "data" y otros la devuelven directa
	var envoltura struct {
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(cuerpo, &envoltura) == nil && len(envoltura.Data) > 0 {
		cuerpo = envoltura.Data
	}

	err = json.Unmarshal(cuerpo, destino)
	if err != nil {
		logs.Warn("Home: la respuesta de %s no tiene el formato esperado: %v", url, err)
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (c *HomeController) contarForosNoticias(urlForos string, urlNoticias string) *int {
	var foros []models.TemaForoCrud
	_, err := c.obtener(urlForos+"?limit=-1", &foros)
	if err != nil {
		return nil
	}

	var noticias []models.NoticiaCrud
	_, err = c.obtener(urlNoticias+"?limit=-1&query=Estado:publicada,AccesoLimitado:false", &noticias)
	if err != nil {
		return nil
	}

	total := 0
	for _, foro := range foros {
		if foro.Activo {
			total++
		}
	}
	for _, noticia := range noticias {
		if noticia.PublicadaEn != nil {
			total++
		}
	}

	return &total
}

func (c *HomeController) contarVeterinarias(urlVeterinarias string) *int {
	var veterinarias []models.VeterinariaCrud
	_, err := c.obtener(urlVeterinarias+"?limit=-1", &veterinarias)
	if err != nil {
		return nil
	}

	total := 0
	for _, veterinaria := range veterinarias {
		if veterinaria.Activo && veterinaria.Aprobado {
			total++
		}
	}

	return &total
}

func (c *HomeController) calcularEficiencia(urlSolicitudes string) *int {
	var solicitudes []models.SolicitudCrud
	_, err := c.obtener(urlSolicitudes+"?limit=-1", &solicitudes)
	if err != nil {
		return nil
	}

	total := 0
	resueltas := 0
	for _, solicitud := range solicitudes {
		if !solicitud.Activo {
			continue
		}

		total++
		if solicitud.CompletadoEn != nil && solicitud.Error == "" {
			resueltas++
		}
	}

	eficiencia := 0
	if total > 0 {
		eficiencia = resueltas * 100 / total
	}

	return &eficiencia
}

func (c *HomeController) ultimasNoticias(urlNoticias string) []models.NoticiaHome {
	ultimas := []models.NoticiaHome{}

	var noticias []models.NoticiaCrud
	_, err := c.obtener(urlNoticias+"?sortby=FechaNoticia&order=desc&limit=5&query=Estado:publicada,AccesoLimitado:false", &noticias)
	if err != nil {
		return ultimas
	}

	for _, noticia := range noticias {
		if noticia.PublicadaEn == nil {
			continue
		}

		ultimas = append(ultimas, models.NoticiaHome{
			Id:      noticia.Id,
			Titulo:  noticia.Titulo,
			Resumen: noticia.Resumen,
			Fuente:  noticia.Fuente,
			Imagen:  noticia.Imagen,
			Fecha:   *noticia.PublicadaEn,
		})

		if len(ultimas) == 5 {
			break
		}
	}

	return ultimas
}

// GetAll ...
// @Title GetAll
// @Description get Home
// @Success 200 {object} models.HomeResponse
// @Failure 500 {object} map[string]string
// @router / [get]
func (c *HomeController) contarProductores() *int {
	respuesta, status, err := peticionModulo(c.Ctx.Request, "crud_usuarios", "cuenta/contar-productores", http.MethodPost, []byte(`{}`))
	var total int
	if err != nil || status != http.StatusOK || !respuesta.Success || json.Unmarshal(respuesta.Data, &total) != nil {
		return nil
	}
	return &total
}

func (c *HomeController) GetAll() {
	urlForos, _ := beego.AppConfig.String("url_foros")
	if urlForos == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_foros no esta configurada en el app.conf")
		return
	}

	urlNoticias, _ := beego.AppConfig.String("url_noticias")
	if urlNoticias == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_noticias no esta configurada en el app.conf")
		return
	}

	urlVeterinarias, _ := beego.AppConfig.String("url_veterinarias")
	if urlVeterinarias == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_veterinarias no esta configurada en el app.conf")
		return
	}

	urlSolicitudes, _ := beego.AppConfig.String("url_solicitudes")
	if urlSolicitudes == "" {
		c.responderError(http.StatusInternalServerError, "La variable url_solicitudes no esta configurada en el app.conf")
		return
	}

	var resumen models.HomeResponse
	var tareas sync.WaitGroup
	tareas.Go(func() { resumen.Contadores.ForosNoticias = c.contarForosNoticias(urlForos, urlNoticias) })
	tareas.Go(func() { resumen.Contadores.Veterinarias = c.contarVeterinarias(urlVeterinarias) })
	tareas.Go(func() { resumen.Contadores.Productores = c.contarProductores() })
	tareas.Go(func() { resumen.Contadores.Eficiencia = c.calcularEficiencia(urlSolicitudes) })
	tareas.Go(func() { resumen.Noticias = c.ultimasNoticias(urlNoticias) })
	tareas.Wait()
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = resumen
	c.ServeJSON()
}
