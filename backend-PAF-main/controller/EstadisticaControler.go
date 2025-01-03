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

	semestre := ctx.Param("semestreId")
	// Llamar al servicio para obtener las estadísticas
	estadisticas, err := c.Service.ObtenerEstadisticas(semestre)
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
	semestre := ctx.Param("semestre")

	// Llamar al servicio para contar los registros
	count, err := c.Service.ContarRegistrosPorNombreUnidadMayor(nombreUnidadMayor, semestre)
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
	semestre := ctx.Param("semestre")
	count, porcentaje, err := c.Service.ContarRegistrosExcluyendoEstados(semestre)
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
	semestre := ctx.Param("semestre")
	frecuencia, err := c.Service.ObtenerFrecuenciaNombreUnidadMayor(semestre)
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
	semestre := ctx.Param("semestre")

	// Validar que el parámetro no esté vacío
	if unidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidad-mayor' es obligatorio"})
		return
	}

	// Llamar al servicio
	resp, err := c.Service.ObtenerEstadisticasPorUnidadMayor(unidadMayor, semestre)
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
	semestre := ctx.Param("semestre")

	// Llamar al servicio
	frecuencia, err := c.Service.ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayor(nombreUnidadMayor, semestre)
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
	semestre := ctx.Param("semestre")
	if unidadMayor == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro unidadMayor es requerido",
		})
		return
	}

	// Llamar al servicio para obtener los datos
	count, totalRUNs, err := c.Service.ContarRegistrosPorUnidadMayorConRuns(unidadMayor, semestre)
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
	// Obtener el parámetro codEstadoPAF desde la URL
	semestre := ctx.Param("semestre")
	if semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'codEstadoPAF' es obligatorio",
		})
		return
	}

	// Llamar al servicio para obtener los datos
	unidades, err := c.Service.ObtenerUnidadesMayoresConProfesores(semestre)
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
	// Obtener el parámetro codEstadoPAF desde la URL
	semestre := ctx.Param("semestre")
	if semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'semestre' es obligatorio",
		})
		return
	}

	// Llamar al servicio con el filtro de semestre
	resultado, err := c.Service.ObtenerUnidadesMayoresSinProfesoresEnPipelsoftPorSemestre(semestre)
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
	// Obtener el parámetro semestre desde la URL
	semestre := ctx.Param("semestre")
	if semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'semestre' es obligatorio",
		})
		return
	}

	// Llamar al servicio con el filtro de semestre
	resultado, err := c.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorSemestre(semestre)
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
	// Obtener el parámetro semestre desde la URL
	semestre := ctx.Param("semestre")
	if semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "El parámetro 'semestre' es obligatorio",
		})
		return
	}

	// Llamar al servicio con el filtro de semestre
	resultado, err := c.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivosPorSemestre(semestre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// 5
func (c *EstadisticasController) ObtenerUnidadesMayoresPorCodEstadoPAF(ctx *gin.Context) {
	// Obtener los parámetros codEstadoPAF y semestre desde la URL
	codEstadoPAF := ctx.Param("codEstadoPAF")
	semestre := ctx.Param("semestre")

	// Validar parámetros
	if codEstadoPAF == "" || semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Los parámetros 'codEstadoPAF' y 'semestre' son obligatorios",
		})
		return
	}

	// Llamar al servicio
	resultado, err := c.Service.ObtenerUnidadesMayoresPorCodEstadoPAF(codEstadoPAF, semestre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// 6
// ObtenerEstadisticasPorUnidad maneja la solicitud HTTP para obtener estadísticas por unidad
func (ctrl *EstadisticasController) ObtenerEstadisticasPorUnidad(c *gin.Context) {
	// Obtener parámetros de la ruta
	unidadMayor := c.Param("unidadMayor")
	unidadMenor := c.Param("unidadMenor")

	// Obtener el parámetro 'semestre' como query parameter
	semestre := c.Param("semestre")

	// Validar los parámetros obligatorios
	if unidadMayor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMayor' es obligatorio"})
		return
	}
	if unidadMenor == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'unidadMenor' es obligatorio"})
		return
	}

	// Llamar al servicio para obtener las estadísticas
	resp, err := ctrl.Service.ObtenerEstadisticasPorUnidad(unidadMayor, unidadMenor, semestre)
	if err != nil {
		log.Printf("Error al obtener estadísticas: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta en formato JSON
	c.JSON(http.StatusOK, resp)
}

// 7
// ObtenerEstadisticasPorUnidadTOTO maneja la solicitud HTTP para obtener estadísticas por unidad
func (ctrl *EstadisticasController) ContarRegistrosPorUnidadMayorYUnidadMenor(c *gin.Context) {
	// Obtener los parámetros de la URL
	unidadMayor := c.Param("unidadMayor")
	unidadMenor := c.Param("unidadMenor")
	semestre := c.Param("semestre") // Capturar el parámetro semestre

	// Validar que todos los parámetros no estén vacíos
	if unidadMayor == "" || unidadMenor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio para obtener los conteos
	count, totalRUNs, err := ctrl.Service.ContarRegistrosPorUnidadMayorYUnidadMenor(unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornar la respuesta
	c.JSON(http.StatusOK, gin.H{
		"totalRegistros": count,
		"totalRUNs":      totalRUNs,
	})
}

// 8.1
// ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos obtiene las unidades menores con profesores activos (PAF).
func (h *EstadisticasController) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(c *gin.Context) {
	// Obtener el parámetro unidadMayor y semestre desde la URL
	unidadMayor := c.Param("unidadMayor")
	semestre := c.Param("semestre") // Capturar el parámetro semestre

	// Validar que ambos parámetros no estén vacíos
	if unidadMayor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio
	resultado, err := h.Service.ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(unidadMayor, semestre)
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
	// Obtener los parámetros 'unidadMayor' y 'semestre' desde la URL
	unidadMayor := c.Param("unidadMayor")
	semestre := c.Param("semestre")

	// Validar que se hayan recibido ambos parámetros
	if unidadMayor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio para obtener las unidades mayores con profesores activos filtrados por 'unidadMayor' y 'semestre'
	unidadesInactivas, err := ctrl.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor(unidadMayor, semestre)
	if err != nil {
		// Si hubo un error, devolver una respuesta con el mensaje de error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver la respuesta con el mapa de unidades mayores y el total de profesores
	c.JSON(http.StatusOK, gin.H{"unidades": unidadesInactivas})
}

// 8.2
// ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3 maneja la solicitud para obtener unidades menores sin profesores en Pipelsoft (versión 8.3).
func (ctrl *EstadisticasController) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(c *gin.Context) {
	// Obtener los parámetros 'unidadMayor' y 'semestre' desde la URL
	unidadMayor := c.Param("unidadMayor")
	semestre := c.Param("semestre")

	// Validar que se hayan recibido ambos parámetros
	if unidadMayor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio para obtener las unidades menores sin profesores en Pipelsoft, filtrados por 'unidadMayor' y 'semestre'
	unidadesSinProfesores, err := ctrl.Service.ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(unidadMayor, semestre)
	if err != nil {
		// Si hubo un error, devolver una respuesta con el mensaje de error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devolver la respuesta con el mapa de unidades menores y el total de profesores
	c.JSON(http.StatusOK, gin.H{"unidades": unidadesSinProfesores})
}

// 8.4
// ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos maneja la solicitud para obtener unidades menores con profesores filtrados por PAF activos.
func (c *EstadisticasController) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos(ctx *gin.Context) {
	// Recuperamos los parámetros 'unidadMayor' y 'semestre' desde la ruta
	unidadMayor := ctx.Param("unidadMayor")
	semestre := ctx.Param("semestre")

	// Validar que se hayan recibido ambos parámetros
	if unidadMayor == "" || semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor' y 'semestre' son obligatorios"})
		return
	}

	// Llamamos al servicio para obtener los datos, pasando los parámetros 'unidadMayor' y 'semestre'
	resultado, err := c.Service.ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos(unidadMayor, semestre)
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
	// Obtener los parámetros codEstadoPAF, unidadMayor y semestre desde la URL
	codEstadoPAF := c.Param("codEstadoPAF")
	unidadMayor := c.Param("unidadMayor")
	semestre := c.Param("semestre")

	// Validar que se hayan recibido todos los parámetros
	if codEstadoPAF == "" || unidadMayor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros codEstadoPAF, unidadMayor y semestre son obligatorios"})
		return
	}

	// Llamar al servicio con los parámetros obtenidos
	resultado, err := h.Service.ObtenerUnidadesMenoresPorCodEstadoPAF(codEstadoPAF, unidadMayor, semestre)
	if err != nil {
		// Si ocurre un error, devolvemos una respuesta 400 con el mensaje de error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

func (h *EstadisticasController) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor9_1(c *gin.Context) {
	// Obtener los parámetros desde la URL
	unidadMayor := c.Param("unidadMayor") // Parámetro de la unidad mayor
	unidadMenor := c.Param("unidadMenor") // Parámetro de la unidad menor (opcional)
	semestre := c.Param("semestre")       // Parámetro del semestre

	// Validar que el parámetro unidadMayor y semestre no estén vacíos
	if unidadMayor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros unidadMayor y semestre son obligatorios"})
		return
	}

	// Llamar al servicio para obtener las unidades menores con profesores
	unidadesMenores, err := h.Service.ObtenerUnidadesMenoresConProfesoresPorUnidadMayor9_1(unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, unidadesMenores)
}

// 9.2
func (h *EstadisticasController) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_9_2(c *gin.Context) {
	// Obtener los parámetros desde la URL
	unidadMayor := c.Param("unidadMayor") // Parámetro de la unidad mayor
	unidadMenor := c.Param("unidadMenor") // Parámetro de la unidad menor
	semestre := c.Param("semestre")       // Parámetro del semestre

	// Validar que ambos parámetros no estén vacíos
	if unidadMayor == "" || unidadMenor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio con los parámetros obtenidos
	resultado, err := h.Service.ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_9_2(unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// 9.3
func (h *EstadisticasController) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayorYUnidadMenor9_3(c *gin.Context) {
	// Obtener los parámetros desde la URL
	unidadMayor := c.Param("unidadMayor") // Parámetro de la unidad mayor
	unidadMenor := c.Param("unidadMenor") // Parámetro de la unidad menor
	semestre := c.Param("semestre")       // Parámetro del semestre

	// Validar que los parámetros no estén vacíos
	if unidadMayor == "" || unidadMenor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio con los parámetros obtenidos
	resultado, err := h.Service.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayorYUnidadMenor9_3(unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// 9.4
func (h *EstadisticasController) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivosPorUnidadMayorYUnidadMenor9_4(c *gin.Context) {
	// Obtener los parámetros desde la URL
	unidadMayor := c.Param("unidadMayor") // Parámetro de la unidad mayor
	unidadMenor := c.Param("unidadMenor") // Parámetro de la unidad menor
	semestre := c.Param("semestre")       // Parámetro del semestre

	// Validar que los parámetros no estén vacíos
	if unidadMayor == "" || unidadMenor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio con los parámetros obtenidos
	resultado, err := h.Service.ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivosPorUnidadMayorYUnidadMenor9_4(unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}

// 9.5
func (h *EstadisticasController) ObtenerUnidadesMenoresPorCodEstadoPAFPorCodEstadoYUnidadMayorYUnidadMenor9_5(c *gin.Context) {
	// Obtener los parámetros desde la URL
	codEstadoPAF := c.Param("codEstadoPAF") // Parámetro del estado del PAF
	unidadMayor := c.Param("unidadMayor")   // Parámetro de la unidad mayor
	unidadMenor := c.Param("unidadMenor")   // Parámetro de la unidad menor
	semestre := c.Param("semestre")         // Parámetro del semestre

	// Validar que todos los parámetros no estén vacíos
	if codEstadoPAF == "" || unidadMayor == "" || unidadMenor == "" || semestre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Los parámetros 'codEstadoPAF', 'unidadMayor', 'unidadMenor' y 'semestre' son obligatorios"})
		return
	}

	// Llamar al servicio con los parámetros obtenidos
	resultado, err := h.Service.ObtenerUnidadesMenoresPorCodEstadoPAFPorCodEstadoYUnidadMayorYUnidadMenor9_5(codEstadoPAF, unidadMayor, unidadMenor, semestre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responder con los datos en formato JSON
	c.JSON(http.StatusOK, resultado)
}
