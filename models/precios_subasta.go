package models

// PrecioSubasta es el precio por kilo de una categoria de ganado en una subasta
type PrecioSubasta struct {
	Id                int     `json:"id"`
	SubastaId         int     `json:"subasta_id"`
	CategoriaGanadoId int     `json:"categoria_ganado_id"`
	PrecioPorKg       float64 `json:"precio_por_kg"`
	Unidad            string  `json:"unidad"`
}
