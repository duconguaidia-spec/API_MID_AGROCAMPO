package models

// VeterinariaCRUD representa los campos que devuelve la CRUD para una clínica.
type VeterinariaCRUD struct {
	ID              int      `json:"id"`
	Nombre          string   `json:"nombre"`
	Descripcion     string   `json:"descripcion"`
	Direccion       string   `json:"direccion"`
	Ciudad          string   `json:"ciudad"`
	Departamento    string   `json:"departamento"`
	Telefono        string   `json:"telefono"`
	CorreoPublico   string   `json:"correo_publico"`
	HorarioAtencion string   `json:"horario_atencion"`
	Latitud         *float64 `json:"latitud"`
	Longitud        *float64 `json:"longitud"`
	Aprobado        bool     `json:"aprobado"`
	Activo          bool     `json:"activo"`
}

// RespuestaVeterinariaCRUD es el sobre JSON que entrega la CRUD al consultar una clínica.
type RespuestaVeterinariaCRUD struct {
	Success bool            `json:"success"`
	Data    VeterinariaCRUD `json:"data"`
}

// VeterinariaServicioCRUD representa un servicio del catálogo asociado a la clínica.
type VeterinariaServicioCRUD struct {
	Nombre string `json:"nombre"`
	Activo bool   `json:"activo"`
}

// RelacionServicioCRUD representa la respuesta de la tabla intermedia de servicios.
type RelacionServicioCRUD struct {
	Activo          bool                    `json:"activo"`
	ServicioGeneral VeterinariaServicioCRUD `json:"servicio_general"`
}

// RespuestaRelacionesServiciosCRUD es el sobre JSON que entrega la CRUD para la relación.
type RespuestaRelacionesServiciosCRUD struct {
	Success bool                   `json:"success"`
	Data    []RelacionServicioCRUD `json:"data"`
}

// ProductoCRUD representa un producto tal como lo devuelve la CRUD.
type ProductoCRUD struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	ImagenURL   string  `json:"imagen_url"`
	Precio      float64 `json:"precio"`
	Activo      bool    `json:"activo"`
	Veterinaria struct {
		ID int `json:"id"`
	} `json:"veterinaria"`
}

// RespuestaProductosCRUD es el sobre JSON que entrega la CRUD al listar productos.
type RespuestaProductosCRUD struct {
	Success bool           `json:"success"`
	Data    []ProductoCRUD `json:"data"`
}

// ServicioCRUD representa un servicio propio que ofrece una veterinaria.
type ServicioCRUD struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	Activo      bool    `json:"activo"`
}

// RespuestaServiciosCRUD es el sobre JSON que entrega la CRUD al listar servicios propios.
type RespuestaServiciosCRUD struct {
	Success bool           `json:"success"`
	Data    []ServicioCRUD `json:"data"`
}

// ProductoVeterinaria es la versión simplificada de un producto para el frontend.
type ProductoVeterinaria struct {
	ID          int     `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion,omitempty"`
	ImagenURL   string  `json:"imagen_url,omitempty"`
	Precio      float64 `json:"precio"`
}

// ServicioVeterinaria es la versión simplificada de un servicio para el frontend.
type ServicioVeterinaria struct {
	ID          int     `json:"id,omitempty"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion,omitempty"`
	Precio      float64 `json:"precio,omitempty"`
}

// DetalleVeterinaria reúne en una sola respuesta los datos relacionados que muestra la página.
type DetalleVeterinaria struct {
	Veterinaria VeterinariaCRUD       `json:"veterinaria"`
	Servicios   []ServicioVeterinaria `json:"servicios"`
	Productos   []ProductoVeterinaria `json:"productos"`
}
// Datos que recibe el MID desde el frontend.
type CrearVeterinariaRequest struct {
	PropietarioID   *int                           `json:"propietario_id,omitempty"`
	Nombre          string                         `json:"nombre"`
	Descripcion     string                         `json:"descripcion,omitempty"`
	Direccion       string                         `json:"direccion"`
	Ciudad          string                         `json:"ciudad"`
	Departamento    string                         `json:"departamento,omitempty"`
	Telefono        string                         `json:"telefono"`
	CorreoPublico   string                         `json:"correo_publico,omitempty"`
	HorarioAtencion string                         `json:"horario_atencion,omitempty"`
	Latitud         *float64                       `json:"latitud,omitempty"`
	Longitud        *float64                       `json:"longitud,omitempty"`
	Productos       []CrearProductoRequest         `json:"productos"`
}

// Datos de cada producto que recibe el MID.
type CrearProductoRequest struct {
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion,omitempty"`
	ImagenURL   string  `json:"imagen_url,omitempty"`
	Precio      float64 `json:"precio"`
}

// Formato que el MID envía a la CRUD de veterinarias.
type CrearVeterinariaCRUD struct {
	PropietarioID   *int     `json:"propietario_id,omitempty"`
	Nombre          string   `json:"nombre"`
	Descripcion     string   `json:"descripcion,omitempty"`
	Direccion       string   `json:"direccion"`
	Ciudad          string   `json:"ciudad"`
	Departamento    string   `json:"departamento,omitempty"`
	Telefono        string   `json:"telefono"`
	CorreoPublico   string   `json:"correo_publico,omitempty"`
	HorarioAtencion string   `json:"horario_atencion,omitempty"`
	Latitud         *float64 `json:"latitud,omitempty"`
	Longitud        *float64 `json:"longitud,omitempty"`
	Aprobado        bool     `json:"aprobado"`
	Activo          bool     `json:"activo"`
}

// Referencia a la veterinaria que se anida en el JSON de la CRUD de productos.
type ReferenciaVeterinariaCRUD struct {
	ID int `json:"id"`
}

// Formato que el MID envía a la CRUD de productos.
type CrearProductoVeterinariaCRUD struct {
	Veterinaria ReferenciaVeterinariaCRUD `json:"veterinaria"`
	Nombre      string                    `json:"nombre"`
	Descripcion string                    `json:"descripcion,omitempty"`
	ImagenURL   string                    `json:"imagen_url,omitempty"`
	Precio      float64                   `json:"precio"`
	Activo      bool                      `json:"activo"`
}

// Respuesta de la CRUD al crear un solo producto.
type RespuestaCrearProductoCRUD struct {
	Success bool         `json:"success"`
	Data    ProductoCRUD `json:"data"`
}

type CrearVeterinariaDataResponse struct {
	Veterinaria VeterinariaCRUD       `json:"veterinaria"`
	Productos   []ProductoCRUD        `json:"productos"`
}

type CrearVeterinariaResponse struct {
	Success bool                          `json:"success"`
	Message string                        `json:"message"`
	Data    CrearVeterinariaDataResponse `json:"data"`
}
