// controllers/historial_paf_aceptadas_controller.go
package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

// HistorialPafAceptadasController gestiona las solicitudes relacionadas con los historiales de PAF aceptados
type HistorialPafAceptadasController struct {
	Service *service.HistorialPafAceptadasService
}

// CrearHistorialHandler maneja la creación de un nuevo historial
func (h *HistorialPafAceptadasController) CrearHistorialHandler(c *gin.Context) {
	// Obtener el código PAF desde los parámetros de la URL
	codigoPAFStr := c.Param("codigoPAF")
	comentario := c.Param("comentario")
	codigo_asignatura_pipelsoft := c.Param("cod_asignatura_pipelsoft")
	if codigoPAFStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'codigoPAF' es obligatorio"})
		return
	}

	// Convertir el código PAF a entero
	codigoPAF, err := strconv.Atoi(codigoPAFStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'codigoPAF' debe ser un número entero válido"})
		return
	}

	// Crear un struct auxiliar para parsear el cuerpo de la solicitud
	var request struct {
		Profesor models.ProfesorDB `json:"profesor"`
		Bloque   []string          `json:"bloque"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error al parsear el cuerpo de la solicitud: %v", err)})
		return
	}

	// Llamar al servicio para crear el historial
	historial, err := h.Service.CrearHistorial(codigoPAF, request.Profesor, request.Bloque, codigo_asignatura_pipelsoft, comentario)
	if err != nil {
		log.Printf("Error al crear el historial: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al crear el historial: %v", err)})
		return
	}

	// Responder con el historial creado
	c.JSON(http.StatusCreated, historial)
}

// ObtenerTodosLosHistorialesHandler maneja la solicitud para obtener todos los historiales
func (h *HistorialPafAceptadasController) ObtenerTodosLosHistorialesHandler(c *gin.Context) {
	historiales, err := h.Service.ObtenerTodosLosHistoriales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los historiales"})
		return
	}

	c.JSON(http.StatusOK, historiales)
}

// EliminarHistorialHandler maneja la solicitud para eliminar un historial por CodigoPAF
func (h *HistorialPafAceptadasController) EliminarHistorialHandler(c *gin.Context) {
	// Obtener el CodigoPAF desde los parámetros de la URL
	codigoPAFStr := c.Param("codigo_paf")
	if codigoPAFStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'codigo_paf' es obligatorio"})
		return
	}

	// Convertir el string a int64
	codigoPAF, err := strconv.ParseInt(codigoPAFStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'codigo_paf' debe ser un número válido"})
		return
	}

	// Llamar al servicio para eliminar el historial por CodigoPAF
	if err := h.Service.EliminarHistorial(codigoPAF); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el historial"})
		return
	}

	// Responder con 204 No Content si la eliminación fue exitosa
	c.Status(http.StatusNoContent)
}

// ExportExcel genera y envía un archivo Excel con los datos
// ExportarExcelHandler genera y envía un archivo Excel con los datos
func (h *HistorialPafAceptadasController) ExportarExcelHandler(c *gin.Context) {
	historiales, err := h.Service.ObtenerTodosLosHistoriales()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos"})
		return
	}

	if len(historiales) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No hay historiales para exportar"})
		return
	}

	// Convertir models.HistorialPafAceptadas a service.HistorialPafAceptadas
	var serviceHistoriales []service.HistorialPafAceptadas
	for _, h := range historiales {
		serviceHistoriales = append(serviceHistoriales, service.HistorialPafAceptadas{
			Run:                      h.Run,
			IdPaf:                    h.IdPaf,
			FechaInicioContrato:      h.FechaInicioContrato,
			FechaFinContrato:         h.FechaFinContrato,
			CodigoAsignatura:         h.CodigoAsignatura,
			NombreAsignatura:         h.NombreAsignatura,
			CantidadHoras:            h.CantidadHoras,
			Jerarquia:                h.Jerarquia,
			Calidad:                  h.Calidad,
			SemestrePaf:              h.SemestrePaf,
			DesEstado:                h.DesEstado,
			EstadoProceso:            h.EstadoProceso,
			ProfesorRun:              h.ProfesorRun,
			Semestre:                 h.Semestre,
			Tipo:                     h.Tipo,
			ProfesorCodigoAsignatura: h.ProfesorCodigoAsignatura,
			Seccion:                  h.Seccion,
			Cupo:                     h.Cupo,
			UltimaModificacion:       h.UltimaModificacion,
			Comentario:               h.Comentario,
			Llave:                    h.Llave,
			SemestreInicioPaf:        h.SemestreInicioPaf,
		})
	}

	filePath := "historial.xlsx"
	err = service.GenerateExcel(serviceHistoriales, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar Excel"})
		return
	}

	c.FileAttachment(filePath, "historial.xlsx")
}
