package controller

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type PipelsoftController struct {
	Service *service.PipelsoftService
}

// Constructor del controlador
func NewPipelsoftController(service *service.PipelsoftService) *PipelsoftController {
	return &PipelsoftController{Service: service}
}

// Obtener contratos por código de curso
func (c *PipelsoftController) ObtenerContratosPorCodigoCurso(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codigoCurso := vars["codigo_curso"]

	contratos, err := c.Service.ObtenerContratosPorCodigoCurso(codigoCurso)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// Obtener todos los contratos
func (c *PipelsoftController) ObtenerTodosLosContratos(w http.ResponseWriter, r *http.Request) {
	contratos, err := c.Service.ObtenerTodosLosContratos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// Obtener contratos por RUN
func (c *PipelsoftController) ObtenerContratosPorRUN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	run := vars["run"]

	contratos, err := c.Service.ObtenerContratosPorRUN(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// Obtener contratos por Código PAF
func (c *PipelsoftController) ObtenerPorCodigoPAF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	codigoPAF := vars["codigo_paf"]

	datos, err := c.Service.ObtenerPorCodigoPAF(codigoPAF)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(datos)
}

// Obtener PAF de los últimos 7 días
func (c *PipelsoftController) ObtenerPAFUltimos7Dias(w http.ResponseWriter, r *http.Request) {
	datos, err := c.Service.ObtenerPAFUltimos7Dias()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(datos)
}

// Obtener PAF del último mes
func (c *PipelsoftController) ObtenerPAFUltimoMes(w http.ResponseWriter, r *http.Request) {
	datos, err := c.Service.ObtenerPAFUltimoMes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(datos)
}

// Obtener contratos por nombre de unidad mayor
func (c *PipelsoftController) ObtenerContratosPorNombreUnidadMayor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nombreUnidadMayor := vars["nombreUnidadMayor"]

	contratos, err := c.Service.ObtenerContratosPorNombreUnidadMayor(nombreUnidadMayor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// Obtener contratos por nombre de unidad contratante
func (c *PipelsoftController) ObtenerContratosPorNombreUnidadContratante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nombreUnidadContratante := vars["nombreUnidadContratante"]

	contratos, err := c.Service.ObtenerContratosPorNombreUnidadContratante(nombreUnidadContratante)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}
