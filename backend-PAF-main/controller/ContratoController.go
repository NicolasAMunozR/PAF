package controller

import (
	"fmt"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

type ContratoController struct {
	Service *service.ContratoService
}

// NewContratoController crea un nuevo controlador de contrato.
func NewContratoController(service *service.ContratoService) *ContratoController {
	return &ContratoController{
		Service: service,
	}
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

	contratos, numElementos, err := c.Service.GetContratosByUnidadMayor(unidad)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear un objeto de respuesta con contratos y el número de elementos
	response := gin.H{
		"contratos":    contratos,
		"numElementos": numElementos,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ContratoController) ProfesoresUnidadMayorPafHandler(ctx *gin.Context) {
	contratoCounts, pipelsoftCounts, err := c.Service.ProfesorUnidadMayorYPaf()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Depuración adicional (si es necesario)
	fmt.Println("Contrato counts:", contratoCounts)
	fmt.Println("Pipelsoft counts:", pipelsoftCounts)

	ctx.JSON(http.StatusOK, gin.H{
		"contratoCounts":  contratoCounts,
		"pipelsoftCounts": pipelsoftCounts,
	})
}

// ProfesorUnidadMayorYPafHandler maneja las solicitudes para contar los contratos y profesores por unidad mayor.
func (c *ContratoController) ProfesorUnidadMayorNOPafHandler(ctx *gin.Context) {
	// Llamar al servicio para obtener los conteos
	contratoCounts, pipelsoftCounts, err := c.Service.ProfesorUnidadMayorNOPaf()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los conteos
	ctx.JSON(http.StatusOK, gin.H{
		"contratoCounts":  contratoCounts,
		"pipelsoftCounts": pipelsoftCounts,
	})
}

// Controlador para manejar la petición de PAF por Unidad Mayor
func (c *ContratoController) GetPafByUnidadMayorHandler(ctx *gin.Context) {
	// Obtener el nombre de la unidad mayor de los parámetros de la URL
	nombreUnidadMayor := ctx.Param("nombreUnidadMayor")

	// Llamar al servicio
	pafs, err := c.Service.GetPafByUnidadMayor(nombreUnidadMayor)
	if err != nil {
		// Responder con un error si el servicio falla
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	ctx.JSON(http.StatusOK, pafs)
}
