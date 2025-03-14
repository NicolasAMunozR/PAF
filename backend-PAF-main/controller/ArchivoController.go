package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

// ArchivoController maneja las solicitudes relacionadas con la generación de PDFs.
type ArchivoController struct {
	Service *service.ArchivoService
}

// NewArchivoController crea un nuevo controlador para archivos.
func NewArchivoController(service *service.ArchivoService) *ArchivoController {
	return &ArchivoController{Service: service}
}

// GenerarPDFHandler maneja la solicitud para generar un PDF basado en un RUN.
func (c *ArchivoController) GenerarPDFHandler(ctx *gin.Context) {
	run := ctx.Param("run")
	if run == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "RUN es requerido"})
		return
	}

	err := service.CrearPDF(c.Service.DB, run)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "PDF generado y guardado correctamente"})
}

// Handler para descargar el archivo PDF
func (ctrl *ArchivoController) DescargarArchivo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pdfData, err := ctrl.Service.GetArchivoPDF(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Archivo no encontrado"})
		return
	}

	// Configurar la respuesta HTTP para descargar el archivo
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=archivo_%d.pdf", id))
	c.Header("Content-Type", "application/pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

// FALTA AÑADIR LOS PROFESORES ACADEMICO
// CreaRContratoHandler maneja la solicitud para obtener los profesores no comunes y académicos sin contrato, filtrando por semestre.
func (c *ArchivoController) CreaRContratoHandler(ctx *gin.Context) {
	// Obtener el unidadMenor desde la URL
	unidadMenor := ctx.Param("unidadMenor")

	// Validar que el unidadMenor no esté vacío
	if unidadMenor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMenor' es requerido"})
		return
	}

	// Llamar al servicio con el unidadMenorproporcionado
	rutsNoComunes, err := service.CreaRContratoAutomaticamentePorUnidad(c.Service.DB, unidadMenor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con ambas listas en la respuesta JSON
	ctx.JSON(http.StatusOK, gin.H{
		"profesores que no se genero paf": rutsNoComunes,
	})
}

// AgregarComentarioHandler maneja la solicitud para agregar un comentario a un archivo por ID
func (c *ArchivoController) AgregarComentarioHandler(ctx *gin.Context) {
	// Obtener el ID desde la URL
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener el comentario desde el cuerpo de la solicitud
	var requestBody struct {
		Comentario string `json:"comentario"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de JSON inválido"})
		return
	}

	// Llamar al servicio para agregar el comentario
	if err := c.Service.AgregarComentario(uint(id), requestBody.Comentario); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"mensaje": "Comentario agregado exitosamente"})
}

// SubirArchivoHandler maneja la subida de archivos
func (c *ArchivoController) SubirArchivoHandler(ctx *gin.Context) {
	// Obtener el ID desde la URL
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener el archivo desde el formulario
	file, err := ctx.FormFile("archivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el archivo"})
		return
	}

	// Abrir el archivo
	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo"})
		return
	}
	defer openedFile.Close()

	// Leer los datos del archivo
	fileData, err := ioutil.ReadAll(openedFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el archivo"})
		return
	}

	// Llamar al servicio para guardar el archivo
	if err := c.Service.SubirArchivo(uint(id), file.Filename, fileData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"mensaje": "Archivo subido exitosamente"})
}

// ObtenerListaArchivosHandler devuelve la lista de archivos adjuntos de un Archivo específico
func (c *ArchivoController) ObtenerListaArchivosHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	archivos, err := c.Service.ObtenerListaArchivos(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, archivos)
}

// ModificarArchivoHandler permite actualizar un archivo adjunto
func (c *ArchivoController) ModificarArchivoHandler(ctx *gin.Context) {
	// Obtener el ID del archivo adjunto
	idParam := ctx.Param("archivoAdjuntoID")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener el archivo desde el formulario
	file, err := ctx.FormFile("archivo")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el archivo"})
		return
	}

	// Abrir el archivo
	openedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al abrir el archivo"})
		return
	}
	defer openedFile.Close()

	// Leer los datos del archivo
	fileData, err := ioutil.ReadAll(openedFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer el archivo"})
		return
	}

	// Llamar al servicio para actualizar el archivo
	if err := c.Service.ModificarArchivo(uint(id), file.Filename, fileData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"mensaje": "Archivo modificado exitosamente"})
}

// GetArchivosAdjuntosByRutHandler obtiene archivos adjuntos por RUT
func (c *ArchivoController) GetArchivosAdjuntosByRutHandler(ctx *gin.Context) {
	rut := ctx.Param("rut")
	if rut == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El RUT es requerido"})
		return
	}

	archivosAdjuntos, err := c.Service.GetArchivosAdjuntosByRut(rut)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"archivos_adjuntos": archivosAdjuntos})
}

// DownloadArchivoHandler permite descargar un archivo adjunto por ID
func (c *ArchivoController) DownloadArchivoHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El ID del archivo es requerido"})
		return
	}

	// Convertir el ID a uint
	var archivoID uint
	if _, err := fmt.Sscanf(idParam, "%d", &archivoID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener el archivo desde el servicio
	archivo, err := c.Service.GetArchivoAdjuntoByID(archivoID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el archivo tiene datos
	if len(archivo.Datos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "El archivo está vacío"})
		return
	}

	// Asegurar que el nombre del archivo tiene la extensión .pdf
	nombreArchivo := archivo.Nombre
	if !strings.HasSuffix(strings.ToLower(nombreArchivo), ".pdf") {
		nombreArchivo += ".pdf"
	}

	// Configurar los encabezados para la descarga
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", nombreArchivo))
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(archivo.Datos)))

	// Enviar los datos del archivo
	ctx.Data(http.StatusOK, "application/pdf", archivo.Datos)
}

//ObtenerProfesoresQueNoSePuedeGenerarContrato

// ObtenerProfesoresHandler maneja la solicitud para obtener los profesores que no pueden generar contrato.
func (c *ArchivoController) ObtenerProfesoresQueNoSePuedeGenerarContrato(ctx *gin.Context) {
	// Llamar al servicio para obtener los profesores sin filtro por unidad mayor
	rutsNoComunes, rutsContratables, rutsConContrato, err := service.ObtenerProfesoresQueNoSePuedeGenerarContrato(c.Service.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con las listas de profesores
	ctx.JSON(http.StatusOK, gin.H{
		"profesores_no_comunes":   rutsNoComunes,
		"profesores_contratables": rutsContratables,
		"profesores_con_contrato": rutsConContrato,
	})
}

func (c *ArchivoController) ObtenerProfesoresQueNoSePuedeGenerarContratoUnidadMayor(ctx *gin.Context) {
	// Obtener el parámetro 'unidad_mayor' desde la URL
	unidadMayor := ctx.Param("unidad_mayor")

	// Validar que el parámetro 'unidad_mayor' no esté vacío
	if unidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidad_mayor' es requerido"})
		return
	}

	// Llamar al servicio para obtener los profesores filtrados por la unidad mayor
	rutsNoComunes, rutsAcademicos, rutsContratables, rutsConContrato, err := service.ObtenerProfesoresQueNoSePuedeGenerarContratoUnidadMayor(c.Service.DB, unidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con las listas de profesores
	ctx.JSON(http.StatusOK, gin.H{
		"profesores_no_comunes":   rutsNoComunes,
		"profesores_academicos":   rutsAcademicos,
		"profesores_contratables": rutsContratables,
		"profesores_con_contrato": rutsConContrato,
	})
}

func (c *ArchivoController) GenerarPDFSinDatos(ctx *gin.Context) {
	var request models.Archivo1

	// Decodificar JSON en la estructura
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos de entrada inválidos"})
		return
	}

	// Llamar al servicio para crear el PDF
	id, err := service.CrearPDFSinData(
		c.Service.DB,
		request.UnidadMayor,
		request.UnidadMenor,
		request.NumeroCentroDeCostos,
		request.CelulaIdentidad,
		request.LugarNacimiento,
		request.FechaYHoraNacimiento,
		request.ApellidoP,
		request.ApellidoM,
		request.Nombres,
		request.Nacionalidad,
		request.Domicio,
		request.Correo,
		request.Titulo,
		request.Institucion,
		request.FechaObtencion,
		request.NumeroSemestre,
		request.GradoAcademico,
		request.InstitucionGradoAcademico,
		request.FechaObtencionGradoAcademico,
		request.TipoIngreso,
		request.Cargo,
		request.Nivel,
		request.Grado,
		request.Rango,
		request.Funcion,
		request.Jerarquia,
		request.Asignatura,
		request.NumeroHoras,
		request.Categoria,
		request.Calidad,
		request.LugarDesempeño,
		request.CargoOtroPublico,
		request.GradoOtroPublico,
		request.NivelOtroPublico,
		request.RangoOtroPublico,
		request.NumeroHorasOtroPublico,
		request.CalidadOtroPublico,
		request.FechaInicioContrato,
		request.FechaFinContrato,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "PDF generado correctamente",
		"id": id,
	})
}

// Endpoint para buscar archivos por célula de identidad (RUT del docente)
func (c *ArchivoController) GetArchivoPorCelula(ctx *gin.Context) {
	celulaIdentidad := ctx.Param("celulaIdentidad") // Capturar parámetro de la URL

	archivos, err := c.Service.BuscarArchivoPorCelula(celulaIdentidad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar los archivos"})
		return
	}
	if len(archivos) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron archivos para este RUT"})
		return
	}

	ctx.JSON(http.StatusOK, archivos)
}
