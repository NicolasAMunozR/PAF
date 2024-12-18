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

// Endpoint para obtener una lista de usuarios por su Run
func (uc *UsuariosController) GetUsuariosByRun(c *gin.Context) {
	// Obtener el parámetro "run" de la ruta
	run := c.Param("run")
	if run == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El parámetro 'run' es obligatorio"})
		return
	}

	// Llamar al servicio para obtener la lista de usuarios
	usuarios, err := uc.UsuariosService.GetUsuariosByRun(run)
	if err != nil {
		if err.Error() == "no se encontraron usuarios con el RUN proporcionado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los usuarios"})
		}
		return
	}

	// Responder con la lista de usuarios
	c.JSON(http.StatusOK, usuarios)
}
