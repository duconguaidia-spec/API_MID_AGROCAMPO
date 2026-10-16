package routers

import (
	"API_MID_AGROCAMPO/controllers"
	"API_MID_AGROCAMPO/controllers_InicioSesion"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//Inicio Sesion
	beego.Router("/api/creacion/usuario", &controllers_InicioSesion.CrearUsuarioController{}, "post:CrearUsuario")
	beego.Router("/api/recuperar_contrasena", &controllers_InicioSesion.Recuperar_contrasenaController{}, "post:Post")
	//Estadisticas_del_sector
	beego.Router("/v1/estadisticas_del_sector", &controllers.Estadisticas_del_sectorController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadistica_del_sector/:id([0-9]+)", &controllers.Estadisticas_del_sectorController{}, "get:GetOne;put:Put;delete:Delete")


	//Estadisticas_del_ganado(precio de la subasta)
	beego.Router("/v1/estadisticas_del_ganado", &controllers.Estadisticas_del_ganadoController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadisticas_del_ganado/:id([0-9]+", &controllers.Estadisticas_del_ganadoController{},"get:GetOne;put:Put;delete:Delete")

}
