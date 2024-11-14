// controller/pipelsoft_controller.go
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

func NewPipelsoftController(service *service.PipelsoftService) *PipelsoftController {
	return &PipelsoftController{Service: service}
}

func (c *PipelsoftController) ObtenerContratosUltimos7Dias(w http.ResponseWriter, r *http.Request) {
	pipelsofts, err := c.Service.ObtenerContratosUltimos7Dias()
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
}

func (c *PipelsoftController) ObtenerContratosPorCodigoAsignatura(w http.ResponseWriter, r *http.Request) {
	codigoAsignatura := mux.Vars(r)["codigoAsignatura"]
	pipelsofts, err := c.Service.ObtenerContratosPorCodigoAsignatura(codigoAsignatura)
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
}

func (c *PipelsoftController) ObtenerContratosUltimoMes(w http.ResponseWriter, r *http.Request) {
	pipelsofts, err := c.Service.ObtenerContratosUltimoMes()
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
}

func (c *PipelsoftController) ObtenerPersonaPorCorreo(w http.ResponseWriter, r *http.Request) {
	correo := mux.Vars(r)["correo"]
	pipelsoft, err := c.Service.ObtenerPersonaPorCorreo(correo)
	if err != nil {
		http.Error(w, "Persona no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pipelsoft)
}

func (c *PipelsoftController) ObtenerPersonaPorRUT(w http.ResponseWriter, r *http.Request) {
	run := mux.Vars(r)["run"]
	pipelsoft, err := c.Service.ObtenerPersonaPorRUT(run)
	if err != nil {
		http.Error(w, "Persona no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pipelsoft)
}

func (c *PipelsoftController) ObtenerProcesoPorEstado(w http.ResponseWriter, r *http.Request) {
	estado := mux.Vars(r)["estado"]
	pipelsofts, err := c.Service.ObtenerProcesoPorEstado(estado)
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
}

func (c *PipelsoftController) ObtenerUnidadPorCodigo(w http.ResponseWriter, r *http.Request) {
	codigo := mux.Vars(r)["codigo"]
	pipelsoft, err := c.Service.ObtenerUnidadPorCodigo(codigo)
	if err != nil {
		http.Error(w, "Unidad no encontrada", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(pipelsoft)
}

func (c *PipelsoftController) ObtenerListaPersonas(w http.ResponseWriter, r *http.Request) {
	pipelsofts, err := c.Service.ObtenerListaPersonas()
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
}
