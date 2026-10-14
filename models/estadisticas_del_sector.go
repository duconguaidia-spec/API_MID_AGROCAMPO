package models

import (
	"github.com/beego/beego/v2/client/orm"

	
)

type Estadisticas_del_sector struct {
	Id       int     `orm:"column(id); pk;auto" json:"id"`
	Tipo     string  `orm:"column(tipo)"json:"tipo"`
	Etiqueta string  `orm:"column(etiqueta)"json:"etiqueta"`
	Valor    float64 `orm:"column(valor)"json:"valor"`
	Unidad   string  `orm:"column(unidad)"json:"unidad"`
	Año      int     `orm:"column(año)"json:"año"`
	Activo   bool    `orm:"column(activo)"json:"activo"`
}

func (t *Estadisticas_del_sector) TableName() string {
	return "estadisticas_del_sector"
}

func init() {
	orm.RegisterModel(new(Estadisticas_del_sector))
}
