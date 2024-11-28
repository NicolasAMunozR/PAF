// controllers/horario_controller.go
package controller

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

// HorarioController gestiona las solicitudes relacionadas con los horarios
type HorarioController struct {
	HorarioService *service.HorarioService
}

// NewHorarioController crea un nuevo controlador para los horarios
func NewHorarioController(horarioService *service.HorarioService) *HorarioController {
	return &HorarioController{HorarioService: horarioService}
}

// ObtenerHorariosPorRun maneja la solicitud para obtener todos los horarios por Run
func (h *HorarioController) ObtenerHorariosPorRun(c *gin.Context) {
	// Obtener el parámetro "run" de la URL
	run := c.Param("run")

	// Llamar al servicio para obtener los horarios
	horarios, err := h.HorarioService.ObtenerHorariosPorRun(run)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error al obtener los horarios: " + err.Error()})
		return
	}

	// Si no se encuentran horarios
	if len(horarios) == 0 {
		c.JSON(404, gin.H{"error": "No se encontraron horarios para el Run especificado"})
		return
	}

	// Retornar los horarios encontrados como JSON
	c.JSON(200, horarios)
}
