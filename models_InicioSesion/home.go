package models_InicioSesion

// Claims del token del login (supuesto: JWT HS256)
type TokenClaims struct {
	Id  int   `json:"id"`
	Exp int64 `json:"exp"`
}

// Lo que se lee del CRUD de solicitudes (solo lo necesario)
type SolicitudCrud struct {
	Id     int    `json:"Id"`
	Estado string `json:"Estado"`
}

// Lo que se lee del CRUD de noticias/eventos y se devuelve al cliente
type NoticiaCrud struct {
	Id          int    `json:"Id"`
	Titulo      string `json:"Titulo"`
	Descripcion string `json:"Descripcion"`
	Fecha       string `json:"Fecha"`
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

// Un contador en null significa "no disponible"
type ContadoresHome struct {
	Foros        *int `json:"Foros"`
	Productores  *int `json:"Productores"`
	Veterinarias *int `json:"Veterinarias"`
	Eficiencia   *int `json:"Eficiencia"`
}

// Respuesta final del MID
type HomeResponse struct {
	Perfil     PerfilHome     `json:"Perfil"`
	Contadores ContadoresHome `json:"Contadores"`
	Noticias   []NoticiaCrud  `json:"Noticias"`
}