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

// 1
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

// 2
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

// 3
func (c *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltradosHandler(ctx *gin.Context) {
	resultado, err := c.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivas()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// 4
// ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivos obtiene las unidades mayores con profesores activos (PAF).
func (c *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivos(ctx *gin.Context) {
	// Llamamos al servicio para obtener los datos
	resultado, err := c.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivos()
	if err != nil {
		// Si ocurre un error, devolvemos una respuesta 500 con el mensaje de error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Si la operación fue exitosa, devolvemos los resultados en formato JSON
	ctx.JSON(http.StatusOK, resultado)
}

// 5
func (h *EstadisticasController) ObtenerUnidadesMayoresPorCodEstadoPAF(c *gin.Context) {
	// Obtener el parámetro codEstadoPAF desde la URL
	codEstadoPAF := c.Param("codEstadoPAF")

	// Llamar al servicio
	resultado, err := h.Service.ObtenerUnidadesMayoresPorCodEstadoPAF(codEstadoPAF)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// 6
// ObtenerEstadisticasPorUnidad maneja la solicitud HTTP para obtener estadísticas por unidad
func (ctrl *EstadisticasController) ObtenerEstadisticasPorUnidad(c *gin.Context) {
	// Obtener parámetros de la ruta
	unidadMayor := c.Param("unidadMayor") // Parámetro de la ruta
	unidadMenor := c.Param("unidadMenor") // Parámetro de la ruta

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

// 7
// ObtenerEstadisticasPorUnidadTOTO maneja la solicitud HTTP para obtener estadísticas por unidad
func (ctrl *EstadisticasController) ObtenerEstadisticasPorUnidadTOTO(c *gin.Context) {
	// Obtener parámetros de la ruta
	unidadMayor := c.Param("unidadMayor") // Parámetro de la ruta
	unidadMenor := c.Param("unidadMenor") // Parámetro de la ruta (opcional)

	// Validar que al menos 'unidadMayor' esté presente
	if unidadMayor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMayor' debe ser proporcionado"})
		return
	}

	// Llamar al servicio para obtener las estadísticas
	resp, err := ctrl.Service.ObtenerEstadisticasPorUnidadTOTO(unidadMayor, unidadMenor)
	if err != nil {
		log.Println("Error al obtener estadísticas:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, resp)
}

// 8.1
// ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos obtiene las unidades menores con profesores activos (PAF).
func (h *EstadisticasController) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(c *gin.Context) {
	// Obtener el parámetro unidadMayor desde la URL
	unidadMayor := c.Param("unidadMayor")

	// Llamar al servicio
	resultado, err := h.Service.ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(unidadMayor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// 8.3
// ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor obtiene unidades mayores filtradas por 'unidadMayor' y profesores activos
func (ctrl *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor(c *gin.Context) {
	// Obtener el parámetro 'unidadMayor' desde la URL o como parámetro de consulta
	unidadMayor := c.DefaultQuery("unidadMayor", "")

	// Validar que se haya recibido el parámetro 'unidadMayor'
	if unidadMayor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMayor' es obligatorio"})
		return
	}

	// Llamar al servicio para obtener las unidades mayores con profesores activos filtrados por 'unidadMayor'
	unidadesInactivas, err := ctrl.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor(unidadMayor)
	if err != nil {
		// Si hubo un error, devolver una respuesta con el mensaje de error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver la respuesta con el mapa de unidades mayores y el total de profesores
	c.JSON(http.StatusOK, gin.H{"unidades": unidadesInactivas})
}

// 8.3
// ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3 maneja la solicitud para obtener unidades menores sin profesores en Pipelsoft (versión 8.3).
func (c *EstadisticasController) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(ctx *gin.Context) {
	// Recuperamos el parámetro 'unidadMayor' desde la ruta
	unidadMayor := ctx.Param("unidadMayor")

	// Llamamos al servicio para obtener los datos, pasando el parámetro 'unidadMayor' que se ha obtenido de la ruta
	resultado, err := c.Service.ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(unidadMayor)
	if err != nil {
		// Si ocurre un error, devolvemos una respuesta 500 con el mensaje de error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Si la operación fue exitosa, devolvemos los resultados en formato JSON
	ctx.JSON(http.StatusOK, resultado)
}

// 8.4
// ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos maneja la solicitud para obtener unidades menores con profesores filtrados por PAF activos.
func (c *EstadisticasController) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos(ctx *gin.Context) {
	// Recuperamos el parámetro 'unidadMayor' desde la ruta
	unidadMayor := ctx.Param("unidadMayor")

	// Llamamos al servicio para obtener los datos, pasando el parámetro 'unidadMayor' que se ha obtenido de la ruta
	resultado, err := c.Service.ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos(unidadMayor)
	if err != nil {
		// Si ocurre un error, devolvemos una respuesta 500 con el mensaje de error
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Si la operación fue exitosa, devolvemos los resultados en formato JSON
	ctx.JSON(http.StatusOK, resultado)
}

// 8.5
func (h *EstadisticasController) ObtenerUnidadesMenoresPorCodEstadoPAF(c *gin.Context) {
	// Obtener el parámetro codEstadoPAF desde la URL
	codEstadoPAF := c.Param("codEstadoPAF")

	// Llamar al servicio
	resultado, err := h.Service.ObtenerUnidadesMenoresPorCodEstadoPAF(codEstadoPAF)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}
