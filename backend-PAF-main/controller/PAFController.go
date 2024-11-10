package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

// PAFController define las operaciones HTTP sobre PAF
type PAFController struct {
	PafService *service.PAFService
}


