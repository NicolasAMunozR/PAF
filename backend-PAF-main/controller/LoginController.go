package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
)

// LoginController maneja el inicio de sesión
func LoginController(c *gin.Context) {
	var loginRequest models.Login

	// Vincular el JSON recibido a la estructura Login
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Llamar a la función DoRequestLogin para autenticar al usuario
	responseLogin, err := models.DoRequestLogin(loginRequest.User, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Responder con el token de autenticación
	c.JSON(http.StatusOK, gin.H{
		"token":  responseLogin.Token,
		"expire": responseLogin.Expire,
		"data":   responseLogin.Data,
	})
}
