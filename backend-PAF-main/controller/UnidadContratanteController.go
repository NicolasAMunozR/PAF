package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type UnidadContratanteController struct {
	UnidadService *service.UnidadContratanteService
}

// NewUnidadContratanteController crea una nueva instancia de UnidadContratanteController
func NewUnidadContratanteController(service *service.UnidadContratanteService) *UnidadContratanteController {
	return &UnidadContratanteController{
		UnidadService: service,
	}
}

// ObtenerUnidadPorCodigo maneja la obtención de una unidad contratante por código
func (c *UnidadContratanteController) ObtenerUnidadPorCodigo(w http.ResponseWriter, r *http.Request) {
	codigo := mux.Vars(r)["codigo"]
	unidad, err := c.UnidadService.ObtenerUnidadPorCodigo(codigo)
	if err != nil {
		http.Error(w, "Unidad no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(unidad)
}

// ObtenerUnidadPorID maneja la obtención de una unidad contratante por ID
func (c *UnidadContratanteController) ObtenerUnidadPorID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	unidad, err := c.UnidadService.ObtenerUnidadPorID(uint(id))
	if err != nil {
		http.Error(w, "Unidad no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(unidad)
}

// ObtenerTodasUnidades maneja la obtención de todas las unidades contratantes
func (c *UnidadContratanteController) ObtenerTodasUnidades(w http.ResponseWriter, r *http.Request) {
	unidades, err := c.UnidadService.ObtenerTodasUnidades()
	if err != nil {
		http.Error(w, "Error al obtener unidades", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(unidades)
}
