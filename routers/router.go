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

	// Recuperacion de contraseña
	beego.Router(
		"/api/recuperar_contrasena",
		&controllers.Recuperar_contrasenaController{},
		"post:Solicitar",
	)
	beego.Router(
		"/api/restablecer_contrasena",
		&controllers.Recuperar_contrasenaController{},
		"post:Restablecer",
	)

	// Home: contadores y ultimas noticias (publico, sin token)
	beego.Router(
		"/api/v1/home",
		&controllers.HomeController{},
		"get:GetAll",
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
	beego.Router("/v1/estadisticas_del_ganado", &controllers.Estadisticas_del_ganadoController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadisticas_del_ganado/:id([0-9]+)", &controllers.Estadisticas_del_ganadoController{}, "get:GetOne;put:Put;delete:Delete")

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
}