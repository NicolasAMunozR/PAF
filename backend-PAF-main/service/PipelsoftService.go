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
		Where("bandera_aceptacion = ?", 1). // Condición para filtrar por BanderaAceptacion = 1
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

func (s *PipelsoftService) comprobarYCombinarDatosPorCodigoPAF(pipelsofts []models.Pipelsoft) ([]models.DatosCombinados, error) {
	// Crear un conjunto único de los códigos PAF presentes en los registros de Pipelsoft
	codigoPAFSet := make(map[string]struct{})
	for _, pipel := range pipelsofts {
		codigoPAFSet[pipel.CodigoPAF] = struct{}{}
	}

	// Crear una lista para almacenar los datos de HistorialPafAceptadas
	var historialPafAceptadas []models.HistorialPafAceptadas
	codigosPAF := make([]string, 0, len(codigoPAFSet))
	for codigoPAF := range codigoPAFSet {
		codigosPAF = append(codigosPAF, codigoPAF)
	}

	// Consultar los datos de HistorialPafAceptadas cuyo CódigoPAF esté en la lista de códigos
	if err := s.DBPersonal.Where("codigo_paf IN ?", codigosPAF).Find(&historialPafAceptadas).Error; err != nil {
		return nil, err
	}

	// Crear un map para almacenar los registros de HistorialPafAceptadas por CódigoPAF
	historialPafMap := make(map[string]models.HistorialPafAceptadas)
	for _, historial := range historialPafAceptadas {
		historialPafMap[historial.CodigoPAF] = historial
	}

	// Crear una lista para almacenar los datos combinados
	var datosCombinados []models.DatosCombinados
	for _, pipel := range pipelsofts {
		// Crear un nuevo objeto DatosCombinados solo con los datos de Pipelsoft
		dato := models.DatosCombinados{
			PipelsoftData: pipel,
		}

		// Buscar si el CódigoPAF del registro de Pipelsoft está en HistorialPafAceptadas
		if historial, found := historialPafMap[pipel.CodigoPAF]; found {
			// Si se encuentra, agregar los datos del historialPafAceptadas al objeto DatosCombinados
			dato.HistorialPafData = historial
		} else {
			// Si no se encuentra, dejar HistorialPafData vacío (nulo o una estructura vacía)
			dato.HistorialPafData = models.HistorialPafAceptadas{}
		}

		// Agregar el objeto DatosCombinados a la lista
		datosCombinados = append(datosCombinados, dato)
	}

	// Devolver los datos combinados
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
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

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
}

// Crear una nueva estructura de Pipelsoft sin guardar nada en la base de datos
func (s *PipelsoftService) crearEstructuraPipelsoft() models.Pipelsoft {
	// Crear una nueva instancia del modelo Pipelsoft con valores vacíos o predeterminados
	pipelsoft := models.Pipelsoft{
		Run:                            "",          // valor vacío, puede ser ajustado
		Nombres:                        "",          // valor vacío
		PrimerApellido:                 "",          // valor vacío
		SegundoApellido:                "",          // valor vacío
		Correo:                         "",          // valor vacío
		CodigoUnidadContratante:        "",          // valor vacío
		NombreUnidadContratante:        "",          // valor vacío
		NombreUnidadMayor:              "",          // valor vacío
		CodigoPAF:                      "",          // valor vacío
		FechaInicioContrato:            "",          // valor vacío
		FechaFinContrato:               "",          // valor vacío
		CodigoAsignatura:               "",          // valor vacío
		NombreAsignatura:               "",          // valor vacío
		CantidadHoras:                  0,           // valor predeterminado (0)
		Jerarquia:                      "",          // valor vacío
		Calidad:                        "",          // valor vacío
		EstadoProceso:                  0,           // valor predeterminado (0)
		FechaUltimaModificacionProceso: time.Time{}, // valor vacío (cero)
	}

	return pipelsoft
}
