package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service" // Cambiar según la ruta correcta del paquete
	"github.com/gorilla/mux"
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
func (c *EstadisticasController) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para obtener las estadísticas
	estadisticas, err := c.Service.ObtenerEstadisticas()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener estadísticas: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con las estadísticas en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(estadisticas); err != nil {
		http.Error(w, fmt.Sprintf("Error al codificar la respuesta: %v", err), http.StatusInternalServerError)
	}
}

// ContarRegistrosPorUnidadContratante maneja la solicitud para contar registros por nombre de unidad contratante
func (c *EstadisticasController) ContarRegistrosPorUnidadContratante(w http.ResponseWriter, r *http.Request) {
	// Obtener el nombre de la unidad contratante desde los parámetros de la URL
	nombreUnidadContratante := mux.Vars(r)["nombreUnidadContratante"]

	// Llamar al servicio para contar los registros
	count, err := c.Service.ContarRegistrosPorNombreUnidadContratante(nombreUnidadContratante)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al contar registros: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con el conteo en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"nombreUnidadContratante": nombreUnidadContratante,
		"conteo":                  count,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error al codificar la respuesta: %v", err), http.StatusInternalServerError)
	}
}

// ContarRegistrosPorCodEstado maneja la solicitud para contar registros de Pipelsoft donde el `cod_estado` no sea "F1", "F9" ni "A9"
func (c *EstadisticasController) ContarRegistrosPorCodEstado(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para obtener el conteo y el porcentaje
	count, porcentaje, err := c.Service.ContarRegistrosExcluyendoEstados()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al contar registros excluyendo estados: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con el conteo y porcentaje en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"conteo":     count,
		"porcentaje": porcentaje,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error al codificar la respuesta en JSON: %v", err), http.StatusInternalServerError)
	}
}

//esto NO DEBERIA IR AQUI MOVER

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
	err = ctrl.Service.ActualizarBanderaAceptacion(codigoPAF, requestBody.NuevaBanderaAceptacion)
	if err != nil {

		http.Error(w, fmt.Sprintf("Error al actualizar BanderaAceptacion: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "BanderaAceptacion actualizada correctamente"}`))
}

// ObtenerFrecuenciaNombreUnidadMayor maneja la solicitud para obtener la frecuencia de NombreUnidadMayor
func (c *EstadisticasController) ObtenerFrecuenciaNombreUnidadMayor(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para obtener los datos
	frecuencia, err := c.Service.ObtenerFrecuenciaNombreUnidadMayor()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener la frecuencia de NombreUnidadMayor: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con los datos en formato JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(frecuencia); err != nil {
		http.Error(w, fmt.Sprintf("Error al codificar la respuesta en JSON: %v", err), http.StatusInternalServerError)
	}
}
