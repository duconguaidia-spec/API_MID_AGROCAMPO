package models

// MedicionSector es el valor de un indicador en una fecha y una region
// (la fecha va como texto AAAA-MM-DD)
type MedicionSector struct {
	Id          int     `json:"id"`
	IndicadorId int     `json:"indicador_id"`
	Fecha       string  `json:"fecha,omitempty"`
	Region      string  `json:"region"`
	Valor       float64 `json:"valor"`
	Fuente      string  `json:"fuente"`
}
