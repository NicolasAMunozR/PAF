// controllers/profesor_db_controller.go
package controller

import (
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

type ProfesorDBController struct {
	Service service.ProfesorDBService
}

// Constructor del controlador
func NewProfesorDBController(service service.ProfesorDBService) *ProfesorDBController {
	return &ProfesorDBController{Service: service}
}

// Obtener Profesor por RUN
func (P *ProfesorDBController) ObtenerProfesorDBPorRun(ctx *gin.Context) {
	run := ctx.Param("run")

	profesor, err := P.Service.ObtenerProfesorPorRUT(run)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, profesor)
}

// Obtener la cantidad de profesores que no están en Pipelsoft
func (ctrl *ProfesorDBController) GetCountProfesoresNotInPipelsoft(ctx *gin.Context) {
	count, err := ctrl.Service.GetCountProfesoresNotInPipelsoft()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la cantidad de profesores que no están en Pipelsoft"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

func (c *ProfesorDBController) GetProfesoresNoContrato(ctx *gin.Context) {
	// Obtener el semestre desde la URL
	semestre := ctx.Param("semestre")

	// Validar que el semestre no esté vacío
	if semestre == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'semestre' es requerido"})
		return
	}

	// Llamar al servicio con el semestre como argumento
	resultado, err := c.Service.ObtenerProfesoresSinContratoYNoAcademico(semestre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enviar la respuesta con los profesores filtrados
	ctx.JSON(http.StatusOK, resultado)
}
