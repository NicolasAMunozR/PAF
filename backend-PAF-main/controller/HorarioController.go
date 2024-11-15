package controller

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type HorarioController struct {
	HorarioService *service.HorarioService
}

// NewHorarioController crea un nuevo controlador para los horarios
func NewHorarioController(horarioService *service.HorarioService) *HorarioController {
	return &HorarioController{HorarioService: horarioService}
}

// ObtenerHorariosPorRun maneja la solicitud para obtener todos los horarios por Run
func (h *HorarioController) ObtenerHorariosPorRun(w http.ResponseWriter, r *http.Request) {
	run := mux.Vars(r)["run"] // Obtener el parámetro "run" de la URL

	// Llamar al servicio para obtener los horarios
	horarios, err := h.HorarioService.ObtenerHorariosPorRun(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Si no se encuentran horarios
	if len(horarios) == 0 {
		http.Error(w, "No se encontraron horarios para este Run", http.StatusNotFound)
		return
	}

	// Retornar los horarios encontrados como JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(horarios); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
