package main

import (
	"log"
	"net/http"

	"github.com/NicolasAMunozR/PAF/backend-PAF/DB"
	"github.com/NicolasAMunozR/PAF/backend-PAF/controller"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Configuración de CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3001"}, // Asegúrate de que el puerto coincida con el de tu frontend Nuxt
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Crear el enrutador y aplicar CORS
	r := mux.NewRouter()
	r.Use(c.Handler)
	// Conectar a la base de datos
	DB.DBconnection()

	// Instanciar el servicio y controlador de HistorialPafAceptadas
	historialPafAceptadasService := service.NewHistorialPafAceptadasService(DB.DB)
	historialPafAceptadasController := controller.HistorialPafAceptadasController{
		Service: historialPafAceptadasService,
	}
	// Rutas para el controlador HistorialPafAceptadas
	r.HandleFunc("/historial", historialPafAceptadasController.CrearHistorialHandler).Methods("POST")
	r.HandleFunc("/historial", historialPafAceptadasController.ObtenerTodosLosHistorialesHandler).Methods("GET")
	r.HandleFunc("/historial/{codigo_paf}", historialPafAceptadasController.EliminarHistorialHandler).Methods("DELETE")

	pipelsoftService := service.NewPipelsoftService(DB.DB)
	pipelsoftController := controller.NewPipelsoftController(pipelsoftService)

	r.HandleFunc("/pipelsoft/contratos/ultimos7dias", pipelsoftController.ObtenerContratosUltimos7Dias).Methods("GET")
	r.HandleFunc("/pipelsoft/contratos/codigo/{codigoAsignatura}", pipelsoftController.ObtenerContratosPorCodigoAsignatura).Methods("GET")
	r.HandleFunc("/pipelsoft/contratos/ultimomes", pipelsoftController.ObtenerContratosUltimoMes).Methods("GET")
	r.HandleFunc("/pipelsoft/persona/correo/{correo}", pipelsoftController.ObtenerPersonaPorCorreo).Methods("GET")
	r.HandleFunc("/pipelsoft/persona/rut/{run}", pipelsoftController.ObtenerPersonaPorRUT).Methods("GET")
	r.HandleFunc("/pipelsoft/personas/rut/{run}", pipelsoftController.ObtenerPersonasPorRUT).Methods("GET")
	r.HandleFunc("/pipelsoft/proceso/estado/{estado}", pipelsoftController.ObtenerProcesoPorEstado).Methods("GET")
	r.HandleFunc("/pipelsoft/unidad/codigo/{codigo}", pipelsoftController.ObtenerUnidadPorCodigo).Methods("GET")
	r.HandleFunc("/pipelsoft/persona", pipelsoftController.ObtenerListaPersonas).Methods("GET")

	// Iniciar el servidor
	log.Println("Servidor escuchando en el puerto 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
