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
	"github.com/robfig/cron"
)

// Inicializar servicios y controladores
func setupRoutes(r *gin.Engine) {
	// Historial PAF Aceptadas
	historialService := service.NewHistorialPafAceptadasService(DB.DBPersonal)
	historialController := controller.HistorialPafAceptadasController{
		Service: historialService,
	}
	historialRoutes := r.Group("/api/paf-en-linea/historial")
	{
		historialRoutes.POST("/post/:codigoPAF/:cod_asignatura_pipelsoft/:comentario", historialController.CrearHistorialHandler)
		historialRoutes.GET("", historialController.ObtenerTodosLosHistorialesHandler)
		historialRoutes.DELETE("/:codigo_paf", historialController.EliminarHistorialHandler)
		historialRoutes.PUT("/:codigoPAF/actualizarBanderaAceptacion", historialController.ActualizarBanderaAceptacion)
	}

	// Pipelsoft
	pipelsoftService := service.NewPipelsoftService(DB.DBPersonal)
	pipelsoftController := controller.NewPipelsoftController(pipelsoftService)
	pipelsoftRoutes := r.Group("/api/paf-en-linea/pipelsoft")
	{
		pipelsoftRoutes.GET("/contratos-curso/:codigo_curso", pipelsoftController.ObtenerContratosPorCodigoCurso)
		pipelsoftRoutes.GET("/contratos", pipelsoftController.ObtenerTodosLosContratos)
		pipelsoftRoutes.GET("/contratos-run/:run", pipelsoftController.ObtenerContratosPorRUN)
		pipelsoftRoutes.GET("/contratos-nombreUnidadMenor/:nombreUnidadMenor", pipelsoftController.ObtenerContratosPorNombreUnidadMenor)
		pipelsoftRoutes.GET("/contratos-nombreUnidadMayor/:nombreUnidadMayor", pipelsoftController.ObtenerContratosPorNombreUnidadMayor)
		pipelsoftRoutes.GET("/obtenerContratos/mostrarTodo/:rut", pipelsoftController.ObtenerContratosPorRUNMostrarTodo)
		pipelsoftRoutes.GET("/unidades-menores", pipelsoftController.ObtenerUnidadesMenores)
	}

	// Horarios
	horarioService := service.NewHorarioService(DB.DBPersonal)
	horarioController := controller.NewHorarioController(horarioService)
	r.GET("/api/paf-en-linea/horarios/:run", horarioController.ObtenerHorariosPorRun)

	// ProfesorDB
	profesorDBService := service.NewProfesorDBService(DB.DBPersonal)
	profesorDBController := controller.NewProfesorDBController(*profesorDBService)
	r.GET("/api/paf-en-linea/profesorDB/:run", profesorDBController.ObtenerProfesorDBPorRun)

	// Estadísticas
	estadisticasService := service.NewEstadisticasService(DB.DBPersonal)
	estadisticasController := controller.NewEstadisticasController(estadisticasService)
	estadisticasRoutes := r.Group("/api/paf-en-linea/estadisticas")
	{
		estadisticasRoutes.GET("/:semestreId", estadisticasController.ObtenerEstadisticas)
		estadisticasRoutes.GET("/unidad/:nombreUnidadMayor/:semestre", estadisticasController.ContarRegistrosPorUnidadMayor)
		estadisticasRoutes.GET("/frecuencia-unidades-mayores/:semestre", estadisticasController.ObtenerFrecuenciaNombreUnidadMayor)
		// Resto de rutas de estadísticas...
	}

	// Contratos
	contratoService := service.NewContratoService(DB.DBPersonal)
	contratoController := controller.NewContratoController(contratoService)
	contratosRoutes := r.Group("/api/paf-en-linea/contratos")
	{
		contratosRoutes.GET("/", contratoController.GetAllContratosHandler)
		contratosRoutes.GET("/:run", contratoController.GetContratoByRunHandler)
		contratosRoutes.GET("/unidad-mayor/:unidad", contratoController.GetPafByUnidadMayorHandler)
	}

	// Usuarios
	usuariosService := service.NewUsuariosService(DB.DBPersonal)
	usuariosController := controller.NewUsuariosController(usuariosService)
	r.GET("/api/paf-en-linea/usuario/rut/:run", usuariosController.GetUsuarioByRun)
}

// Middleware común
func commonMiddleware(c *gin.Context) {
	// Aquí puedes agregar validaciones, logging, etc.
	log.Printf("Ruta solicitada: %s", c.Request.URL.Path)
	c.Next()
}

// Función principal
func main() {
	// Crear el enrutador Gin
	r := gin.Default()

	// Configurar middleware global
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	r.Use(commonMiddleware)

	// Conectar a la base de datos
	DB.InitDBConnections()

	// Configurar rutas
	setupRoutes(r)

	// Iniciar el cron job para actualización periódica
	// de acuerdo al sai profes sin paf y sin contrato sin contar profesores repetidos, desde pipelsoft estan los ruts con un 0 al inicio y con digito verificador
	// luego se busca en la tabla de contratos, y hay revisamos cuales ruts corresponden,
	go iniciarCronJob()

	// Iniciar el servidor
	log.Println("Servidor iniciado en el puerto 8080")
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %s", err)
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
