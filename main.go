package main

import (
	_ "API_MID_AGROCAMPO/routers"
	"errors"
	"os"
	"path/filepath"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	configurarClaveUsuarios()
	beego.Run()
}

func configurarClaveUsuarios() {
	archivo := filepath.Join("..", "AgroCampo_Apis_Crud_Beego", ".env")
	datos, err := godotenv.Read(archivo)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// Los errores del lector pueden incluir líneas del archivo con credenciales.
		logs.Warn("No se pudo leer el .env compartido para configurar Usuarios")
	}
	if clave, existe := datos["AGROCAMPO_CLAVE_INTERNA"]; existe {
		beego.AppConfig.Set("clave_crud_usuarios", clave)
	}
	if len(beego.AppConfig.DefaultString("clave_crud_usuarios", "")) < 32 {
		logs.Warn("Usuarios requiere AGROCAMPO_CLAVE_INTERNA de al menos 32 caracteres")
	}
}
