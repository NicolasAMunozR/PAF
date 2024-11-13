package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type ContratoController struct {
	ContratoService *service.ContratoService
}

// ObtenerContratoPorID maneja la obtención de un contrato por su ID
func (c *ContratoController) ObtenerContratoPorID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	contrato, err := c.ContratoService.ObtenerContratoPorID(uint(id))
	if err != nil {
		http.Error(w, "Contrato no encontrado", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(contrato)
}

// ObtenerTodosContratos maneja la obtención de todos los contratos
func (c *ContratoController) ObtenerTodosContratos(w http.ResponseWriter, r *http.Request) {
	contratos, err := c.ContratoService.ObtenerTodosContratos()
	if err != nil {
		http.Error(w, "Error al obtener los contratos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contratos)
}

// ObtenerContratosUltimos7Dias maneja la obtención de los contratos creados en los últimos 7 días
func (c *ContratoController) ObtenerContratosUltimos7Dias(w http.ResponseWriter, r *http.Request) {
	contratos, err := c.ContratoService.ObtenerContratosUltimos7Dias()
	if err != nil {
		http.Error(w, "Error al obtener los contratos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contratos)
}

// ObtenerContratosPorCodigoAsignatura maneja la obtención de contratos con el mismo código de asignatura
func (c *ContratoController) ObtenerContratosPorCodigoAsignatura(w http.ResponseWriter, r *http.Request) {
	codigoAsignatura := mux.Vars(r)["codigoAsignatura"]
	contratos, err := c.ContratoService.ObtenerContratosPorCodigoAsignatura(codigoAsignatura)
	if err != nil {
		http.Error(w, "Error al obtener los contratos", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contratos)
}
