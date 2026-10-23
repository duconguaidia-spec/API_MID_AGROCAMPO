package models

// Lo que envia el cliente al MID
type CrearUsuarioRequest struct {
	Nombre     string `json:"Nombre"`
	Apellido   string `json:"Apellido"`
	Correo     string `json:"Correo"`
	Telefono   string `json:"Telefono"`
	IdRol      int    `json:"IdRol"`
	Contrasena string `json:"Contrasena"`
	Avatar     string `json:"Avatar"`
}

// Referencia a otra tabla: el rel(fk) de Beego espera un json {"id": n}
type RefCrud struct {
	Id int `json:"id"`
}

// Lo que se envia al CRUD de usuarios
type UsuarioCrudRequest struct {
	Nombre     string  `json:"nombre"`
	Apellido   string  `json:"apellido"`
	Correo     string  `json:"correo"`
	Telefono   string  `json:"telefono"`
	IdRol      RefCrud `json:"rol"`
	Avatar     string  `json:"avatar_url"`
	Verificado bool    `json:"verificado"`
	Activo     bool    `json:"activo"`
}

// Datos del usuario creado que devuelve el CRUD (solo necesita el id)
type DatosUsuarioCrud struct {
	Id int `json:"id"`
}

// Respuesta del CRUD de usuarios al crear: el id viene dentro de "data"
type UsuarioCrudResponse struct {
	Data DatosUsuarioCrud `json:"data"`
}

// Parte comun de las respuestas del CRUD, sirve para saber si la operacion fue exitosa
type RespuestaCrud struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Lo que se envia al CRUD de contrasenas
type ContrasenaCrudRequest struct {
	IdUsuario      RefCrud `json:"IdUsuario"`
	ContrasenaHash string  `json:"ContrasenaHash"`
	Activa         bool    `json:"Activa"`
}

type CrearUsuarioResponse struct {
	Mensaje                 string `json:"Mensaje"`
	IDUsuario               int    `json:"IDUsuario"`
	CorreoBienvenidaEnviado bool   `json:"CorreoBienvenidaEnviado"`
}