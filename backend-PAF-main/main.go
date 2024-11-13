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
	// Creamos una instancia del servicio
	pafService := service.NewPAFService()

	// Creamos una instancia del controlador
	pafController := controller.NewPAFController(pafService)

	// Crear el servicio y controlador de Persona
	personaService := service.NewPersonaService(DB.DB)
	personaController := controller.NewPersonaController(personaService)

	// Crear el servicio y el controlador para UnidadContratante
	unidadService := service.NewUnidadContratanteService(DB.DB)
	unidadController := controller.NewUnidadContratanteController(unidadService)

	// Crear el servicio y el controlador de contrato
	contratoService := service.NewContratoService(DB.DB)
	contratoController := controller.NewContratoController(contratoService)

	// Inicializar los servicios y controladores
	procesoService := service.NewProcesoService(DB.DB)
	procesoController := controller.NewProcesoController(procesoService)

	// Configuramos el enrutador
	r := mux.NewRouter()

	// Definir las rutas del controlador
	r.HandleFunc("/paf", pafController.CrearPAF).Methods("POST")
	//Filtro Por Id PAF
	r.HandleFunc("/paf/{id:[0-9]+}", pafController.ObtenerPAF).Methods("GET")
	//Obtiene todas las PAFS
	r.HandleFunc("/pafs", pafController.ObtenerTodosPAFs).Methods("GET")
	//Actualiza la PAF del id ingresado
	r.HandleFunc("/paf/{id:[0-9]+}", pafController.ActualizarPAF).Methods("PUT")
	//Elimina una PAF
	r.HandleFunc("/paf/{id:[0-9]+}", pafController.EliminarPAF).Methods("DELETE")

	// Ruta para filtrar por nombre del profesor
	r.HandleFunc("/pafs/buscarNombre", pafController.ObtenerPAFsPorNombreProfesor).Methods("GET")

	// Rutas para los endpoints GET de Persona
	r.HandleFunc("/persona/{id:[0-9]+}", personaController.ObtenerPersonaPorID).Methods("GET")
	r.HandleFunc("/personas", personaController.ObtenerTodasPersonas).Methods("GET")
	r.HandleFunc("/persona/correo/{correo}", personaController.ObtenerPersonaPorCorreo).Methods("GET")
	r.HandleFunc("/persona/rut/{run}", personaController.ObtenerPersonaPorRUT).Methods("GET")

	// Definir las rutas del controlador UnidadContratante
	r.HandleFunc("/unidad/{id:[0-9]+}", unidadController.ObtenerUnidadPorID).Methods("GET")
	r.HandleFunc("/unidad/codigo/{codigo}", unidadController.ObtenerUnidadPorCodigo).Methods("GET")
	r.HandleFunc("/unidades", unidadController.ObtenerTodasUnidades).Methods("GET")

	// Rutas para Contrato
	r.HandleFunc("/contrato/{id:[0-9]+}", contratoController.ObtenerContratoPorID).Methods("GET")
	r.HandleFunc("/contratos", contratoController.ObtenerTodosContratos).Methods("GET")
	r.HandleFunc("/contratos/ultimos7dias", contratoController.ObtenerContratosUltimos7Dias).Methods("GET")
	r.HandleFunc("/contratos/codigo/{codigoAsignatura}", contratoController.ObtenerContratosPorCodigoAsignatura).Methods("GET")
	r.HandleFunc("/contratos/ultimo-mes", contratoController.ObtenerContratosUltimoMes).Methods("GET")

	r.HandleFunc("/procesos", procesoController.ObtenerTodosLosProcesos).Methods("GET")
	r.HandleFunc("/proceso/{id:[0-9]+}", procesoController.ObtenerProcesoPorID).Methods("GET")
	r.HandleFunc("/procesos/estado/{estado}", procesoController.ObtenerProcesoPorEstado).Methods("GET")

	// Iniciar el servidor
	log.Println("Servidor escuchando en el puerto 3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
