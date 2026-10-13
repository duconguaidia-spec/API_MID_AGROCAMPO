package models

// Lo que se lee del CRUD de usuarios
type UsuarioPerfilCrud struct {
	Id                   int    `json:"id,string"`
	NombreCompleto       string `json:"NombreCompleto"`
	Correo               string `json:"Correo"`
	Telefono             string `json:"Telefono"`
	IdRol                string `json:"IdRol"`
	VerificacionDosPasos bool   `json:"VerificacionDosPasos"`
	Avatar               string `json:"Avatar"`
	Activo               bool   `json:"Activo"`
}

// Lo que se lee del CRUD de perfiles (PerfilExtendido)
type PerfilCrud struct {
	UsuarioId int    `json:"usuario_id"`
	Ubicacion string `json:"ubicacion"`
	Activo    bool   `json:"activo"`
}

// Respuesta final del MID
type PerfilUsuarioResponse struct {
	IDUsuario            int    `json:"IDUsuario"`
	NombreCompleto       string `json:"NombreCompleto"`
	Correo               string `json:"Correo"`
	Telefono             string `json:"Telefono"`
	IdRol                string `json:"IdRol"`
	VerificacionDosPasos bool   `json:"VerificacionDosPasos"`
	Avatar               string `json:"Avatar"`
	Ubicacion            string `json:"Ubicacion"`
	Activo               bool   `json:"Activo"`
}