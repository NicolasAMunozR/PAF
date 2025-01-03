// controllers/pipelsoft_controller.go
package controller

import (
	"net/http"

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

// Obtener contratos por nombre de unidad menor
func (c *PipelsoftController) ObtenerContratosPorNombreUnidadMenor(ctx *gin.Context) {
	nombreUnidadMenor := ctx.Param("nombreUnidadMenor")

	contratos, err := c.Service.ObtenerContratosPorNombreUnidadMenor(nombreUnidadMenor)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, contratos)
}

// ObtenerUnidadesMenores obtiene todas las unidades menores asociadas a una unidad mayor
func (c *PipelsoftController) ObtenerUnidadesMenores(ctx *gin.Context) {
	// Leer el parámetro NombreUnidadMayor desde la URL o consulta
	nombreUnidadMayor := ctx.Query("nombreUnidadMayor")
	if nombreUnidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro nombreUnidadMayor es requerido"})
		return
	}

	// Llamar al servicio para obtener las unidades menores
	unidadesMenores, err := c.Service.ObtenerNombreUnidadMenorPorMayor(nombreUnidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar la base de datos"})
		return
	}

	// Responder con la lista de unidades menores
	ctx.JSON(http.StatusOK, gin.H{"unidadesMenores": unidadesMenores})
}

// Handler para buscar por RunEmpleado y retornar múltiples registros
func (controller *PipelsoftController) ObtenerContratosPorRUNMostrarTodo(c *gin.Context) {
	runEmpleado := c.Param("rut")

	// Llamar al servicio
	records, err := controller.Service.ObtenerContratosPorRUNMostrarTodo(runEmpleado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar los registros"})
		return
	}

	// Si no se encontraron registros, devolver un mensaje adecuado
	if len(records) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron registros"})
		return
	}

	// Devolver todos los registros encontrados
	c.JSON(http.StatusOK, records)
}

// Handler para buscar por RunEmpleado y retornar múltiples registros
func (controller *PipelsoftController) ObtenerContratosPorIdPafMostrarTodo(c *gin.Context) {
	runEmpleado := c.Param("rut")

	// Llamar al servicio
	records, err := controller.Service.ObtenerContratosPorIdPafMostrarTodo(runEmpleado)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al buscar los registros"})
		return
	}

	// Si no se encontraron registros, devolver un mensaje adecuado
	if len(records) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron registros"})
		return
	}

	// Devolver todos los registros encontrados
	c.JSON(http.StatusOK, records)
}

func (c *PipelsoftController) GetUniqueUnits(ctx *gin.Context) {
	mayores, menores, err := c.Service.GetUniqueUnits()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"NombreUnidadMayor": mayores,
		"NombreUnidadMenor": menores,
	})
}

func (ctrl *PipelsoftController) GetUnitsByMayor(c *gin.Context) {
	nombreUnidadMayor := c.Param("unidadMayor")
	if nombreUnidadMayor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMayor' es requerido"})
		return
	}

	menores, err := ctrl.Service.GetUnitsByMayor(nombreUnidadMayor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener las unidades menores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unidadMenor": menores})
}

func (ctrl *PipelsoftController) GetBySemester(c *gin.Context) {
	semestre := c.Param("semestre")
	if semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'semestre' es requerido"})
		return
	}

	pipelsofts, err := ctrl.Service.GetBySemester(semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los elementos filtrados por semestre"})
		return
	}

	c.JSON(http.StatusOK, pipelsofts)
}
