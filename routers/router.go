package routers

import (
	"api_mid_AgroCampo/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	// Esta ruta reúne la información que necesita la página de detalle.
	beego.Router(
		"/api/v1/veterinarias/:id",
		&controllers.VeterinariaController{},
		"get:GetDetalleVeterinaria",
	)
}
