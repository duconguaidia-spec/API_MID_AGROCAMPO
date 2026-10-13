package routers

import (
	"API_MID_AGROCAMPO/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	beego.Router("/api/creacion/usuario", &controllers.CrearUsuarioController{}, "post:CrearUsuario")
	beego.Router("/api/recuperar_contrasena", &controllers.Recuperar_contrasenaController{}, "post:Post")
	//Estadisticas_del_sector
	beego.Router("/v1/estadisticas_del_sector", &controllers.Estadisticas_del_sectorController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadistica_sector/:id([0-9]+)", &controllers.Estadisticas_del_sectorController{}, "get:GetOne;put:Put;delete:Delete")
}
