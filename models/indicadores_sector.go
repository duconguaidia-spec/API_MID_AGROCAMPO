package models

// IndicadorSector es un dato que se mide en el sector (por ejemplo "Precio promedio de leche")
type IndicadorSector struct {
	Id          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Unidad      string `json:"unidad"`
	Fuente      string `json:"fuente"`
}
