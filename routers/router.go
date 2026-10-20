package routers

import (
	"API_MID_AGROCAMPO/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//Inicio Sesion
	beego.Router(
		"/api/creacion/usuario", 
		&controllers.CrearUsuarioController{}, 
		"post:CrearUsuario",
	)
	beego.Router(
		"/api/recuperar_contrasena", 
		&controllers.Recuperar_contrasenaController{}, 
		"post:Post",
	)
	//Estadisticas_del_sector
	beego.Router(
		"/v1/estadisticas_del_sector",
		&controllers.Estadisticas_del_sectorController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/estadisticas_del_sector/:id([0-9]+)",
		&controllers.Estadisticas_del_sectorController{},
		"get:GetOne;put:Put;delete:Delete",
	)

	//Estadisticas_del_ganado(precio de la subasta)
	beego.Router(
		"/v1/estadisticas_del_ganado",
		&controllers.Estadisticas_del_ganadoController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/estadisticas_del_ganado/:id([0-9]+)",
		&controllers.Estadisticas_del_ganadoController{},
		"get:GetOne;put:Put;delete:Delete",
	)

	// Esta ruta reúne la información que necesita la página de detalle.
	beego.Router(
		"/api/v1/veterinarias/:id",
		&controllers.VeterinariaController{},
		"get:GetDetalleVeterinaria",
	)
	beego.Router(
		"/api/v1/veterinarias",
		&controllers.VeterinariaController{},
		"post:PostCrearVeterinaria",
	)

	// Subastas
	beego.Router(
		"/v1/subasta",
		&controllers.SubastasController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/subasta/:id([0-9]+)",
		&controllers.SubastasController{},
		"get:GetOne;put:Put;delete:Delete",
	)

	// Precios de las subastas
		beego.Router(
		"/v1/precio_subasta",
		&controllers.PrecioSubastaController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/precio_subasta/:id([0-9]+)",
		&controllers.PrecioSubastaController{},
		"get:GetOne;put:Put;delete:Delete",
	)

	// Indicadores del sector
	beego.Router(
		"/v1/indicador_sector",
		&controllers.IndicadoresSectorController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/indicador_sector/:id([0-9]+)",
		&controllers.IndicadoresSectorController{},
		"get:GetOne;put:Put;delete:Delete",
	)

	// Mediciones del sector
	beego.Router(
		"/v1/medicion_sector",
		&controllers.MedicionesSectorController{},
		"get:GetAll;post:Post",
	)
	beego.Router(
		"/v1/medicion_sector/:id([0-9]+)",
		&controllers.MedicionesSectorController{},
		"get:GetOne;put:Put;delete:Delete",
	)
}
