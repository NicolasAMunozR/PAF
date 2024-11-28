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
	historial, err := h.Service.CrearHistorial(codigoPAF, request.Profesor, request.Bloque)
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
	codigoPAF := c.Param("codigo_paf")
	if codigoPAF == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'codigo_paf' es obligatorio"})
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
