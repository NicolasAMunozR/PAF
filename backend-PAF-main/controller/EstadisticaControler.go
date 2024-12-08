package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service" // Cambiar según la ruta correcta del paquete
	"github.com/gin-gonic/gin"
)

// EstadisticasController gestiona las solicitudes relacionadas con las estadísticas
type EstadisticasController struct {
	Service *service.EstadisticasService
}

// NewEstadisticasController crea una nueva instancia de EstadisticasController
func NewEstadisticasController(service *service.EstadisticasService) *EstadisticasController {
	return &EstadisticasController{Service: service}
}

// ObtenerEstadisticas maneja la solicitud para obtener las estadísticas generales
func (c *EstadisticasController) ObtenerEstadisticas(ctx *gin.Context) {
	// Llamar al servicio para obtener las estadísticas
	estadisticas, err := c.Service.ObtenerEstadisticas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al obtener estadísticas: %v", err)})
		return
	}

	// Responder con las estadísticas en formato JSON
	ctx.JSON(http.StatusOK, estadisticas)
}

// ContarRegistrosPorUnidadMayor maneja la solicitud para contar registros por nombre de unidad Mayor
func (c *EstadisticasController) ContarRegistrosPorUnidadMayor(ctx *gin.Context) {
	// Obtener el nombre de la unidad Mayor desde los parámetros de la URL
	nombreUnidadMayor := ctx.Param("nombreUnidadMayor")

	// Llamar al servicio para contar los registros
	count, err := c.Service.ContarRegistrosPorNombreUnidadMayor(nombreUnidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al contar registros: %v", err)})
		return
	}

	// Responder con el conteo en formato JSON
	ctx.JSON(http.StatusOK, gin.H{
		"nombreUnidadMayor": nombreUnidadMayor,
		"conteo":            count,
	})
}

// ContarRegistrosPorCodEstado maneja la solicitud para contar registros de Pipelsoft donde el `cod_estado` no sea "F1", "F9" ni "A9"
func (c *EstadisticasController) ContarRegistrosPorCodEstado(ctx *gin.Context) {
	// Llamar al servicio para obtener el conteo y el porcentaje
	count, porcentaje, err := c.Service.ContarRegistrosExcluyendoEstados()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al contar registros excluyendo estados: %v", err)})
		return
	}

	// Responder con el conteo y porcentaje en formato JSON
	ctx.JSON(http.StatusOK, gin.H{
		"conteo":     count,
		"porcentaje": porcentaje,
	})
}

// ActualizarBanderaAceptacion maneja la solicitud para actualizar la bandera de aceptación
func (ctrl *HistorialPafAceptadasController) ActualizarBanderaAceptacion(ctx *gin.Context) {
	// Obtener los parámetros del request
	codigoPAF := ctx.Param("codigoPAF") // Se espera el código PAF en la URL

	var requestBody struct {
		NuevaBanderaAceptacion int `json:"nuevaBanderaAceptacion"` // Nueva bandera de aceptación en el cuerpo de la petición
	}

	// Parsear el cuerpo de la solicitud
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Llamar al servicio para actualizar la BanderaAceptacion
	err := ctrl.Service.ActualizarBanderaAceptacion(codigoPAF, requestBody.NuevaBanderaAceptacion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al actualizar BanderaAceptacion: %v", err)})
		return
	}

	// Responder con éxito
	ctx.JSON(http.StatusOK, gin.H{"message": "BanderaAceptacion actualizada correctamente"})
}

// ObtenerFrecuenciaNombreUnidadMayor maneja la solicitud para obtener la frecuencia de NombreUnidadMayor
func (c *EstadisticasController) ObtenerFrecuenciaNombreUnidadMayor(ctx *gin.Context) {
	// Llamar al servicio para obtener los datos
	frecuencia, err := c.Service.ObtenerFrecuenciaNombreUnidadMayor()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error al obtener la frecuencia de NombreUnidadMayor: %v", err)})
		return
	}

	// Responder con los datos en formato JSON
	ctx.JSON(http.StatusOK, frecuencia)
}

func (c *EstadisticasController) ObtenerEstadisticasPorUnidadMayorHandler(ctx *gin.Context) {
	// Obtener el valor de 'unidad-mayor' desde los parámetros de la URL
	unidadMayor := ctx.Param("unidad-mayor")

	// Validar que el parámetro no esté vacío
	if unidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidad-mayor' es obligatorio"})
		return
	}

	// Llamar al servicio
	resp, err := c.Service.ObtenerEstadisticasPorUnidadMayor(unidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	ctx.JSON(http.StatusOK, resp)
}

func (c *EstadisticasController) ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayorHandler(ctx *gin.Context) {
	// Obtener el parámetro de la URL
	nombreUnidadMayor := ctx.Param("unidad-mayor")

	// Llamar al servicio
	frecuencia, err := c.Service.ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayor(nombreUnidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos
	ctx.JSON(http.StatusOK, frecuencia)
}

// ObtenerPafActivasPorUnidadHandler maneja el conteo de registros y RUNs únicos
func (c *EstadisticasController) ObtenerPafActivasPorUnidadHandler(ctx *gin.Context) {
	// Obtener el parámetro unidadMayor desde la URL
	unidadMayor := ctx.Param("unidadMayor")
	if unidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro unidadMayor es requerido",
		})
		return
	}

	// Llamar al servicio para obtener los datos
	count, totalRUNs, err := c.Service.ContarRegistrosPorUnidadMayorConRuns(unidadMayor)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Responder con los resultados
	ctx.JSON(http.StatusOK, gin.H{
		"totalRegistros": count,
		"totalRUNs":      totalRUNs,
	})
}

