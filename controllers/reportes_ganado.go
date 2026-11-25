package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type cotizacionReporte struct {
	Fecha     string  `json:"fecha"`
	Categoria string  `json:"categoria"`
	Precio    float64 `json:"precio"`
}

type datosReporte struct {
	Cotizaciones []cotizacionReporte `json:"cotizaciones"`
}

type reporteGanado struct {
	ID      int    `json:"id"`
	Usuario int    `json:"solicitado_por"`
	Formato string `json:"formato"`
	Estado  string `json:"estado"`
	Fecha   string `json:"fecha_creacion"`
	Filtros string `json:"filtros"`
	Activo  bool   `json:"activo"`
}

func (r reporteGanado) instantaneaValida() bool {
	var datos datosReporte
	return (r.Formato == "pdf" || r.Formato == "xlsx") && json.Unmarshal([]byte(r.Filtros), &datos) == nil && len(datos.Cotizaciones) > 0
}

func cotizacionesParaReporte(req *http.Request) ([]cotizacionReporte, error) {
	resp, status, err := peticionModulo(req, "crud_agropecuario", "cotizaciones/", http.MethodGet, nil)
	if err != nil || status >= 400 || !resp.Success {
		return nil, errors.New("no se pudieron leer las cotizaciones")
	}
	var filas []cotizacionReporte
	err = json.Unmarshal(resp.Data, &filas)
	return filas, err
}

func (c *ModulosController) DescargarReporte() {
	usuario, ok := sesionActual(c.Ctx.Request)
	if !ok {
		c.responder(401, nil, "Inicia sesión para descargar tus reportes")
		return
	}
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		c.responder(400, nil, "ID inválido")
		return
	}
	resp, status, err := peticionModulo(c.Ctx.Request, "crud_reportes", "reporteactividad/"+strconv.Itoa(id), http.MethodGet, nil)
	if err != nil || status >= 500 {
		c.responder(502, nil, "No fue posible consultar el reporte")
		return
	}
	var reporte reporteGanado
	if !resp.Success || json.Unmarshal(resp.Data, &reporte) != nil || reporte.ID != id || reporte.Usuario != usuario.ID || !reporte.Activo || reporte.Estado != "completado" || !reporte.instantaneaValida() {
		c.responder(404, nil, "Reporte no disponible")
		return
	}
	var datos datosReporte
	if json.Unmarshal([]byte(reporte.Filtros), &datos) != nil {
		c.responder(502, nil, "No fue posible leer el reporte")
		return
	}
	var archivo []byte
	tipo := "application/pdf"
	if reporte.Formato == "pdf" {
		archivo = reportePDF(datos.Cotizaciones)
	} else {
		tipo = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		archivo, err = reporteExcel(datos.Cotizaciones)
		if err != nil {
			c.responder(500, nil, "No fue posible generar el archivo")
			return
		}
	}
	c.Ctx.Output.Header("Content-Type", tipo)
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="ganado-%d.%s"`, id, reporte.Formato))
	c.Ctx.Output.Header("Cache-Control", "private, no-store")
	c.Ctx.ResponseWriter.Write(archivo)
}

func reporteExcel(filas []cotizacionReporte) ([]byte, error) {
	var resultado bytes.Buffer
	libro := zip.NewWriter(&resultado)
	var hoja bytes.Buffer
	hoja.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	celda := func(texto string) {
		hoja.WriteString(`<c t="inlineStr"><is><t>`)
		xml.EscapeText(&hoja, []byte(texto))
		hoja.WriteString(`</t></is></c>`)
	}
	hoja.WriteString("<row>")
	for _, encabezado := range []string{"Fecha", "Categoría", "Precio (COP/kg)"} {
		celda(encabezado)
	}
	hoja.WriteString("</row>")
	for _, fila := range filas {
		hoja.WriteString("<row>")
		celda(fila.Fecha)
		celda(fila.Categoria)
		fmt.Fprintf(&hoja, "<c><v>%.2f</v></c></row>", fila.Precio)
	}
	hoja.WriteString(`</sheetData></worksheet>`)
	partes := map[string]string{
		"[Content_Types].xml":        `<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`,
		"_rels/.rels":                `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`,
		"xl/_rels/workbook.xml.rels": `<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/workbook.xml":            `<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="Cotizaciones" sheetId="1" r:id="rId1"/></sheets></workbook>`,
		"xl/worksheets/sheet1.xml":   hoja.String(),
	}
	for ruta, contenido := range partes {
		parte, err := libro.Create(ruta)
		if err != nil {
			return nil, err
		}
		if _, err = parte.Write([]byte(contenido)); err != nil {
			return nil, err
		}
	}
	if err := libro.Close(); err != nil {
		return nil, err
	}
	return resultado.Bytes(), nil
}

func reportePDF(filas []cotizacionReporte) []byte {
	objetos := []string{"<< /Type /Catalog /Pages 2 0 R >>", "", "<< /Type /Font /Subtype /Type1 /BaseFont /Courier >>"}
	var paginas []string
	for inicio := 0; inicio < len(filas); inicio += 42 {
		var pagina strings.Builder
		pagina.WriteString("BT /F1 11 Tf 15 TL 45 795 Td (AgroCampo - Cotizaciones de ganado) Tj T* T* (Fecha       Categoria                                      COP/kg) Tj T* T*\n")
		for _, fila := range filas[inicio:min(inicio+42, len(filas))] {
			categoria := strings.Map(func(r rune) rune {
				if unicode.Is(unicode.Mn, r) {
					return -1
				}
				if r < 32 || r > 126 {
					return '?'
				}
				return r
			}, norm.NFD.String(fila.Categoria))
			linea := fmt.Sprintf("%-11.10s %-43.43s %12.2f", fila.Fecha, categoria, fila.Precio)
			linea = strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`).Replace(linea)
			fmt.Fprintf(&pagina, "(%s) Tj T*\n", linea)
		}
		pagina.WriteString("ET\n")
		id := len(objetos) + 1
		paginas = append(paginas, fmt.Sprintf("%d 0 R", id))
		objetos = append(objetos, fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 3 0 R >> >> /Contents %d 0 R >>", id+1), fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", pagina.Len(), pagina.String()))
	}
	objetos[1] = fmt.Sprintf("<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(paginas, " "), len(paginas))
	var archivo bytes.Buffer
	archivo.WriteString("%PDF-1.4\n")
	posiciones := make([]int, 0, len(objetos))
	for i, objeto := range objetos {
		posiciones = append(posiciones, archivo.Len())
		fmt.Fprintf(&archivo, "%d 0 obj\n%s\nendobj\n", i+1, objeto)
	}
	inicioXref := archivo.Len()
	fmt.Fprintf(&archivo, "xref\n0 %d\n0000000000 65535 f \n", len(objetos)+1)
	for _, posicion := range posiciones {
		fmt.Fprintf(&archivo, "%010d 00000 n \n", posicion)
	}
	fmt.Fprintf(&archivo, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", len(objetos)+1, inicioXref)
	return archivo.Bytes()
}
