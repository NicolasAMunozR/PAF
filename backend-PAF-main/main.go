package main

import (
	"log"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/DB"
	"github.com/NicolasAMunozR/PAF/backend-PAF/controller"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
)

func main() {
	// Conectar a la base de datos
	DB.DBconnection()


	// Configuramos el enrutador
	r := mux.NewRouter()


	// Iniciar el servidor
	log.Println("Servidor escuchando en el puerto 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
