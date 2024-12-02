package controller

import (
	"fmt"
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
