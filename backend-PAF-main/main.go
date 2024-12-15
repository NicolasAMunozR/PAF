package main

import (
	"fmt"
	"log"

	"github.com/NicolasAMunozR/PAF/backend-PAF/DB"
	"github.com/NicolasAMunozR/PAF/backend-PAF/controller"
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/NicolasAMunozR/PAF/backend-PAF/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3" // Añadido para el cron job
)

// Función principal
func main() {
	// Crear el enrutador Gin
	r := gin.Default()

	// Configurar CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Si necesitas enviar cookies
	}))

	// Conectar a la base de datos
	DB.InitDBConnections()

	// Instanciar el servicio y controlador de HistorialPafAceptadas
	historialPafAceptadasService := service.NewHistorialPafAceptadasService(DB.DBPersonal)
	historialPafAceptadasController := controller.HistorialPafAceptadasController{
		Service: historialPafAceptadasService,
	}

	// Rutas para el controlador HistorialPafAceptadas
	r.POST("/api/paf-en-linea/historial/post/:codigoPAF/:cod_asignatura_pipelsoft/:comentario", historialPafAceptadasController.CrearHistorialHandler)
	r.GET("/api/paf-en-linea/historial", historialPafAceptadasController.ObtenerTodosLosHistorialesHandler)
	r.DELETE("/api/paf-en-linea/historial/:codigo_paf", historialPafAceptadasController.EliminarHistorialHandler)
	r.PUT("/api/paf-en-linea/historial/:codigoPAF/actualizarBanderaAceptacion", historialPafAceptadasController.ActualizarBanderaAceptacion)

	// Servicios y controladores adicionales
	pipelsoftService := service.NewPipelsoftService(DB.DBPersonal)
	pipelsoftController := controller.NewPipelsoftController(pipelsoftService)

	// Rutas para el controlador de Pipelsoft (actualizadas)
	r.GET("/api/paf-en-linea/pipelsoft/contratos-curso/:codigo_curso", pipelsoftController.ObtenerContratosPorCodigoCurso)
	r.GET("/api/paf-en-linea/pipelsoft/contratos", pipelsoftController.ObtenerTodosLosContratos)
	r.GET("/api/paf-en-linea/pipelsoft/contratos-run/:run", pipelsoftController.ObtenerContratosPorRUN)
	r.GET("/api/paf-en-linea/pipelsoft/contratos-nombreUnidadMenor/:nombreUnidadMenor", pipelsoftController.ObtenerContratosPorNombreUnidadMenor)
	r.GET("/api/paf-en-linea/pipelsoft/contratos-nombreUnidadMayor/:nombreUnidadMayor", pipelsoftController.ObtenerContratosPorNombreUnidadMayor)
	r.GET("/api/paf-en-linea/contratos/codigo_paf/:codigo_paf", pipelsoftController.ObtenerPorCodigoPAF)
	r.GET("/api/paf-en-linea/contratos/ultimos_7_dias", pipelsoftController.ObtenerPAFUltimos7Dias)
	r.GET("/api/paf-en-linea/contratos/ultimo_mes", pipelsoftController.ObtenerPAFUltimoMes)
	//controaldor igual a contratos-run pero muestra todos aunque esten aceptados
	r.GET("/api/paf-en-linea/pipelsoft/obtenerContratos/mostrarTodo/:rut", pipelsoftController.ObtenerContratosPorRUNMostrarTodo)

	// Ruta para obtener unidades menores por unidad mayor
	r.GET("/api/paf-en-linea/pipelsoft/unidades-menores", pipelsoftController.ObtenerUnidadesMenores)
	// Instanciar el servicio y controlador de Horarios
	horarioService := service.NewHorarioService(DB.DBPersonal)
	horarioController := controller.NewHorarioController(horarioService)

	// Ruta para obtener los horarios por Run
	r.GET("/api/paf-en-linea/horarios/:run", horarioController.ObtenerHorariosPorRun)

	// Instanciar el servicio y controlador de ProfesorDB
	profesorDBService := service.NewProfesorDBService(DB.DBPersonal)
	profesorDBController := controller.NewProfesorDBController(*profesorDBService)
	r.GET("/api/paf-en-linea/profesorDB/:run", profesorDBController.ObtenerProfesorDBPorRun)

	// Instanciar el servicio y controlador de Estadísticas
	estadisticasService := service.NewEstadisticasService(DB.DBPersonal)
	estadisticasController := controller.NewEstadisticasController(estadisticasService)

	// Ruta para obtener las estadísticas
	///api/paf-en-linea
	r.GET("/api/paf-en-linea/estadisticas/:semestre", estadisticasController.ObtenerEstadisticas)
	r.GET("/api/paf-en-linea/estadisticas/unidad/:nombreUnidadMayor/:semestre", estadisticasController.ContarRegistrosPorUnidadMayor)
	r.GET("/api/paf-en-linea/estadisticas/frecuencia-unidades-mayores/:semestre", estadisticasController.ObtenerFrecuenciaNombreUnidadMayor)
	r.GET("/api/paf-en-linea/estadisticas/PafActivas/:semestre", estadisticasController.ContarRegistrosPorCodEstado)
	r.GET("/api/paf-en-linea/estadisticas/pafActivas/unidad-mayor/:unidadMayor/:semestre", estadisticasController.ObtenerPafActivasPorUnidadHandler)
	r.GET("/api/paf-en-linea/estadisticas/unidad-mayor/:unidad-mayor/:semestre", estadisticasController.ObtenerEstadisticasPorUnidadMayorHandler)
	r.GET("/api/paf-en-linea/estadisticas/unidad-mayor/unidades-menores-frecuencia/:unidad-mayor/semestre", estadisticasController.ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayorHandler)

	// Inicializar servicios y controladores
	contratoService := service.NewContratoService(DB.DBPersonal)
	contratoController := controller.NewContratoController(contratoService)

	contrato := r.Group("/api/paf-en-linea/contratos")
	{
		contrato.Use(cors.Default())
		contrato.GET("/", contratoController.GetAllContratosHandler)
		contrato.GET("/:run", contratoController.GetContratoByRunHandler)
		contrato.GET("/unidad-mayor/:unidad", contratoController.GetPafByUnidadMayorHandler)

	}

	// Crear el servicio y controlador
	usuariosService := service.NewUsuariosService(DB.DBPersonal)
	usuariosController := controller.NewUsuariosController(usuariosService)
	// Ruta para obtener un usuario por su Run
	r.GET("/api/paf-en-linea/usuario/rut/:run", usuariosController.GetUsuarioByRun)

	// nuevo y nueva url
	// Crear servicio y controlador
	historialService := service.NewHistorialPasosPafService(DB.DBPersonal)
	historialController := controller.NewHistorialPasosPafController(historialService)

	r.GET("/api/paf-en-linea/historialPaso/:id_paf/:run_docente", historialController.ObtenerHistorialYDuracionesPorIdYRun)

	//controlador para obtener los profesores que no poseen paf
	r.GET("/api/paf-en-linea/estadisticas/obtener-y-comparar-runs", estadisticasController.ObtenerYCompararRunsHandler)

	// 1
	r.GET("/api/paf-en-linea/estadisticas/unidades-mayores/cant_profesores/:semestre", estadisticasController.ObtenerUnidadesMayoresHandler)
	// 2
	r.GET("/api/paf-en-linea/estadisticas/unidades-mayores/sin_profesores/:semestre", estadisticasController.ObtenerUnidadesMayoresSinProfesoresEnPipelsoftHandler)
	// 3
	r.GET("/api/paf-en-linea/estadisticas/unidades-mayores/profesores-filtrados/:semestre", estadisticasController.ObtenerUnidadesMayoresConProfesoresFiltradosHandler)
	// 4
	r.GET("/api/paf-en-linea/estadisticas/unidades-mayores/profesores-codestado/:semestre", estadisticasController.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivos)
	// 5
	r.GET("/api/paf-en-linea/estadisticas/profesores/estado/:codEstadoPAF/:semestre", estadisticasController.ObtenerUnidadesMayoresPorCodEstadoPAF)
	// 6 arreglar
	r.GET("/api/paf-en-linea/estadisticas/6/:unidadMayor/:unidadMenor/semestre", estadisticasController.ObtenerEstadisticasPorUnidad)
	// 7 arreglar
	r.GET("/api/paf-en-linea/estadisticas/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ContarRegistrosPorUnidadMayorYUnidadMenor)

	// Ruta para obtener las unidades menores con profesores filtrados (PAF activos)
	// 8.1
	r.GET("/api/paf-en-linea/estadistica/unidades-menores-con-profesores-activos/8_1/:unidadMayor/:semestre", estadisticasController.ObtenerUnidadesMenoresConProfesoresPorUnidadMayor)
	// 8.3
	r.GET("/api/paf-en-linea/estadistica/unidades-menores-sin-profesores/8_3/:unidadMayor/:semestre", estadisticasController.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor)

	// 8.2
	// Ruta para obtener unidades menores sin profesores en Pipelsoft (8.3)
	r.GET("/api/paf-en-linea/estadistica/unidades-menores-sin-profesores-8-2/:unidadMayor/:semestre", estadisticasController.ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3)

	// 8.4
	r.GET("/api/paf-en-linea/estadistica/unidades-menores-con-profesores-paf-activos/8_4/:unidadMayor/:semestre", estadisticasController.ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos)

	// 8.5
	r.GET("/api/paf-en-linea/estadistica/unidades-menores/:codEstadoPAF/:unidadMayor/:semestre", estadisticasController.ObtenerUnidadesMenoresPorCodEstadoPAF)
	// 9.1
	r.GET("/api/paf-en-linea/unidadesmenores/profesores/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ObtenerUnidadesMenoresConProfesoresPorUnidadMayor9_1)
	// 9.2
	r.GET("/api/paf-en-linea/unidadesmenores/sinprofesores/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_9_2)
	// 9.3
	r.GET("/api/paf-en-linea/unidadesmayores/filtradospafactivos/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayorYUnidadMenor9_3)
	// 9.4
	r.GET("/api/paf-en-linea/unidadesmenores/filtradospafactivos/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivosPorUnidadMayorYUnidadMenor9_4)
	// 9.5
	r.GET("/api/paf-en-linea/unidadesmenores/porcodestadopaf/:codEstadoPAF/:unidadMayor/:unidadMenor/:semestre", estadisticasController.ObtenerUnidadesMenoresPorCodEstadoPAFPorCodEstadoYUnidadMayorYUnidadMenor9_5)

	// Iniciar el cron job para actualización periódica
	// de acuerdo al sai profes sin paf y sin contrato sin contar profesores repetidos, desde pipelsoft estan los ruts con un 0 al inicio y con digito verificador
	// luego se busca en la tabla de contratos, y hay revisamos cuales ruts corresponden,
	go iniciarCronJob()

	// Iniciar el servidor Gin
	log.Println("Servidor escuchando en el puerto 3000...")
	if err := r.Run(":3000"); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

// Función para iniciar el cron job de actualización
func iniciarCronJob() {
	c := cron.New()

	// Ejecutar cada 30 minutos
	c.AddFunc("@every 45m", func() {
		actualizarModificaciones()
	})
	c.Start()

	// Mantener el cron job corriendo en segundo plano
	select {}
}

// Función que compara y actualiza los registros
func actualizarModificaciones() {
	// Conexión a la base de datos
	db := DB.DBPersonal

	var historial []models.HistorialPafAceptadas

	// Obtener todos los registros de HistorialPafAceptadas
	if err := db.Find(&historial).Error; err != nil {
		log.Println("Error al obtener historial:", err)
		return
	}

	for _, h := range historial {
		var pipelsoft models.Pipelsoft // Inicializa la variable limpia para evitar conflictos.

		// Buscar el registro correspondiente en Pipelsoft filtrando por id_paf, codigo_asignatura y fecha_inicio_contrato
		if err := db.Where("id_paf = ? AND codigo_asignatura = ? AND fecha_inicio_contrato = ?", h.IdPaf, h.CodigoAsignatura, h.FechaInicioContrato).Take(&pipelsoft).Error; err != nil {
			// Si no se encuentra en Pipelsoft, marcamos como eliminada
			if err := db.Model(&h).Updates(map[string]interface{}{
				"codigo_modificacion":      1, // Se marca como modificada (en este caso, eliminada)
				"bandera_modificacion":     2, // 2 = Eliminada
				"descripcion_modificacion": fmt.Sprintf("PAF con id_paf: %d, cod_asignatura: %s y fecha_inicio_contrato: %s eliminada, no se encuentra en Pipelsoft", h.IdPaf, h.CodigoAsignatura, h.FechaInicioContrato),
			}).Error; err != nil {
				log.Println("Error al actualizar HistorialPafAceptadas como eliminada:", err)
			} else {
				fmt.Printf("PAF eliminada, no encontrada en Pipelsoft: id_paf %d, cod_asignatura %s, fecha_inicio_contrato %s\n", h.IdPaf, h.CodigoAsignatura, h.FechaInicioContrato)
			}
			continue // Si no se encuentra, seguimos con el siguiente registro
		}

		// Si se encuentra el registro, puedes continuar con otras operaciones aquí

		// Variable para almacenar los cambios detectados
		var cambios []string

		// Comparar los valores y detectar qué cambió
		if h.FechaInicioContrato != pipelsoft.FechaInicioContrato {
			cambios = append(cambios, fmt.Sprintf("FechaInicioContrato cambiado de %s a %s", h.FechaInicioContrato, pipelsoft.FechaInicioContrato))
		}
		if h.FechaFinContrato != pipelsoft.FechaFinContrato {
			cambios = append(cambios, fmt.Sprintf("FechaFinContrato cambiado de %s a %s", h.FechaFinContrato, pipelsoft.FechaFinContrato))
		}
		if h.CodigoAsignatura != pipelsoft.CodigoAsignatura {
			cambios = append(cambios, fmt.Sprintf("CodigoAsignatura cambiado de %s a %s", h.CodigoAsignatura, pipelsoft.CodigoAsignatura))
		}
		if h.NombreAsignatura != pipelsoft.NombreAsignatura {
			cambios = append(cambios, fmt.Sprintf("NombreAsignatura cambiado de %s a %s", h.NombreAsignatura, pipelsoft.NombreAsignatura))
		}
		if h.CantidadHoras != pipelsoft.HorasAsignatura {
			cambios = append(cambios, fmt.Sprintf("CantidadHoras cambiadas de %d a %d", h.CantidadHoras, pipelsoft.HorasAsignatura))
		}
		if h.Jerarquia != pipelsoft.Jerarquia {
			cambios = append(cambios, fmt.Sprintf("Jerarquia cambiada de %s a %s", h.Jerarquia, pipelsoft.Jerarquia))
		}

		// Si hay cambios detectados, se marca como modificado
		if len(cambios) > 0 {
			// Descripción de la modificación basada en los cambios
			descripcion := fmt.Sprintf("Modificación detectada: %v", cambios)

			// Actualizamos el campo CodigoModificacion a 1 (modificado)
			// También actualizamos la descripción de la modificación
			if err := db.Model(&h).Updates(map[string]interface{}{
				"codigo_modificacion":      1,
				"descripcion_modificacion": descripcion,
				"bandera_modificacion":     1, // 1 = Modificado
			}).Error; err != nil {
				log.Println("Error al actualizar HistorialPafAceptadas:", err)
			} else {
				fmt.Printf("Registro actualizado para CodigoPAF: %d\n", h.IdPaf)
			}
		}
	}

	fmt.Println("Modificaciones verificadas y actualizadas.")
}
