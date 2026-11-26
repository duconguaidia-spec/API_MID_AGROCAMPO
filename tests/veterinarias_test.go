package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	_ "API_MID_AGROCAMPO/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	beego.TestBeegoInit(filepath.Dir(filepath.Dir(file)))
}

func pedir(t *testing.T, metodo, ruta, cuerpo string, codigo int, cookies ...*http.Cookie) map[string]interface{} {
	t.Helper()
	r := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	r.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != codigo {
		t.Fatalf("%s %s: esperaba %d, recibió %d: %s", metodo, ruta, codigo, w.Code, w.Body.String())
	}
	var respuesta map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respuesta); err != nil {
		t.Fatalf("respuesta inválida: %s", w.Body.String())
	}
	if respuesta["success"] != (codigo < 400) {
		t.Fatalf("success incoherente: %s", w.Body.String())
	}
	return respuesta
}

func crudPrueba(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	servidor := httptest.NewServer(handler)
	t.Cleanup(servidor.Close)
	for clave, ruta := range map[string]string{
		"url_veterinaria":                      "/v1/veterinaria",
		"url_productos_veterinaria":            "/v1/producto_veterinaria",
		"url_servicios_veterinaria":            "/v1/servicioveterinaria",
		"url_relaciones_servicios_veterinaria": "/v1/tr_veterinariaservicio",
	} {
		anterior, _ := beego.AppConfig.String(clave)
		beego.AppConfig.Set(clave, servidor.URL+ruta)
		t.Cleanup(func() { beego.AppConfig.Set(clave, anterior) })
	}
}

func TestLecturaYCreacionVeterinaria(t *testing.T) {
	crudPrueba(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method + " " + strings.TrimRight(r.URL.Path, "/") {
		case "GET /v1/veterinaria":
			w.Write([]byte(`{"success":true,"data":[{"id":7,"nombre":"Clínica","activo":true},{"id":8,"activo":false}]}`))
		case "GET /v1/veterinaria/7":
			w.Write([]byte(`{"success":true,"data":{"id":7,"nombre":"Clínica","activo":true}}`))
		case "GET /v1/producto_veterinaria":
			if r.URL.Query().Get("query") != "IdVeterinaria__Id:7" {
				t.Error("productos sin filtro de clínica")
			}
			w.Write([]byte(`{"success":true,"data":[{"id":3,"nombre":"Producto","activo":true}]}`))
		case "GET /v1/servicioveterinaria", "GET /v1/tr_veterinariaservicio":
			if r.URL.Query().Get("query") != "IdVeterinaria__Id:7" {
				t.Error("servicios sin filtro de clínica")
			}
			w.Write([]byte(`{"success":true,"data":[]}`))
		case "POST /v1/veterinaria":
			if !strings.HasSuffix(r.URL.Path, "/") {
				t.Error("POST debe usar la ruta de colección del CRUD")
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success":true,"data":{"id":7,"nombre":"Nueva","activo":true}}`))
		case "POST /v1/producto_veterinaria", "POST /v1/servicioveterinaria":
			var cuerpo struct {
				Veterinaria struct {
					ID int `json:"id"`
				} `json:"veterinaria"`
			}
			json.NewDecoder(r.Body).Decode(&cuerpo)
			if cuerpo.Veterinaria.ID != 7 {
				t.Error("relación sin ID generado por el CRUD")
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success":true,"data":{"id":3,"nombre":"Producto","activo":true}}`))
		default:
			t.Errorf("petición inesperada: %s %s", r.Method, r.URL)
			w.WriteHeader(500)
		}
	})
	lista := pedir(t, "GET", "/api/v1/veterinarias", "", 200)
	if len(lista["data"].([]interface{})) != 1 {
		t.Fatal("se deben excluir las clínicas inactivas")
	}
	detalle := pedir(t, "GET", "/api/v1/veterinarias/7", "", 200)["data"].(map[string]interface{})
	if len(detalle["productos"].([]interface{})) != 1 {
		t.Fatal("faltan los productos del detalle")
	}
	pedir(t, "POST", "/api/v1/veterinarias", `{"nombre":"Nueva","direccion":"Calle 1","ciudad":"Yopal","telefono":"123","productos":[{"nombre":"Producto","precio":20}],"servicios":[{"nombre":"Consulta","precio":10}]}`, 201)
}

func TestEdicionConservaDatosYValidaRelacion(t *testing.T) {
	sesion := iniciarSesionPrueba(t, "veterinario")
	crudPrueba(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut && r.Header.Get("X-Agrocampo-Interno") != claveEdicionPrueba {
			t.Error("falta la autenticación interna de la edición")
		}
		switch r.Method + " " + r.URL.Path {
		case "GET /v1/veterinaria/7":
			w.Write([]byte(`{"success":true,"data":{"id":7,"nombre":"Clínica","direccion":"Calle 1","ciudad":"Yopal","telefono":"123","propietario_id":9,"activo":true,"aprobado":true}}`))
		case "PUT /v1/veterinaria/7":
			var cuerpo map[string]interface{}
			json.NewDecoder(r.Body).Decode(&cuerpo)
			if cuerpo["propietario_id"] != float64(9) || cuerpo["ciudad"] != "Yopal" || cuerpo["aprobado"] != true || cuerpo["nombre"] != "Editada" {
				t.Errorf("datos perdidos o campos protegidos modificados: %v", cuerpo)
			}
			w.Write([]byte(`{"success":true}`))
		case "GET /v1/producto_veterinaria/3":
			w.Write([]byte(`{"success":true,"data":{"id":3,"nombre":"Producto","activo":true,"veterinaria":{"id":7},"fecha_creacion":"2026-01-01T00:00:00Z"}}`))
		case "PUT /v1/producto_veterinaria/3":
			var cuerpo map[string]interface{}
			json.NewDecoder(r.Body).Decode(&cuerpo)
			if cuerpo["nombre"] != "Actualizado" || cuerpo["fecha_creacion"] != "2026-01-01T00:00:00Z" || cuerpo["activo"] != true {
				t.Errorf("producto inválido: %v", cuerpo)
			}
			w.Write([]byte(`{"success":true}`))
		case "POST /v1/producto_veterinaria/":
			w.WriteHeader(201)
			w.Write([]byte(`{"success":true,"data":{"id":12,"nombre":"Nuevo","activo":true}}`))
		case "DELETE /v1/producto_veterinaria/3":
			w.Write([]byte(`{"success":true}`))
		default:
			t.Errorf("petición inesperada: %s %s", r.Method, r.URL)
			w.WriteHeader(500)
		}
	})
	pedir(t, "PUT", "/api/v1/veterinarias/7", `{"nombre":"Editada","aprobado":false,"propietario_id":22}`, 200, sesion)
	pedir(t, "POST", "/api/v1/veterinarias/7/productos", `{"nombre":"Nuevo","precio":10}`, 201)
	pedir(t, "PUT", "/api/v1/veterinarias/7/productos/3", `{"nombre":"Actualizado","precio":15}`, 200, sesion)
	pedir(t, "DELETE", "/api/v1/veterinarias/8/productos/3", "", 404)
	pedir(t, "DELETE", "/api/v1/veterinarias/7/productos/3", "", 200)
}

