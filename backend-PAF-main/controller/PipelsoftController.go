package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
)

type PipelsoftController struct {
	PipelsoftService *service.PipelsoftService
}

// NewPipelsoftController crea un nuevo controlador que maneja las rutas relacionadas con Pipelsoft.
func NewPipelsoftController(pipelsoftService *service.PipelsoftService) *PipelsoftController {
	return &PipelsoftController{
		PipelsoftService: pipelsoftService,
	}
}

// ObtenerContratosUltimos7Dias maneja la solicitud para obtener contratos de los últimos 7 días.
func (pc *PipelsoftController) ObtenerContratosUltimos7Dias(w http.ResponseWriter, r *http.Request) {
	contratos, err := pc.PipelsoftService.ObtenerContratosUltimos7Dias()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// ObtenerContratosPorCodigoAsignatura maneja la solicitud para obtener contratos por código de asignatura.
func (pc *PipelsoftController) ObtenerContratosPorCodigoAsignatura(w http.ResponseWriter, r *http.Request) {
	codigoAsignatura := r.URL.Query().Get("codigoAsignatura")
	if codigoAsignatura == "" {
		http.Error(w, "El parámetro 'codigoAsignatura' es obligatorio", http.StatusBadRequest)
		return
	}

	contratos, err := pc.PipelsoftService.ObtenerContratosPorCodigoAsignatura(codigoAsignatura)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// ObtenerContratosUltimoMes maneja la solicitud para obtener contratos del último mes.
func (pc *PipelsoftController) ObtenerContratosUltimoMes(w http.ResponseWriter, r *http.Request) {
	contratos, err := pc.PipelsoftService.ObtenerContratosUltimoMes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contratos)
}

// ObtenerPersonaPorCorreo maneja la solicitud para obtener una persona por correo.
func (pc *PipelsoftController) ObtenerPersonaPorCorreo(w http.ResponseWriter, r *http.Request) {
	correo := r.URL.Query().Get("correo")
	if correo == "" {
		http.Error(w, "El parámetro 'correo' es obligatorio", http.StatusBadRequest)
		return
	}

	persona, err := pc.PipelsoftService.ObtenerPersonaPorCorreo(correo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if persona == nil {
		http.Error(w, "Persona no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persona)
}

<<<<<<< HEAD
// ObtenerUnidadPorCodigo maneja la solicitud para obtener una unidad por código.
func (pc *PipelsoftController) ObtenerUnidadPorCodigo(w http.ResponseWriter, r *http.Request) {
	codigo := r.URL.Query().Get("codigo")
	if codigo == "" {
		http.Error(w, "El parámetro 'codigo' es obligatorio", http.StatusBadRequest)
		return
	}

	unidad, err := pc.PipelsoftService.ObtenerUnidadPorCodigo(codigo)
=======

func (c *PipelsoftController) ObtenerUnidadPorCodigo(w http.ResponseWriter, r *http.Request) {
	codigo := mux.Vars(r)["codigo"]
	pipelsoft, err := c.Service.ObtenerUnidadPorCodigo(codigo)
>>>>>>> a7d672ab4f0028fefa78de043e7169f61c75e505
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if unidad == nil {
		http.Error(w, "Unidad no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unidad)
}

// ObtenerListaPersonas maneja la solicitud para obtener la lista de personas.
func (pc *PipelsoftController) ObtenerListaPersonas(w http.ResponseWriter, r *http.Request) {
	personas, err := pc.PipelsoftService.ObtenerListaPersonas()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(personas)
}

// ObtenerEstadisticas maneja la solicitud para obtener las estadísticas calculadas.
func (pc *PipelsoftController) ObtenerEstadisticas(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para calcular las estadísticas
	estadisticas, err := pc.PipelsoftService.CalcularEstadisticas()
	if err != nil {
		http.Error(w, "No se pudieron calcular las estadísticas: "+err.Error(), http.StatusInternalServerError)
		return
	}

<<<<<<< HEAD
	// Preparar la respuesta en formato JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estadisticas)
}
=======
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
>>>>>>> a7d672ab4f0028fefa78de043e7169f61c75e505
