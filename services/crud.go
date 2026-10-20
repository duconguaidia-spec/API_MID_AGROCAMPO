package services

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
)

// RespuestaCrud es el formato en que responde el CRUD
type RespuestaCrud struct {
	Success bool            `json:"Success"`
	Status  int             `json:"Status"`
	Message string          `json:"Message"`
	Data    json.RawMessage `json:"Data"`
}

// Pedir le hace una petición al CRUD
// metodo: GET, POST, PUT o DELETE
// ruta: por ejemplo "/noticia" o "/noticia/3"
// cuerpo: lo que se manda en el JSON (puede ser nil)
func Pedir(metodo string, ruta string, cuerpo interface{}) (RespuestaCrud, error) {
	var respuesta RespuestaCrud

	urlCrud, err := beego.AppConfig.String("UrlCrud")
	if err != nil {
		return respuesta, err
	}

	if urlCrud == "" {
    return respuesta, fmt.Errorf("la configuración UrlCrud está vacía en app.conf")
}

	var datos []byte
	if cuerpo != nil {
		datos, err = json.Marshal(cuerpo)
		if err != nil {
			return respuesta, err
		}
	}

	peticion, err := http.NewRequest(metodo, urlCrud+ruta, bytes.NewBuffer(datos))
	if err != nil {
		return respuesta, err
	}
	peticion.Header.Set("Content-Type", "application/json")

	cliente := &http.Client{Timeout: 10 * time.Second}
	resp, err := cliente.Do(peticion)
	if err != nil {
		return respuesta, err
	}
	defer resp.Body.Close()

	texto, err := io.ReadAll(resp.Body)
	if err != nil {
		return respuesta, err
	}

	err = json.Unmarshal(texto, &respuesta)
	return respuesta, err
}
