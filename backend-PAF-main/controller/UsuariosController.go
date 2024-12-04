package controller

import (
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-gonic/gin"
)

type UsuariosController struct {
	UsuariosService *service.UsuariosService
}

// Nuevo controlador para usuarios
func NewUsuariosController(usuariosService *service.UsuariosService) *UsuariosController {
	return &UsuariosController{UsuariosService: usuariosService}
}

// Endpoint para obtener un usuario por su Run
func (uc *UsuariosController) GetUsuarioByRun(c *gin.Context) {
	// Obtener el parámetro "run" de la ruta
	run := c.Param("run")
	if run == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'run' es obligatorio"})
		return
	}

	// Llamar al servicio para obtener el usuario
	usuario, err := uc.UsuariosService.GetUsuarioByRun(run)
	if err != nil {
		if err.Error() == "usuario no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el usuario"})
		}
		return
	}

	// Responder con la información del usuario
	c.JSON(http.StatusOK, usuario)
}
