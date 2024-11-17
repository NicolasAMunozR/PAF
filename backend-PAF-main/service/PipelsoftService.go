package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type PipelsoftService struct {
	DBPipelsoft *gorm.DB
	DBPersonal  *gorm.DB
	DBProfesor  *gorm.DB // Nueva conexión a la base de datos de profesores
}

// NewPipelsoftService crea un nuevo servicio que utiliza las tres bases de datos.
func NewPipelsoftService(dbPipelsoft, dbPersonal, dbProfesor *gorm.DB) *PipelsoftService {
	return &PipelsoftService{
		DBPipelsoft: dbPipelsoft,
		DBPersonal:  dbPersonal,
		DBProfesor:  dbProfesor, // Añadido para la base de datos de profesores
	}
}

// Obtener códigos aceptados de la base de datos "personal"
func (s *PipelsoftService) getAcceptedCodes() (map[string]struct{}, error) {
	var historialAceptadas []struct {
		CodigoPaf string `gorm:"column:codigo_paf"`
	}
	if err := s.DBPersonal.Table("historial_paf_aceptadas").
		Select("codigo_paf").
		Find(&historialAceptadas).Error; err != nil {
		return nil, err
	}

	aceptadosMap := make(map[string]struct{})
	for _, item := range historialAceptadas {
		aceptadosMap[item.CodigoPaf] = struct{}{}
	}
	return aceptadosMap, nil
}

