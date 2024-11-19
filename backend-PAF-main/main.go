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
	DB.InitDBConnections()

	// Instanciar el servicio y controlador de HistorialPafAceptadas
	historialPafAceptadasService := service.NewHistorialPafAceptadasService(DB.DBPersonal)
	historialPafAceptadasController := controller.HistorialPafAceptadasController{
		Service: historialPafAceptadasService,
	}

	// Rutas para el controlador HistorialPafAceptadas
	r.HandleFunc("/historial", historialPafAceptadasController.CrearHistorialHandler).Methods("POST")
	r.HandleFunc("/historial", historialPafAceptadasController.ObtenerTodosLosHistorialesHandler).Methods("GET")
	r.HandleFunc("/historial/{codigo_paf}", historialPafAceptadasController.EliminarHistorialHandler).Methods("DELETE")

	// Servicios y controladores adicionales
	pipelsoftService := service.NewPipelsoftService(DB.DBPersonal)
	pipelsoftController := controller.NewPipelsoftController(pipelsoftService)

	// Rutas para el controlador de Pipelsoft (actualizadas)
	r.HandleFunc("/pipelsoft/contratos-curso/{codigo_curso}", pipelsoftController.ObtenerContratosPorCodigoCurso).Methods("GET")
	r.HandleFunc("/pipelsoft/contratos", pipelsoftController.ObtenerTodosLosContratos).Methods("GET")
	r.HandleFunc("/pipelsoft/contratos-run/{run}", pipelsoftController.ObtenerContratosPorRUN).Methods("GET")
	r.HandleFunc("/contratos/codigo_paf/{codigo_paf}", pipelsoftController.ObtenerPorCodigoPAF).Methods("GET")
	r.HandleFunc("/contratos/ultimos_7_dias", pipelsoftController.ObtenerPAFUltimos7Dias).Methods("GET")
	r.HandleFunc("/contratos/ultimo_mes", pipelsoftController.ObtenerPAFUltimoMes).Methods("GET")

	// Instanciar el servicio y controlador de Horarios
	horarioService := service.NewHorarioService(DB.DBPersonal)
	horarioController := controller.NewHorarioController(horarioService)

	// Ruta para obtener los horarios por Run
	r.HandleFunc("/horarios/{run}", horarioController.ObtenerHorariosPorRun).Methods("GET")

	// Aplicar CORS al enrutador
	handler := c.Handler(r)

	// Iniciar el servidor
	log.Println("Servidor escuchando en el puerto 3000...")
	log.Fatal(http.ListenAndServe(":3000", handler))
}
