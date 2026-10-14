package routers

import (
	"API_MID_AGROCAMPO/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//Estadisticas_del_sector
	beego.Router("/v1/estadisticas_del_sector", &controllers.Estadisticas_del_sectorController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadistica_del_sector/:id([0-9]+)", &controllers.Estadisticas_del_sectorController{}, "get:GetOne;put:Put;delete:Delete")


	//Estadisticas_del_ganado(precio de la subasta)
	beego.Router("/v1/estadisticas_del_ganado", &controllers.Estadisticas_del_ganadoController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadisticas_del_ganado/:id([0-9]+", &controllers.Estadisticas_del_ganadoController{},"get:GetOne;put:Put;delete:Delete")

}


  //  beego.Router("/", &controllers.MainController{})
//}
