// controller/pipelsoft_controller.go
package controller

import (
	"encoding/json"
	"fmt"
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

func (c *PipelsoftController) ObtenerPersonasPorRUT(w http.ResponseWriter, r *http.Request) {
	run := mux.Vars(r)["run"]
	pipelsofts, err := c.Service.ObtenerPersonasPorRUT(run)
	if err != nil {
		http.Error(w, "Error al obtener registros", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pipelsofts)
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

func (c *PipelsoftController) ObtenerPersonaPorPaf(w http.ResponseWriter, r *http.Request) {
	codigoPaf := mux.Vars(r)["codigoPaf"]

	// Llamada al servicio para obtener los datos
	pipelsoft, err := c.Service.ObtenerPersonaPorPaf(codigoPaf)
	if err != nil {
		// Verifica si el error es de "registro no encontrado"
		if err.Error() == fmt.Sprintf("registro con codigo_paf %s no encontrado", codigoPaf) {
			http.Error(w, fmt.Sprintf("Persona con código PAF %s no encontrada", codigoPaf), http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener los datos de la persona", http.StatusInternalServerError)
		}
		return
	}

	// Codifica el objeto de la persona a JSON
	w.Header().Set("Content-Type", "application/json")
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