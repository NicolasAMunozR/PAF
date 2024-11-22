package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"

	// Cambiar con la ruta del paquete del servicio
	"github.com/gorilla/mux"
)

// Controlador que maneja las estadísticas
type EstadisticasController struct {
	Service *service.EstadisticasService
}

// Nueva instancia de EstadisticasController
func NewEstadisticasController(service *service.EstadisticasService) *EstadisticasController {
	return &EstadisticasController{Service: service}
}

// ObtenerEstadisticas maneja la solicitud para obtener las estadísticas
func (c *EstadisticasController) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	// Obtener las estadísticas del servicio
	estadisticas, err := c.Service.ObtenerEstadisticas()
	if err != nil {
		// Si ocurre un error, devolver respuesta con el código de error
		http.Error(w, fmt.Sprintf("Error al obtener las estadísticas: %v", err), http.StatusInternalServerError)
		return
	}

	// Configurar el encabezado de la respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Convertir la respuesta a JSON
	if err := json.NewEncoder(w).Encode(estadisticas); err != nil {
		http.Error(w, fmt.Sprintf("Error al convertir la respuesta a JSON: %v", err), http.StatusInternalServerError)
	}
}

// Configurar las rutas del controlador
func (c *EstadisticasController) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/estadisticas", c.ObtenerEstadisticas).Methods("GET")
}

// Nuevo controlador para actualizar la BanderaAceptacion
func (ctrl *HistorialPafAceptadasController) ActualizarBanderaAceptacion(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros del request
	vars := mux.Vars(r)
	codigoPAF := vars["codigoPAF"] // Se espera el código PAF en la URL

	var requestBody struct {
		NuevaBanderaAceptacion int `json:"nuevaBanderaAceptacion"` // Nueva bandera de aceptación en el cuerpo de la petición
	}

	// Parsear el cuerpo de la solicitud
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para actualizar la BanderaAceptacion
	err = ctrl.HistorialService.ActualizarBanderaAceptacion(codigoPAF, requestBody.NuevaBanderaAceptacion)
	if err != nil {
		log.Printf("Error al actualizar BanderaAceptacion: %v", err)
		http.Error(w, fmt.Sprintf("Error al actualizar BanderaAceptacion: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "BanderaAceptacion actualizada correctamente"}`))
}
