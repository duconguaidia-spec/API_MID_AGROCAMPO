package routers

import (
	"API_MID_AGROCAMPO/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	beego.Router("/api/creacion/usuario", &controllers.CrearUsuarioController{}, "post:CrearUsuario")
	beego.Router("/api/recuperar_contrasena", &controllers.Recuperar_contrasenaController{}, "post:Post")
}