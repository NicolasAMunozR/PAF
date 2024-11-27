// controllers/historial_paf_aceptadas_controller.go
package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type HistorialPafAceptadasController struct {
	Service *service.HistorialPafAceptadasService
}

func (h *HistorialPafAceptadasController) CrearHistorialHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el código PAF desde los parámetros de la URL
	vars := mux.Vars(r)
	codigoPAFStr, ok := vars["codigoPAF"]
	if !ok {
		http.Error(w, "El parámetro 'codigoPAF' es obligatorio", http.StatusBadRequest)
		return
	}

	// Convertir el código PAF a entero
	codigoPAF, err := strconv.Atoi(codigoPAFStr)
	if err != nil {
		http.Error(w, "El parámetro 'codigoPAF' debe ser un número entero válido", http.StatusBadRequest)
		return
	}

	// Crear un struct auxiliar para parsear el cuerpo de la solicitud
	type CrearHistorialRequest struct {
		Profesor models.ProfesorDB `json:"profesor"`
		Bloque   []string          `json:"bloque"`
	}

	var request CrearHistorialRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, fmt.Sprintf("Error al parsear el cuerpo de la solicitud: %v", err), http.StatusBadRequest)
		return
	}

	// Llamar al servicio para crear el historial
	historial, err := h.Service.CrearHistorial(codigoPAF, request.Profesor, request.Bloque)
	if err != nil {
		log.Printf("Error al crear el historial: %v\n", err)
		http.Error(w, fmt.Sprintf("Error al crear el historial: %v", err), http.StatusInternalServerError)
		return
	}

	// Responder con el historial creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(historial); err != nil {
		log.Printf("Error al codificar la respuesta: %v\n", err)
		http.Error(w, "Error al generar la respuesta", http.StatusInternalServerError)
	}
}

// ObtenerTodosLosHistorialesHandler maneja la solicitud para obtener todos los historiales
func (c *HistorialPafAceptadasController) ObtenerTodosLosHistorialesHandler(w http.ResponseWriter, r *http.Request) {
	historiales, err := c.Service.ObtenerTodosLosHistoriales()
	if err != nil {
		http.Error(w, "Error al obtener los historiales", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historiales)
}

// EliminarHistorialHandler maneja la solicitud para eliminar un historial por CodigoPAF
func (c *HistorialPafAceptadasController) EliminarHistorialHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el CodigoPAF desde los parámetros de la URL
	codigoPAF := mux.Vars(r)["codigo_paf"]

	// Llamar al servicio para eliminar el historial por CodigoPAF
	if err := c.Service.EliminarHistorial(codigoPAF); err != nil {
		http.Error(w, "Error al eliminar el historial", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Responder con 204 No Content si la eliminación fue exitosa
}
