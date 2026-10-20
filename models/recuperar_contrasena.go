package models

// Lo que envía el cliente al MID
type RecuperarContrasenaRequest struct {
	Correo string `json:"Correo"`
}

// Lo que se lee de cada usuario devuelto por el CRUD (solo lo necesario)
type UsuarioCrudItem struct {
	Id     int    `json:"Id"`
	Correo string `json:"Correo"`
	Activo bool   `json:"Activo"`
}

// Respuesta final del MID
type RecuperarContrasenaResponse struct {
	Mensaje string `json:"Mensaje"`
}