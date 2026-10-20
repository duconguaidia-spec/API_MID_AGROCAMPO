package models

// Subasta es lo que el MID recibe y devuelve (la fecha va como texto AAAA-MM-DD)
type Subasta struct {
	Id            int    `json:"id"`
	RegistradaPor int    `json:"registrada_por"`
	FechaSubasta  string `json:"fecha_subasta,omitempty"`
	Ubicacion     string `json:"ubicacion"`
	Nombre        string `json:"nombre"`
}