// ObtenerRUNUnicosExcluidosHandler obtiene los RUN únicos de ProfesorDB que no están en Pipelsoft
func (ctrl *EstadisticasController) ObtenerRUNUnicosExcluidosHandler(c *gin.Context) {
	// Llamar al servicio para obtener los RUN excluidos
	profesorRuns, excluidos, err := ctrl.Service.ObtenerRUNUnicosExcluidos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los RUNs excluidos"})
		return
	}

	// Retornar los resultados
	c.JSON(http.StatusOK, gin.H{
		"profesorRuns": profesorRuns,
		"excluidos":    excluidos,
	})
}

// CompararRunsHandler compara los RUN excluidos con los RUN en Contrato
func (ctrl *EstadisticasController) CompararRunsHandler(c *gin.Context) {
	// Obtener los RUNs excluidos desde el parámetro de la solicitud
	var runsExcluidos []string
	if err := c.ShouldBindJSON(&runsExcluidos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros inválidos"})
		return
	}

	// Llamar al servicio para comparar los RUNs excluidos
	noEncontrados, cantidad, err := ctrl.Service.CompararRuns(runsExcluidos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al comparar los RUNs"})
		return
	}

	// Retornar los resultados
	c.JSON(http.StatusOK, gin.H{
		"noEncontrados": noEncontrados,
		"cantidad":      cantidad,
	})
}

// ObtenerYCompararRunsHandler maneja la ruta para obtener los RUNs únicos de ProfesorDB que no están en Pipelsoft
// y compara esos RUNs con los RUNs en los contratos.
func (ctrl *EstadisticasController) ObtenerYCompararRunsHandler(c *gin.Context) {
	// Obtener los RUNs excluidos de ProfesorDB que no están en Pipelsoft
	_, excluidos, err := ctrl.Service.ObtenerRUNUnicosExcluidos()
	if err != nil {
		log.Println("Error al obtener los RUNs excluidos:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los RUNs excluidos"})
		return
	}

	// Comparar los RUNs excluidos con los RUNs en los contratos
	noEncontrados, cantidad, err := ctrl.Service.CompararRuns(excluidos)
	if err != nil {
		log.Println("Error al comparar los RUNs:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al comparar los RUNs"})
		return
	}

	// Retornar los resultados en el formato solicitado
	c.JSON(http.StatusOK, gin.H{
		"noEncontrados": noEncontrados,
		"cantidad":      cantidad,
	})
}

func (c *EstadisticasController) ObtenerUnidadesMayoresHandler(ctx *gin.Context) {
	// Llamar al servicio para obtener los datos
	unidades, err := c.Service.ObtenerUnidadesMayoresConProfesores()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Responder con los resultados
	ctx.JSON(http.StatusOK, gin.H{
		"unidadesMayores": unidades,
	})
}

func (c *EstadisticasController) ObtenerUnidadesMayoresSinProfesoresEnPipelsoftHandler(ctx *gin.Context) {
	resultado, err := c.Service.ObtenerUnidadesMayoresSinProfesoresEnPipelsoft()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

func (c *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltradosHandler(ctx *gin.Context) {
	resultado, err := c.Service.ObtenerUnidadesMayoresConProfesoresFiltrados()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// ObtenerUnidadesMayoresConProfesoresFiltrados maneja la solicitud HTTP para obtener unidades mayores con los profesores filtrados
func (ctrl *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltrados(c *gin.Context) {
	// Llamar al servicio para obtener las unidades mayores con los profesores filtrados
	resultado, err := ctrl.Service.ObtenerUnidadesMayoresConProfesoresFiltrados()
	if err != nil {
		log.Println("Error al obtener unidades mayores con profesores filtrados:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// ObtenerUnidadesMayoresSinProfesoresEnPipelsoft maneja la solicitud HTTP para obtener unidades mayores sin profesores en Pipelsoft
func (ctrl *EstadisticasController) ObtenerUnidadesMayoresSinProfesoresEnPipelsoft(c *gin.Context) {
	// Llamar al servicio para obtener las unidades mayores sin profesores en Pipelsoft
	resultado, err := ctrl.Service.ObtenerUnidadesMayoresSinProfesoresEnPipelsoft()
	if err != nil {
		log.Println("Error al obtener unidades mayores sin profesores en Pipelsoft:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// ObtenerEstadisticasPorUnidad maneja la solicitud HTTP para obtener estadísticas por unidad
func (ctrl *EstadisticasController) ObtenerEstadisticasPorUnidad(c *gin.Context) {
	// Obtener parámetros de la solicitud
	unidadMayor := c.DefaultQuery("unidadMayor", "") // "unidadMayor" es requerido
	unidadMenor := c.DefaultQuery("unidadMenor", "")

	// Validar los parámetros
	if unidadMayor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMayor' es obligatorio"})
		return
	}

	// Llamar al servicio para obtener las estadísticas
	resp, err := ctrl.Service.ObtenerEstadisticasPorUnidad(unidadMayor, unidadMenor)
	if err != nil {
		log.Println("Error al obtener estadísticas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, resp)
}
