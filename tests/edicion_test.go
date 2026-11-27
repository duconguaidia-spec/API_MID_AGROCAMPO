package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	beego "github.com/beego/beego/v2/server/web"
)

const claveEdicionPrueba = "clave-simulada-edicion-de-agrocampo"

func iniciarSesionPrueba(t *testing.T, rol string) *http.Cookie {
	t.Helper()
	usuarios := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/cuenta/login" {
			t.Errorf("petición inesperada al CRUD de usuarios simulado: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, `{"success":true,"data":{"id":7,"rol":%q}}`, rol)
	}))
	t.Cleanup(usuarios.Close)
	for clave, valor := range map[string]string{"crud_usuarios": usuarios.URL, "clave_crud_usuarios": claveEdicionPrueba} {
		anterior, _ := beego.AppConfig.String(clave)
		beego.AppConfig.Set(clave, valor)
		t.Cleanup(func() { beego.AppConfig.Set(clave, anterior) })
	}
	r := httptest.NewRequest(http.MethodPost, "/api/v1/cuenta/login", strings.NewReader(`{"correo":"prueba@example.invalid"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if cookies := w.Result().Cookies(); w.Code == http.StatusOK && len(cookies) == 1 {
		return cookies[0]
	}
	t.Fatalf("no se obtuvo una sesión simulada: %d %s", w.Code, w.Body)
	return nil
}

func TestEdicionRequiereRol(t *testing.T) {
	usuario := iniciarSesionPrueba(t, "usuario")
	for _, ruta := range []string{
		"/api/v1/modulos/temas/7",
		"/api/v1/veterinarias/7",
		"/api/v1/veterinarias/7/productos/3",
	} {
		pedir(t, http.MethodPut, ruta, `{}`, http.StatusUnauthorized)
		pedir(t, http.MethodPut, ruta, `{}`, http.StatusForbidden, usuario)
	}
}

func TestEditarTemaEnviaSoloCamposEditables(t *testing.T) {
	admin := iniciarSesionPrueba(t, "administrador")
	contenido := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/temaforo/12" || r.Header.Get("X-Agrocampo-Interno") != claveEdicionPrueba {
			t.Errorf("petición inesperada al CRUD de contenido simulado: %s %s", r.Method, r.URL.Path)
		}
		var datos map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&datos); err != nil || len(datos) != 3 || datos["titulo"] != "Nuevo" || datos["categoria_id"] != float64(3) || datos["descripcion"] != "Actualizada" {
			t.Errorf("campos inesperados en la edición: %v", datos)
		}
		fmt.Fprint(w, `{"success":true,"data":{"titulo":"Nuevo","descripcion":"Actualizada","categoria_id":3}}`)
	}))
	t.Cleanup(contenido.Close)
	anterior, _ := beego.AppConfig.String("crud_contenido")
	beego.AppConfig.Set("crud_contenido", contenido.URL)
	t.Cleanup(func() { beego.AppConfig.Set("crud_contenido", anterior) })
	pedir(t, http.MethodPut, "/api/v1/modulos/noticias/12", `{}`, http.StatusMethodNotAllowed, admin)
	pedir(t, http.MethodPut, "/api/v1/modulos/temas/12", `{"titulo":" "}`, http.StatusBadRequest, admin)
	respuesta := pedir(t, http.MethodPut, "/api/v1/modulos/temas/12", `{"titulo":" Nuevo ","descripcion":"Actualizada","categoria_id":3,"autor_id":20,"vistas":0}`, http.StatusOK, admin)
	if respuesta["data"].(map[string]interface{})["titulo"] != "Nuevo" {
		t.Fatalf("la MID no devolvió el tema editado: %v", respuesta)
	}
}