// Filtrar registros que no están en los códigos aceptados
func (s *PipelsoftService) filterRecords(records []models.Pipelsoft, aceptadosMap map[string]struct{}) []models.Pipelsoft {
	var filtered []models.Pipelsoft
	for _, record := range records {
		if _, found := aceptadosMap[record.CodigoPAF]; !found {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

// Combina los datos de Pipelsoft y ProfesorDB por RUT
func (s *PipelsoftService) combinarDatosPorRUT(pipelsofts []models.Pipelsoft) ([]models.DatosCombinados, error) {
	// Obtener todos los RUN únicos de los registros de Pipelsoft
	runSet := make(map[string]struct{})
	for _, pipel := range pipelsofts {
		runSet[pipel.Run] = struct{}{}
	}

	// Consultar la tabla de Profesores (profesor_dbs) en la base de datos "profesor"
	var profesores []models.ProfesorDB
	runs := make([]string, 0, len(runSet))
	for run := range runSet {
		runs = append(runs, run)
	}

	if err := s.DBProfesor.Where("run IN ?", runs).Find(&profesores).Error; err != nil {
		return nil, err
	}

	// Crear un mapa para buscar profesores por RUN
	profesoresMap := make(map[string]models.ProfesorDB)
	for _, profesor := range profesores {
		profesoresMap[profesor.RUN] = profesor
	}

	// Combinar los datos de Pipelsoft y ProfesorDB
	var datosCombinados []models.DatosCombinados
	for _, pipel := range pipelsofts {
		dato := models.DatosCombinados{
			PipelsoftData: pipel,
		}
		if profesor, found := profesoresMap[pipel.Run]; found {
			dato.ProfesorData = profesor
		}
		datosCombinados = append(datosCombinados, dato)
	}

	return datosCombinados, nil
}

// Obtener contratos de los últimos 7 días de la base de datos "pipelsoft"
func (s *PipelsoftService) ObtenerContratosUltimos7Dias() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// Obtener los registros de la base de datos "pipelsoft" que fueron creados en los últimos 7 días
	if err := s.DBPipelsoft.Where("created_at >= ?", sevenDaysAgo).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	// Obtener los códigos aceptados
	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	// Filtrar registros que no están en los códigos aceptados
	prueba := s.filterRecords(pipelsofts, aceptadosMap)

	// Combinar los datos con la tabla de profesores por RUT
	return s.combinarDatosPorRUT(prueba)
}

func (s *PipelsoftService) ObtenerContratosPorCodigoAsignatura(codigoAsignatura string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPipelsoft.Where("codigo_asignatura = ?", codigoAsignatura).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	// Filtrar los registros que no están aceptados y combinar los datos por RUT
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

func (s *PipelsoftService) ObtenerContratosUltimoMes() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	if err := s.DBPipelsoft.Where("fecha_inicio_contrato >= ?", oneMonthAgo).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	// Filtrar los registros y combinar los datos por RUT
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

func (s *PipelsoftService) ObtenerPersonaPorCorreo(correo string) (*models.DatosCombinados, error) {
	var pipelsoft models.Pipelsoft

	if err := s.DBPipelsoft.Where("correo = ?", correo).First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	// Obtener los códigos aceptados
	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	if _, found := aceptadosMap[pipelsoft.CodigoPAF]; found {
		return nil, nil
	}

	// Combinar los datos con la tabla de profesores por RUT
	datosCombinados, err := s.combinarDatosPorRUT([]models.Pipelsoft{pipelsoft})
	if err != nil {
		return nil, err
	}

	// Retornar los datos combinados para la persona
	return &datosCombinados[0], nil
}

func (s *PipelsoftService) ObtenerUnidadPorCodigo(codigo string) (*models.DatosCombinados, error) {
	var pipelsoft models.Pipelsoft

	if err := s.DBPipelsoft.Where("codigo_unidad_contratante = ?", codigo).First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	// Obtener los códigos aceptados
	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	if _, found := aceptadosMap[pipelsoft.CodigoPAF]; found {
		return nil, nil
	}

	// Combinar los datos con la tabla de profesores por RUT
	datosCombinados, err := s.combinarDatosPorRUT([]models.Pipelsoft{pipelsoft})
	if err != nil {
		return nil, err
	}

	// Retornar los datos combinados para la unidad
	return &datosCombinados[0], nil
}

func (s *PipelsoftService) ObtenerListaPersonas() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPipelsoft.Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	// Obtener los códigos aceptados
	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}

	// Filtrar los registros que no están en los códigos aceptados
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	// Combinar los datos con la tabla de profesores por RUT
	return s.combinarDatosPorRUT(contratosFiltrados)
}

// CalcularEstadisticas calcula estadísticas sobre la lista de DatosCombinados.
func (s *PipelsoftService) CalcularEstadisticas() (map[string]interface{}, error) {
	// Obtener la lista de personas combinadas desde el método ObtenerListaPersonas
	listaPersonas, err := s.ObtenerListaPersonas()
	if err != nil {
		return nil, err
	}

	// Inicializar variables para las estadísticas
	estadisticas := make(map[string]interface{})
	var totalContratos int
	var contratosActivos int
	var contratosInactivos int
	asignaturas := make(map[string]int)

	// Iterar sobre la lista de DatosCombinados y calcular estadísticas
	for _, datos := range listaPersonas {
		// Contar contratos activos e inactivos
		if datos.PipelsoftData.EstadoProceso == "activo" {
			contratosActivos++
		} else if datos.PipelsoftData.EstadoProceso == "inactivo" {
			contratosInactivos++
		}

		// Contar la cantidad de contratos por asignatura (si existe un campo "codigo_asignatura")
		asignaturas[datos.PipelsoftData.CodigoAsignatura]++

		// Contar total de contratos
		totalContratos++
	}

	// Calcular promedios

	// Almacenar las estadísticas calculadas
	estadisticas["total_contratos"] = totalContratos
	estadisticas["contratos_activos"] = contratosActivos
	estadisticas["contratos_inactivos"] = contratosInactivos
	estadisticas["asignaturas"] = asignaturas

	return estadisticas, nil
}

// ObtenerPersonaPorPaf obtiene un registro de `Pipelsoft` por su código PAF
func (s *PipelsoftService) ObtenerPersonaPorPaf(codigoPaf string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("codigo_paf = ?", codigoPaf).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	return &pipelsoft, nil
}