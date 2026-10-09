package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type Estadisticas_del_sector struct {
	Id        int     `orm:"column(id); pk;auto" json:"id"`
	Fecha     string  `orm:"column(fecha)" json:"fecha"`
	Categoria string  `orm:"column(categoria)" json:"categoria"`
	Precio    float64 `orm:"column(precio)" json:"precio"`
	Activo    bool    `orm:"column(activo)" json:"activo"`
}

func (t *Estadisticas_del_sector) TableName() string {
	return "estadisticas_del_sector"
}

func init() {
	orm.RegisterModel(new(Estadisticas_del_sector))
}
