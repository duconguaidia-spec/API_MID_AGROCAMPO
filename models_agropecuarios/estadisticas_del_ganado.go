package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type EstadisticaGanado struct {
	Id        int     `orm:"column(id);pk;auto" json:"id"`
	Fecha     string  `orm:"column(fecha)" json:"fecha"`
	Categoria string  `orm:"column(categoria)" json:"categoria"`
	Precio    float64 `orm:"column(precio)" json:"precio"`
	Activo    bool    `orm:"column(activo)" json:"activo"`
}

func (t *EstadisticaGanado) TableName() string {
	return "estadistica_ganado"
}

func init() {
	orm.RegisterModel(new(EstadisticaGanado))
}
