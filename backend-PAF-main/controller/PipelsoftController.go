// controllers/pipelsoft_controller.go
package controller

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

// PipelsoftController gestiona las solicitudes relacionadas con los contratos
type PipelsoftController struct {
	Service *service.PipelsoftService
}

// Constructor del controlador
func NewPipelsoftController(service *service.PipelsoftService) *PipelsoftController {
	return &PipelsoftController{Service: service}
}

// Obtener contratos por código de curso
func (c *PipelsoftController) ObtenerContratosPorCodigoCurso(ctx *gin.Context) {
	codigoCurso := ctx.Param("codigo_curso")

	contratos, err := c.Service.ObtenerContratosPorCodigoCurso(codigoCurso)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}

// Obtener todos los contratos
func (c *PipelsoftController) ObtenerTodosLosContratos(ctx *gin.Context) {
	contratos, err := c.Service.ObtenerTodosLosContratos()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}

// Obtener contratos por RUN
func (c *PipelsoftController) ObtenerContratosPorRUN(ctx *gin.Context) {
	run := ctx.Param("run")

	contratos, err := c.Service.ObtenerContratosPorRUN(run)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}

// Obtener contratos por Código PAF
func (c *PipelsoftController) ObtenerPorCodigoPAF(ctx *gin.Context) {
	codigoPAF := ctx.Param("codigo_paf")

	datos, err := c.Service.ObtenerPorCodigoPAF(codigoPAF)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, datos)
}

// Obtener PAF de los últimos 7 días
func (c *PipelsoftController) ObtenerPAFUltimos7Dias(ctx *gin.Context) {
	datos, err := c.Service.ObtenerPAFUltimos7Dias()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, datos)
}

// Obtener PAF del último mes
func (c *PipelsoftController) ObtenerPAFUltimoMes(ctx *gin.Context) {
	datos, err := c.Service.ObtenerPAFUltimoMes()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, datos)
}

// Obtener contratos por nombre de unidad mayor
func (c *PipelsoftController) ObtenerContratosPorNombreUnidadMayor(ctx *gin.Context) {
	nombreUnidadMayor := ctx.Param("nombreUnidadMayor")

	contratos, err := c.Service.ObtenerContratosPorNombreUnidadMayor(nombreUnidadMayor)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}

// Obtener contratos por nombre de unidad contratante
func (c *PipelsoftController) ObtenerContratosPorNombreUnidadContratante(ctx *gin.Context) {
	nombreUnidadContratante := ctx.Param("nombreUnidadContratante")

	contratos, err := c.Service.ObtenerContratosPorNombreUnidadContratante(nombreUnidadContratante)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}
