package service

import (
	"time"

	"strings"

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

func (s *PipelsoftService) getAcceptedCodes() (map[int]struct{}, error) {
	var historialAceptadas []struct {
		CodigoPaf int `gorm:"column:id_paf"`
	}

	// Modificación: agregamos la condición BanderaAceptacion = 1
	if err := s.DBPersonal.Table("historial_paf_aceptadas").
		Select("id_paf").
		Where("bandera_aceptacion = ?", 1).
		Find(&historialAceptadas).Error; err != nil {
		return nil, err
	}

	aceptadosMap := make(map[int]struct{})
	for _, item := range historialAceptadas {
		aceptadosMap[item.CodigoPaf] = struct{}{}
	}
	return aceptadosMap, nil
}

func (s *PipelsoftService) filterRecords(records []models.Pipelsoft, aceptadosMap map[int]struct{}) []models.Pipelsoft {
	var filtered []models.Pipelsoft
	for _, record := range records {
		if _, found := aceptadosMap[record.IdPaf]; !found {
			filtered = append(filtered, record)
		}
	}
	return filtered
}

func (s *PipelsoftService) comprobarYCombinarDatosPorCodigoPAF(pipelsofts []models.Pipelsoft) ([]models.DatosCombinados, error) {
	// Crear un conjunto único de los códigos PAF presentes en los registros de Pipelsoft
	codigoPAFSet := make(map[int]struct{})
	for _, pipel := range pipelsofts {
		codigoPAFSet[pipel.IdPaf] = struct{}{}
	}

	// Crear una lista para almacenar los datos de HistorialPafAceptadas
	var historialPafAceptadas []models.HistorialPafAceptadas
	codigosPAF := make([]int, 0, len(codigoPAFSet))
	for codigoPAF := range codigoPAFSet {
		codigosPAF = append(codigosPAF, codigoPAF)
	}

	// Consultar los datos de HistorialPafAceptadas cuyo CódigoPAF esté en la lista de códigos
	if err := s.DBPersonal.Where("id_paf IN ?", codigosPAF).Find(&historialPafAceptadas).Error; err != nil {
		return nil, err
	}

	// Crear un mapa para almacenar los registros de HistorialPafAceptadas por CódigoPAF y CodigoAsignatura
	historialPafMap := make(map[int]map[string]models.HistorialPafAceptadas)
	for _, historial := range historialPafAceptadas {
		if _, exists := historialPafMap[historial.IdPaf]; !exists {
			historialPafMap[historial.IdPaf] = make(map[string]models.HistorialPafAceptadas)
		}
		historialPafMap[historial.IdPaf][historial.CodigoAsignatura] = historial
	}

	// Crear una lista para almacenar los datos combinados
	var datosCombinados []models.DatosCombinados
	for _, pipel := range pipelsofts {
		// Crear un nuevo objeto DatosCombinados solo con los datos de Pipelsoft
		dato := models.DatosCombinados{
			PipelsoftData: pipel,
		}

		// Buscar si el CódigoPAF y CodigoAsignatura del registro de Pipelsoft están en HistorialPafAceptadas
		if historial, found := historialPafMap[pipel.IdPaf][pipel.CodigoAsignatura]; found {
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

	if err := s.DBPersonal.Where("run_empleado = ?", run).Find(&pipelsofts).Error; err != nil {
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

	if err := s.DBPersonal.Where("id_paf = ?", codigoPAF).Find(&pipelsofts).Error; err != nil {
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

// Obtener contratos por nombre de unidad mayor
func (s *PipelsoftService) ObtenerContratosPorNombreUnidadMayor(nombreUnidadMayor string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("nombre_unidad_mayor = ?", nombreUnidadMayor).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
}

// Obtener contratos por nombre de unidad menor
func (s *PipelsoftService) ObtenerContratosPorNombreUnidadMenor(nombreUnidadMenor string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("nombre_unidad_menor = ?", nombreUnidadMenor).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	aceptadosMap, err := s.getAcceptedCodes()
	if err != nil {
		return nil, err
	}
	contratosFiltrados := s.filterRecords(pipelsofts, aceptadosMap)

	return s.comprobarYCombinarDatosPorCodigoPAF(contratosFiltrados)
}

// ObtenerNombreUnidadMenorPorMayor obtiene todas las unidades menores asociadas a un NombreUnidadMayor
func (s *PipelsoftService) ObtenerNombreUnidadMenorPorMayor(nombreUnidadMayor string) ([]string, error) {
	var unidadesMenores []string

	// Consultar la base de datos para obtener los valores únicos de NombreUnidadMenor
	err := s.DBPersonal.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_mayor = ?", nombreUnidadMayor).
		Distinct("nombre_unidad_menor").
		Pluck("nombre_unidad_menor", &unidadesMenores).
		Error

	if err != nil {
		return nil, err
	}

	return unidadesMenores, nil
}

func (s *PipelsoftService) ObtenerContratosPorRUNMostrarTodo(run string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("run_empleado = ?", run).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return s.comprobarYCombinarDatosPorCodigoPAF(pipelsofts)
}

func (s *PipelsoftService) ObtenerContratosPorIdPafMostrarTodo(run string) ([]models.DatosCombinados, error) {
	var pipelsofts []models.Pipelsoft

	if err := s.DBPersonal.Where("id_paf = ?", run).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return s.comprobarYCombinarDatosPorCodigoPAF(pipelsofts)
}

func (s *PipelsoftService) GetUniqueUnits() ([]string, []string, error) {
	var pipelsofts []models.Pipelsoft
	err := s.DBPersonal.Find(&pipelsofts).Error
	if err != nil {
		return nil, nil, err
	}

	// Mapas para eliminar duplicados
	uniqueMayores := map[string]bool{}
	uniqueMenores := map[string]bool{}

	for _, pipel := range pipelsofts {
		// Limpiar espacios adicionales y validar
		mayor := strings.TrimSpace(pipel.NombreUnidadMayor)
		menor := strings.TrimSpace(pipel.NombreUnidadMenor)

		if mayor != "" {
			uniqueMayores[mayor] = true
		}
		if menor != "" {
			uniqueMenores[menor] = true
		}
	}

	// Convertir mapas a listas
	mayores := make([]string, 0, len(uniqueMayores))
	for mayor := range uniqueMayores {
		mayores = append(mayores, mayor)
	}

	menores := make([]string, 0, len(uniqueMenores))
	for menor := range uniqueMenores {
		menores = append(menores, menor)
	}

	return mayores, menores, nil
}

func (s *PipelsoftService) GetUnitsByMayor(nombreUnidadMayor string) ([]string, error) {
	var pipelsofts []models.Pipelsoft
	err := s.DBPersonal.Where("nombre_unidad_mayor = ?", nombreUnidadMayor).Find(&pipelsofts).Error
	if err != nil {
		return nil, err
	}

	uniqueMenores := map[string]bool{}
	for _, pipel := range pipelsofts {
		menor := strings.TrimSpace(pipel.NombreUnidadMenor)
		if menor != "" {
			uniqueMenores[menor] = true
		}
	}

	menores := make([]string, 0, len(uniqueMenores))
	for menor := range uniqueMenores {
		menores = append(menores, menor)
	}

	return menores, nil
}

func (s *PipelsoftService) GetAll() ([]models.PipelsoftDTO, error) {
	var pipelsofts []models.Pipelsoft
	err := s.DBPersonal.Find(&pipelsofts).Error // Recupera todos los registros sin filtrar por semestre
	if err != nil {
		return nil, err
	}

	// Mapa para agrupar por IdPaf y recolectar los códigos de asignatura
	grouped := make(map[int]*models.PipelsoftDTO)

	for _, pipel := range pipelsofts {
		if existing, exists := grouped[pipel.IdPaf]; exists {
			// Agregar el CódigoAsignatura si es diferente
			if !contains(existing.CodigoAsignaturaList, pipel.CodigoAsignatura) {
				existing.CodigoAsignaturaList = append(existing.CodigoAsignaturaList, pipel.CodigoAsignatura)
			}
		} else {
			// Crear una nueva entrada en el mapa con el primer elemento
			grouped[pipel.IdPaf] = &models.PipelsoftDTO{
				RunEmpleado:          pipel.RunEmpleado,
				Nombres:              pipel.Nombres,
				PrimerApp:            pipel.PrimerApp,
				SegundoApp:           pipel.SegundoApp,
				NombreUnidadMayor:    pipel.NombreUnidadMayor,
				NombreUnidadMenor:    pipel.NombreUnidadMenor,
				IdPaf:                pipel.IdPaf,
				FechaInicioContrato:  pipel.FechaInicioContrato,
				FechaFinContrato:     pipel.FechaFinContrato,
				NombreAsignatura:     pipel.NombreAsignatura,
				HorasAsignatura:      pipel.HorasAsignatura,
				CantidadHorasPaf:     pipel.CantidadHorasPaf,
				Jerarquia:            pipel.Jerarquia,
				Semestre:             pipel.Semestre,
				UltimaModificacion:   pipel.UltimaModificacion,
				Categoria:            pipel.Categoria,
				CodEstado:            pipel.CodEstado,
				DesEstado:            pipel.DesEstado,
				Llave:                pipel.Llave,
				Veces:                pipel.Veces,
				CodigoAsignaturaList: []string{pipel.CodigoAsignatura}, // Lista inicial con un solo código
			}
		}
	}

	// Convertir el mapa en una lista de resultados
	result := make([]models.PipelsoftDTO, 0, len(grouped))
	for _, pipel := range grouped {
		result = append(result, *pipel)
	}

	return result, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
