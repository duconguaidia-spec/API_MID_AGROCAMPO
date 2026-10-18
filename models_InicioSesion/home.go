package models_InicioSesion

// Lo que se lee del CRUD de solicitudes (solo lo necesario)
type SolicitudCrud struct {
	Id     int    `json:"Id"`
	Estado string `json:"Estado"`
}

type PerfilHome struct {
	IDUsuario      int    `json:"IDUsuario"`
	NombreCompleto string `json:"NombreCompleto"`
	Correo         string `json:"Correo"`
	Telefono       string `json:"Telefono"`
	IdRol          string `json:"IdRol"`
	Avatar         string `json:"Avatar"`
	Ubicacion      string `json:"Ubicacion"`
}


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
	Productores   *int `json:"Productores"`
	Veterinarias  *int `json:"Veterinarias"`
	Eficiencia    *int `json:"Eficiencia"`
}


// Respuesta final del MID
type HomeResponse struct {
	Perfil     PerfilHome     `json:"Perfil"`
	Contadores ContadoresHome `json:"Contadores"`
	Noticias   []NoticiaCrud  `json:"Noticias"`
}