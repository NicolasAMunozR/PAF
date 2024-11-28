package controller

import (
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

type ContratoController struct {
	Service *service.ContratoService
}

// GetAllContratosHandler maneja las solicitudes para obtener todos los contratos.
func (c *ContratoController) GetAllContratosHandler(ctx *gin.Context) {
	contratos, err := c.Service.GetAllContratos()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contratos)
}

// GetContratoByRunHandler maneja las solicitudes para obtener un contrato por el RUN del docente.
func (c *ContratoController) GetContratoByRunHandler(ctx *gin.Context) {
	run := ctx.Param("run") // Extrae el parámetro 'run' de la URL.
	if run == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'run' es obligatorio"})
		return
	}

	contrato, err := c.Service.GetContratoByRun(run)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if contrato == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Contrato no encontrado"})
		return
	}
	ctx.JSON(http.StatusOK, contrato)
}

// GetContratosByUnidadMayorHandler maneja las solicitudes para obtener contratos por unidad mayor.
func (c *ContratoController) GetContratosByUnidadMayorHandler(ctx *gin.Context) {
	unidad := ctx.DefaultQuery("unidad", "") // Obtiene el parámetro 'unidad' de la query string.
	if unidad == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidad' es obligatorio"})
		return
	}

	contratos, err := c.Service.GetContratosByUnidadMayor(unidad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contratos)
}
