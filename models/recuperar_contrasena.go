package models

// SolicitudRecuperacion es el cuerpo esperado en POST /api/recuperar_contrasena
type SolicitudRecuperacion struct {
	Correo string `json:"correo"`
}

// SolicitudRestablecer es el cuerpo esperado en POST /api/restablecer_contrasena
type SolicitudRestablecer struct {
	Correo          string `json:"correo"`
	Codigo          string `json:"codigo"`
	NuevaContrasena string `json:"nueva_contrasena"`
}