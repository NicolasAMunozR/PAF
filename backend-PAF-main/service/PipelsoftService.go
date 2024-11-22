package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type PipelsoftService struct {
	DBPersonal *gorm.DB
}

// Constructor del servicio
func NewPipelsoftService(dbPersonal *gorm.DB) *PipelsoftService {
	return &PipelsoftService{
		DBPersonal: dbPersonal,
	}
}

func (s *PipelsoftService) getAcceptedCodes() (map[string]struct{}, error) {
	var historialAceptadas []struct {
		CodigoPaf string `gorm:"column:codigo_paf"`
	}

	// Modificación: agregamos la condición BanderaAceptacion = 1
	if err := s.DBPersonal.Table("historial_paf_aceptadas").
		Select("codigo_paf").
		Where("BanderaAceptacion = ?", 1). // Condición para filtrar por BanderaAceptacion = 1
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

// Combinar datos de Pipelsoft con datos de la tabla "profesor_dbs" por RUT
func (s *PipelsoftService) combinarDatosPorRUT(pipelsofts []models.Pipelsoft) ([]models.DatosCombinados, error) {
	runSet := make(map[string]struct{})
	for _, pipel := range pipelsofts {
		runSet[pipel.Run] = struct{}{}
	}

	var profesores []models.ProfesorDB
	runs := make([]string, 0, len(runSet))
	for run := range runSet {
		runs = append(runs, run)
	}

	if err := s.DBPersonal.Where("run IN ?", runs).Find(&profesores).Error; err != nil {
		return nil, err
	}

	profesoresMap := make(map[string]models.ProfesorDB)
	for _, profesor := range profesores {
		profesoresMap[profesor.RUN] = profesor
	}

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

// Obtener contratos por código de curso
func (s *PipelsoftService) ObtenerContratosPorCodigoCurso(codigoCurso string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("codigo_asignatura= ?", codigoCurso).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

// Obtener todos los contratos
func (s *PipelsoftService) ObtenerTodosLosContratos() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

// Obtener contratos por profesor (RUN)
func (s *PipelsoftService) ObtenerContratosPorRUN(run string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("run = ?", run).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

// Obtener contratos por Código PAF
func (s *PipelsoftService) ObtenerPorCodigoPAF(codigoPAF string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("codigo_paf = ?", codigoPAF).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

// Obtener PAF de los últimos 7 días
func (s *PipelsoftService) ObtenerPAFUltimos7Dias() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft
	hace7Dias := time.Now().AddDate(0, 0, -7)

	if err := s.DBPersonal.Where("fecha_inicio_contrato >= ?", hace7Dias).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}

// Obtener PAF del último mes
func (s *PipelsoftService) ObtenerPAFUltimoMes() ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft
	haceUnMes := time.Now().AddDate(0, -1, 0)

	if err := s.DBPersonal.Where("fecha_inicio_contrato >= ?", haceUnMes).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.combinarDatosPorRUT(contratosFiltrados)
}




// Crear una nueva estructura de Pipelsoft sin guardar nada en la base de datos
func (s *PipelsoftService) crearEstructuraPipelsoft() models.Pipelsoft {
	// Crear una nueva instancia del modelo Pipelsoft con valores vacíos o predeterminados
	pipelsoft := models.Pipelsoft{
		Run:                    "",  // valor vacío, puede ser ajustado
		Nombres:                "",  // valor vacío
		PrimerApellido:         "",  // valor vacío
		SegundoApellido:        "",  // valor vacío
		Correo:                 "",  // valor vacío
		CodigoUnidadContratante: "",  // valor vacío
		NombreUnidadContratante: "",  // valor vacío
		NombreUnidadMayor:      "",  // valor vacío
		CodigoPAF:              "",  // valor vacío
		FechaInicioContrato:    "",  // valor vacío
		FechaFinContrato:       "",  // valor vacío
		CodigoAsignatura:       "",  // valor vacío
		NombreAsignatura:       "",  // valor vacío
		CantidadHoras:          0,   // valor predeterminado (0)
		Jerarquia:              "",  // valor vacío
		Calidad:                "",  // valor vacío
		EstadoProceso:          0,   // valor predeterminado (0)
		FechaUltimaModificacionProceso: time.Time{}, // valor vacío (cero)
	}

	return pipelsoft
}