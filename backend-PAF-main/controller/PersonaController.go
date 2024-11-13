package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

// PersonaController maneja las solicitudes HTTP para Personas
type PersonaController struct {
	PersonaService *service.PersonaService
}

// NewPersonaController crea una nueva instancia de PersonaController
func NewPersonaController(personaService *service.PersonaService) *PersonaController {
	return &PersonaController{PersonaService: personaService}
}

// ObtenerPersonaPorID maneja la obtención de una Persona por ID
func (c *PersonaController) ObtenerPersonaPorID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	persona, err := c.PersonaService.ObtenerPersonaPorID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(persona)
}

// ObtenerTodasPersonas maneja la obtención de todas las Personas
func (c *PersonaController) ObtenerTodasPersonas(w http.ResponseWriter, r *http.Request) {
	personas, err := c.PersonaService.ObtenerTodasPersonas()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(personas)
}

// ObtenerPersonaPorCorreo maneja la obtención de una Persona por correo
func (c *PersonaController) ObtenerPersonaPorCorreo(w http.ResponseWriter, r *http.Request) {
	correo := mux.Vars(r)["correo"]

	persona, err := c.PersonaService.ObtenerPersonaPorCorreo(correo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(persona)
}

// ObtenerPersonaPorRUT maneja la obtención de una Persona por RUT
func (c *PersonaController) ObtenerPersonaPorRUT(w http.ResponseWriter, r *http.Request) {
	rut := mux.Vars(r)["run"]

	persona, err := c.PersonaService.ObtenerPersonaPorRUT(rut)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(persona)
}
