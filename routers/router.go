package routers

import (
	controllers "API_MID_AGROCAMPO/controllers/gestion-agropecuaria"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	//Estadisticas_del_sector
	beego.Router("/v1/estadisticas_del_sector", &controllers.Estadisticas_del_sectorController{}, "get:GetAll;post:Post")
	beego.Router("/v1/estadistica_sector/:id([0-9]+)", &controllers.Estadisticas_del_sectorController{}, "get:GetOne;put:Put;delete:Delete")
}


  //  beego.Router("/", &controllers.MainController{})
//}