func TestErroresVeterinarias(t *testing.T) {
	sesion := iniciarSesionPrueba(t, "administrador")
	crudPrueba(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/veterinaria/7" {
			w.Write([]byte(`{"success":true,"data":{"id":7}}`))
			return
		}
		// Los controladores generados pueden fallar con HTTP 200.
		w.Write([]byte(`{"success":false,"status":400}`))
	})
	pedir(t, "GET", "/api/v1/veterinarias", "", 502)
	pedir(t, "GET", "/api/v1/veterinarias/0", "", 400)
	pedir(t, "GET", "/api/v1/veterinarias/99", "", 404)
	pedir(t, "POST", "/api/v1/veterinarias", `{"nombre":"Incompleta"}`, 400)
	pedir(t, "PUT", "/api/v1/veterinarias/7", `{"nombre":null}`, 400, sesion)
	pedir(t, "POST", "/api/v1/veterinarias/7/productos", `{"nombre":"Producto","precio":-1}`, 400)
	pedir(t, "POST", "/api/v1/veterinarias/7/productos", `{"nombre":"Producto"}`, 502)
	pedir(t, "DELETE", "/api/v1/veterinarias/7/productos/no-numero", "", 400)
}

func TestCrudNoDisponible(t *testing.T) {
	servidor := httptest.NewServer(http.NotFoundHandler())
	servidor.Close()
	anterior, _ := beego.AppConfig.String("url_veterinaria")
	beego.AppConfig.Set("url_veterinaria", servidor.URL)
	t.Cleanup(func() { beego.AppConfig.Set("url_veterinaria", anterior) })
	pedir(t, "GET", "/api/v1/veterinarias", "", 502)
}
