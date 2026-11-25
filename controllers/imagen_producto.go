package controllers

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	beego "github.com/beego/beego/v2/server/web"
)

type ImagenProductoController struct{ beego.Controller }

const maxImagenProducto = 2 * 1024 * 1024

var nombreImagenProducto = regexp.MustCompile(`^[a-f0-9]{32}\.(png|jpg|gif)$`)

func (c *ImagenProductoController) respuesta(status int, data interface{}, mensaje string) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = map[string]interface{}{"success": status < 400, "data": data, "message": mensaje}
	c.ServeJSON()
}

func (c *ImagenProductoController) Post() {
	c.Ctx.Request.Body = http.MaxBytesReader(c.Ctx.ResponseWriter, c.Ctx.Request.Body, maxImagenProducto+64*1024)
	if err := c.Ctx.Request.ParseMultipartForm(maxImagenProducto + 64*1024); err != nil {
		c.respuesta(400, nil, "Envía una imagen PNG, JPG o GIF de máximo 2 MB")
		return
	}
	if c.Ctx.Request.MultipartForm != nil {
		defer c.Ctx.Request.MultipartForm.RemoveAll()
	}
	archivo, _, err := c.Ctx.Request.FormFile("imagen")
	if err != nil {
		c.respuesta(400, nil, "Selecciona la imagen del producto")
		return
	}
	defer archivo.Close()
	contenido, err := io.ReadAll(io.LimitReader(archivo, maxImagenProducto+1))
	if err != nil || len(contenido) > maxImagenProducto {
		c.respuesta(400, nil, "La imagen supera el máximo de 2 MB")
		return
	}
	config, formato, err := image.DecodeConfig(bytes.NewReader(contenido))
	extensiones := map[string]string{"png": "png", "jpeg": "jpg", "gif": "gif"}
	extension, valido := extensiones[formato]
	if err != nil || !valido || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 40000000 {
		c.respuesta(400, nil, "El archivo debe ser una imagen PNG, JPG o GIF válida")
		return
	}
	if _, _, err := image.Decode(bytes.NewReader(contenido)); err != nil {
		c.respuesta(400, nil, "La imagen está dañada")
		return
	}
	var identificador [16]byte
	if _, err := rand.Read(identificador[:]); err != nil {
		c.respuesta(500, nil, "No fue posible guardar la imagen")
		return
	}
	nombre := hex.EncodeToString(identificador[:]) + "." + extension
	carpeta := filepath.Join("static", "uploads", "productos")
	if err := os.MkdirAll(carpeta, 0755); err != nil {
		c.respuesta(500, nil, "No fue posible guardar la imagen")
		return
	}
	if err := os.WriteFile(filepath.Join(carpeta, nombre), contenido, 0644); err != nil {
		c.respuesta(500, nil, "No fue posible guardar la imagen")
		return
	}
	c.respuesta(201, map[string]string{"url": "/api/v1/imagenes/productos/" + nombre}, "")
}

func (c *ImagenProductoController) PostCuenta() {
	if _, ok := sesionActual(c.Ctx.Request); !ok {
		c.respuesta(401, nil, "Inicia sesión para subir una imagen")
		return
	}
	c.Post()
}

func (c *ImagenProductoController) Get() {
	nombre := c.Ctx.Input.Param(":archivo")
	if !nombreImagenProducto.MatchString(nombre) {
		c.respuesta(404, nil, "Imagen no encontrada")
		return
	}
	c.Ctx.Output.Header("X-Content-Type-Options", "nosniff")
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filepath.Join("static", "uploads", "productos", nombre))
}
