package controller

import (
	"encoding/json"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

type ProfesorDBController struct {
	Service service.ProfesorDBService
}

func NewProfesorDBController(service service.ProfesorDBService) *ProfesorDBController {
	return &ProfesorDBController{Service: service}
}

func (P *ProfesorDBController) ObtenerProfesorDBPorRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	run := vars["run"]

	profesor, err := P.Service.ObtenerProfesorPorRUT(run)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(profesor)
}