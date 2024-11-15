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
		AllowedOrigins:   []string{"http://localhost:3001"},                   // Tu frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Agregar OPTIONS si lo necesitas
		AllowedHeaders:   []string{"Content-Type", "Authorization"},           // Agregar Authorization si es necesario
		ExposedHeaders:   []string{"X-Total-Count"},
		AllowCredentials: true,
		Debug:            true, // Agregar para más información sobre los errores de CORS
	})

	// Crear el enrutador
	r := mux.NewRouter()

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

	// Servicios y controladores adicionales
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
	r.HandleFunc("/pipelsoft/persona/paf/{codigoPaf}", pipelsoftController.ObtenerPersonaPorPaf).Methods("GET")


	// Instanciar el servicio y controlador de Horarios
	horarioService := service.NewHorarioService(DB.DB)
	horarioController := controller.NewHorarioController(horarioService)

	// Ruta para obtener los horarios por Run
	r.HandleFunc("/horarios/{run}", horarioController.ObtenerHorariosPorRun).Methods("GET")

	// Aplicar CORS al enrutador
	handler := c.Handler(r)

	// Iniciar el servidor
	log.Println("Servidor escuchando en el puerto 3000...")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
