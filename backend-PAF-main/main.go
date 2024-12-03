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
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
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
	r.POST("/historial/post/:codigoPAF", historialPafAceptadasController.CrearHistorialHandler)
	r.GET("/historial", historialPafAceptadasController.ObtenerTodosLosHistorialesHandler)
	r.DELETE("/historial/:codigo_paf", historialPafAceptadasController.EliminarHistorialHandler)
	r.PUT("/historial/:codigoPAF/actualizarBanderaAceptacion", historialPafAceptadasController.ActualizarBanderaAceptacion)

	// Servicios y controladores adicionales
	pipelsoftService := service.NewPipelsoftService(DB.DBPersonal)
	pipelsoftController := controller.NewPipelsoftController(pipelsoftService)

	// Rutas para el controlador de Pipelsoft (actualizadas)
	r.GET("/pipelsoft/contratos-curso/:codigo_curso", pipelsoftController.ObtenerContratosPorCodigoCurso)
	r.GET("/pipelsoft/contratos", pipelsoftController.ObtenerTodosLosContratos)
	r.GET("/pipelsoft/contratos-run/:run", pipelsoftController.ObtenerContratosPorRUN)
	r.GET("/pipelsoft/contratos-nombreUnidadMenor/:nombreUnidadMenor", pipelsoftController.ObtenerContratosPorNombreUnidadMenor)
	r.GET("/pipelsoft/contratos-nombreUnidadMayor/:nombreUnidadMayor", pipelsoftController.ObtenerContratosPorNombreUnidadMayor)
	r.GET("/contratos/codigo_paf/:codigo_paf", pipelsoftController.ObtenerPorCodigoPAF)
	r.GET("/contratos/ultimos_7_dias", pipelsoftController.ObtenerPAFUltimos7Dias)
	r.GET("/contratos/ultimo_mes", pipelsoftController.ObtenerPAFUltimoMes)
	// Ruta para obtener unidades menores por unidad mayor
	r.GET("/pipelsoft/unidades-menores", pipelsoftController.ObtenerUnidadesMenores)
	// Instanciar el servicio y controlador de Horarios
	horarioService := service.NewHorarioService(DB.DBPersonal)
	horarioController := controller.NewHorarioController(horarioService)

	// Ruta para obtener los horarios por Run
	r.GET("/horarios/:run", horarioController.ObtenerHorariosPorRun)

	// Instanciar el servicio y controlador de ProfesorDB
	profesorDBService := service.NewProfesorDBService(DB.DBPersonal)
	profesorDBController := controller.NewProfesorDBController(*profesorDBService)
	r.GET("/profesorDB/:run", profesorDBController.ObtenerProfesorDBPorRun)

	// Instanciar el servicio y controlador de Estadísticas
	estadisticasService := service.NewEstadisticasService(DB.DBPersonal)
	estadisticasController := controller.NewEstadisticasController(estadisticasService)

	// Ruta para obtener las estadísticas
	r.GET("/estadisticas", estadisticasController.ObtenerEstadisticas)
	r.GET("/estadisticas/unidad/:nombreUnidadMayor", estadisticasController.ContarRegistrosPorUnidadMayor)
	r.GET("/estadisticas/frecuencia-unidades-mayores", estadisticasController.ObtenerFrecuenciaNombreUnidadMayor)
	r.GET("/estadisticas/PafActivas", estadisticasController.ContarRegistrosPorCodEstado)
	r.GET("/estadisticas/unidad-mayor/:unidad-mayor", estadisticasController.ObtenerEstadisticasPorUnidadMayorHandler)

	// Inicializar servicios y controladores
	contratoService := service.NewContratoService(DB.DBPersonal)
	contratoController := controller.NewContratoController(contratoService)

	contrato := r.Group("/contratos")
	{
		contrato.Use(cors.Default())
		contrato.GET("/", contratoController.GetAllContratosHandler)
		contrato.GET("/:run", contratoController.GetContratoByRunHandler)
		contrato.GET("/unidad-mayor/:unidad", contratoController.GetPafByUnidadMayorHandler)

	}

	// Iniciar el cron job para actualización periódica
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
	c.AddFunc("@every 1m", func() {
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

		// Buscar el registro correspondiente en Pipelsoft
		if err := db.Where("id_paf = ?", h.IdPaf).Take(&pipelsoft).Error; err != nil {
			// Si no se encuentra en Pipelsoft, marcamos como eliminada
			if err := db.Model(&h).Updates(map[string]interface{}{
				"codigo_modificacion":      1, // Se marca como modificada (en este caso, eliminada)
				"bandera_modificacion":     2, // 2 = Eliminada
				"descripcion_modificacion": fmt.Sprintf("PAF con CodigoPAF: %d eliminada, no se encuentra en Pipelsoft", h.IdPaf),
			}).Error; err != nil {
				log.Println("Error al actualizar HistorialPafAceptadas como eliminada:", err)
			} else {
				fmt.Printf("PAF eliminada, no encontrada en Pipelsoft: id_paf %d\n", h.IdPaf)
			}
			continue // Si no se encuentra, seguimos con el siguiente registro
		}

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
