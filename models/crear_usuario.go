package models

type CrearUsuarioRequest struct {
	NombreCompleto       string `json:"NombreCompleto"`
	Correo               string `json:"Correo"`
	Telefono             string `json:"Telefono"`
	IdRol                int    `json:"IdRol"`
	VerificacionDosPasos bool   `json:"VerificacionDosPasos"`
	Avatar               string `json:"Avatar"`
	Activo               bool   `json:"Activo"`
}


type RolRef struct {
	Id int `json:"Id"`
}

type UsuarioCrudRequest struct {
	NombreCompleto       string `json:"NombreCompleto"`
	Correo               string `json:"Correo"`
	Telefono             string `json:"Telefono"`
	IdRol                string `json:"IdRol"`
	VerificacionDosPasos bool   `json:"VerificacionDosPasos"`
	Avatar               string `json:"Avatar"`
	Activo               bool   `json:"Activo"`
}

// Respuesta del CRUD al crear (solo necesita el Id)
type UsuarioCrudResponse struct {
	Id int `json:"Id"`
}

type CrearUsuarioResponse struct {
	Mensaje   string `json:"Mensaje"`
	IDUsuario int    `json:"IDUsuario"`
}