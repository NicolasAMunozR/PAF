// controllers/historial_paf_aceptadas_controller.go
package controller

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type HistorialPafAceptadasController struct {
	Service *service.HistorialPafAceptadasService
}

// CrearHistorialHandler maneja la solicitud de creación de un nuevo HistorialPafAceptadas
func (c *HistorialPafAceptadasController) CrearHistorialHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CodigoPAF string `json:"codigo_paf"`
	}

	// Decodificar el JSON de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para crear el historial
	historial, err := c.Service.CrearHistorial(input.CodigoPAF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Devolver el historial creado como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(historial)
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
