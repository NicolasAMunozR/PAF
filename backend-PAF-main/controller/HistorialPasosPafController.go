package controller

import (
	"fmt"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"

	"github.com/gin-gonic/gin"
)

type HistorialPasosPafController struct {
	Service *service.HistorialPasosPafService
}

func NewHistorialPasosPafController(service *service.HistorialPasosPafService) *HistorialPasosPafController {
	return &HistorialPasosPafController{Service: service}
}

func (h *HistorialPasosPafController) ObtenerHistorialYDuracionesPorIdYRun(c *gin.Context) {
	idPaf := c.Param("id_paf")
	runDocente := c.Param("run_docente")

	historiales, duraciones, err := h.Service.ObtenerYCalcularPorIdYRun(idPaf, runDocente)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error al obtener los datos: %v", err)})
		return
	}

	// Crear una respuesta que combine los datos y duraciones
	var respuesta []gin.H
	for i, historial := range historiales {
		respuesta = append(respuesta, gin.H{
			"id_paf":                 historial.IdPaf,
			"run_docente":            historial.RunDocente,
			"estado_nuevo_paf":       historial.EstadoNuevoPaf,
			"codigo_estado_paf":      historial.CodigoEstadoPaf,
			"fecha_llegada_paf":      historial.FechaLlegadaPaf,
			"fecha_modificacion_paf": historial.FechaModificacionPaf,
			"duracion":               duraciones[i].String(), // Duración en formato legible
		})
	}

	c.JSON(http.StatusOK, respuesta)
}
