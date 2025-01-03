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
