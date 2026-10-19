package models

// Lo que se lee del CRUD de temas de foro (solo lo necesario)
type TemaForoCrud struct {
	Id     int  `json:"id"`
	Activo bool `json:"activo"`
}

// Lo que se lee del CRUD de noticias (solo lo necesario)
type NoticiaCrud struct {
	Id          int     `json:"id"`
	Titulo      string  `json:"titulo"`
	Resumen     string  `json:"resumen"`
	Fuente      string  `json:"fuente"`
	Imagen      string  `json:"imagen_url"`
	PublicadaEn *string `json:"publicada_en"`
}

// Lo que se lee del CRUD de veterinarias (solo lo necesario)
type VeterinariaCrud struct {
	Id       int  `json:"id"`
	Aprobado bool `json:"aprobado"`
	Activo   bool `json:"activo"`
}

// Lo que se lee del CRUD de solicitudes de reporte (solo lo necesario)
type SolicitudCrud struct {
	Id           int     `json:"id"`
	Error        string  `json:"error"`
	Activo       bool    `json:"activo"`
	CompletadoEn *string `json:"completado_en"`
}

// Noticia que se devuelve al cliente
type NoticiaHome struct {
	Id      int    `json:"Id"`
	Titulo  string `json:"Titulo"`
	Resumen string `json:"Resumen"`
	Fuente  string `json:"Fuente"`
	Imagen  string `json:"Imagen"`
	Fecha   string `json:"Fecha"`
}

// Un contador en null significa "no disponible"
type ContadoresHome struct {
	ForosNoticias *int `json:"ForosNoticias"`
	Veterinarias  *int `json:"Veterinarias"`
	Eficiencia    *int `json:"Eficiencia"`
}

// Respuesta final del MID
type HomeResponse struct {
	Contadores ContadoresHome `json:"Contadores"`
	Noticias   []NoticiaHome  `json:"Noticias"`
}