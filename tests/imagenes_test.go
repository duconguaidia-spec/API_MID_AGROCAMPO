package test

import (
	"bytes"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImagenProducto(t *testing.T) {
	t.Chdir(t.TempDir())
	var pngData bytes.Buffer
	png.Encode(&pngData, image.NewRGBA(image.Rect(0, 0, 2, 2)))
	for _, caso := range []struct {
		nombre string
		datos  []byte
		estado int
	}{
		{"producto.png", pngData.Bytes(), 201},
		{"falso.png", []byte("no es una imagen"), 400},
		{"grande.png", make([]byte, 2*1024*1024+1), 400},
	} {
		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, _ := writer.CreateFormFile("imagen", caso.nombre)
		part.Write(caso.datos)
		writer.Close()
		req := httptest.NewRequest("POST", "/api/v1/imagenes/productos", &body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(res, req)
		if res.Code != caso.estado {
			t.Fatalf("%s: HTTP %d: %s", caso.nombre, res.Code, res.Body.String())
		}
		if res.Code == 201 {
			var respuesta struct {
				Data struct {
					URL string `json:"url"`
				} `json:"data"`
			}
			json.Unmarshal(res.Body.Bytes(), &respuesta)
			imagen := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(imagen, httptest.NewRequest("GET", respuesta.Data.URL, nil))
			if imagen.Code != 200 || !bytes.Equal(imagen.Body.Bytes(), pngData.Bytes()) {
				t.Fatal("La imagen guardada no se puede recuperar")
			}
		}
	}
}

func TestEliminarVeterinaria(t *testing.T) {
	eliminado := false
	crudPrueba(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/veterinaria/7" {
			t.Errorf("ruta inesperada: %s", r.URL.Path)
		}
		if r.Method == "DELETE" {
			eliminado = true
			w.Write([]byte(`{"success":true}`))
			return
		}
		w.Write([]byte(`{"success":true,"data":{"id":7}}`))
	})
	pedir(t, "DELETE", "/api/v1/veterinarias/0", "", 400)
	pedir(t, "DELETE", "/api/v1/veterinarias/7", "", 200)
	if !eliminado {
		t.Fatal("No se llamó al DELETE del CRUD")
	}
}
