package controller

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"

	"github.com/gorilla/mux"
)

type ContratoController struct {
	Service *service.ContratoService
}

// GetAllContratosHandler maneja las solicitudes para obtener todos los contratos.
func (c *ContratoController) GetAllContratosHandler(w http.ResponseWriter, r *http.Request) {
	contratos, err := c.Service.GetAllContratos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// GetContratoByRunHandler maneja las solicitudes para obtener un contrato por el RUN del docente.
func (c *ContratoController) GetContratoByRunHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	run, ok := vars["run"]
	if !ok {
		http.Error(w, "El parámetro 'run' es obligatorio", http.StatusBadRequest)
		return
	}

	contrato, err := c.Service.GetContratoByRun(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if contrato == nil {
		http.Error(w, "Contrato no encontrado", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contrato)
}

// GetContratosByUnidadMayorHandler maneja las solicitudes para obtener contratos por unidad mayor.
func (c *ContratoController) GetContratosByUnidadMayorHandler(w http.ResponseWriter, r *http.Request) {
	unidad := r.URL.Query().Get("unidad")
	if unidad == "" {
		http.Error(w, "El parámetro 'unidad' es obligatorio", http.StatusBadRequest)
		return
	}

	contratos, err := c.Service.GetContratosByUnidadMayor(unidad)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}
