package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

// ProcesoController estructura el controlador para los procesos
type ProcesoController struct {
	ProcesoService *service.ProcesoService
}

// NewProcesoController crea una nueva instancia de ProcesoController
func NewProcesoController(service *service.ProcesoService) *ProcesoController {
	return &ProcesoController{
		ProcesoService: service,
	}
}

// ObtenerProcesoPorID maneja la solicitud para obtener un proceso por ID
func (c *ProcesoController) ObtenerProcesoPorID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	proceso, err := c.ProcesoService.ObtenerProcesoPorID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proceso)
}

// ObtenerTodosLosProcesos maneja la solicitud para obtener todos los procesos
func (c *ProcesoController) ObtenerTodosLosProcesos(w http.ResponseWriter, r *http.Request) {
	procesos, err := c.ProcesoService.ObtenerTodosLosProcesos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procesos)
}

// ObtenerProcesoPorEstado maneja la solicitud para obtener procesos por estado
func (c *ProcesoController) ObtenerProcesoPorEstado(w http.ResponseWriter, r *http.Request) {
	estado := mux.Vars(r)["estado"]

	procesos, err := c.ProcesoService.ObtenerProcesoPorEstado(estado)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procesos)
}
