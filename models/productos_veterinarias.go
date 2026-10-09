package models

// VeterinariaCRUD representa la veterinaria anidada que devuelve la CRUD.
type VeterinariaCRUD struct {
	ID int `json:"id"`
}

// ProductoCRUD representa el formato de producto que devuelve la CRUD.
type ProductoCRUD struct {
	ID          int             `json:"id"`
	Nombre      string          `json:"nombre"`
	Descripcion string          `json:"descripcion"`
	ImagenURL   string          `json:"imagen_url"`
	Precio      float64         `json:"precio"`
	Activo      bool            `json:"activo"`
	Veterinaria VeterinariaCRUD `json:"veterinaria"`
}

// RespuestaProductosCRUD representa la respuesta de la CRUD.
type RespuestaProductosCRUD struct {
	Success bool           `json:"success"`
	Data    []ProductoCRUD `json:"data"`
}

// ProductoVeterinaria es el formato que el MID devuelve al frontend.
type ProductoVeterinaria struct {
	ID            int     `json:"id"`
	VeterinariaID int     `json:"veterinaria_id"`
	Nombre        string  `json:"nombre"`
	Descripcion   string  `json:"descripcion,omitempty"`
	ImagenURL     string  `json:"imagen_url,omitempty"`
	Precio        float64 `json:"precio"`
	Activo        bool    `json:"activo"`
}