package service

import (
	"fmt"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Cambiar con el paquete donde se encuentran los modelos
	"gorm.io/gorm"
)

type EstadisticasResponse struct {
	TotalProfesores            int64 // Cambiar de int a int64
	TotalPipelsoft             int64 // Cambiar de int a int64
	TotalPipelsoftUnicos       int64 // Cambiar de int a int64
	PorcentajeUnicos           float64
	ProfesoresNoEnPipelsoft    int64 // Cambiar de int a int64
	ProfesoresNoEnPipelsoftPct float64
	EstadoProcesoCount         map[string]int     // Claves de tipo string
	EstadoProcesoPct           map[string]float64 // Claves de tipo string
}

// EstadisticasService define los métodos para obtener estadísticas de las tablas
type EstadisticasService struct {
	DB *gorm.DB
}

// Constructor del servicio
func NewEstadisticasService(dbPersonal *gorm.DB) *EstadisticasService {
	return &EstadisticasService{
		DB: dbPersonal,
	}
}

// ObtenerEstadisticas obtiene las estadísticas generales de los profesores y Pipelsofts
func (s *EstadisticasService) ObtenerEstadisticas() (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Contar todos los registros en la tabla profesor_dbs
	if err := s.DB.Model(&models.Contrato{}).Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores: %w", err)
	}

	// Contar todos los registros en la tabla pipelsofts
	if err := s.DB.Model(&models.Pipelsoft{}).Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts
	if err := s.DB.Model(&models.Pipelsoft{}).Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts: %w", err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los Run de los profesores que no existen en pipelsofts
	var profesoresNoEnPipelsoft int64
	if err := s.DB.Table("contratos").
		Where("run_docente NOT IN (SELECT run_empleado FROM pipelsofts)").
		Count(&profesoresNoEnPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores que no están en pipelsofts: %w", err)
	}

	// Asignar el valor a la estructura
	resp.ProfesoresNoEnPipelsoft = profesoresNoEnPipelsoft

	// Calcular porcentaje de profesores no en pipelsoft
	if resp.TotalProfesores > 0 {
		resp.ProfesoresNoEnPipelsoftPct = float64(profesoresNoEnPipelsoft) / float64(resp.TotalProfesores) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso (código de estado como string)
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"A1", "A2", "A3", "B1", "B9", "C1D", "C9D", "F1", "F9", "A9",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		if err := s.DB.Model(&models.Pipelsoft{}).Where("cod_estado = ?", estado).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s: %w", estado, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	return &resp, nil
}

// ContarRegistrosPorNombreUnidadMayor cuenta los registros en Pipelsoft que coinciden con el nombre de la unidad Mayor
func (s *EstadisticasService) ContarRegistrosPorNombreUnidadMayor(nombreUnidadMayor string) (int64, error) {
	var count int64

	// Contar los registros que coinciden con el nombre de unidad Mayor
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_Mayor = ?", nombreUnidadMayor).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error al contar los registros por nombre de unidad Mayor: %w", err)
	}

	return count, nil
}

// ObtenerFrecuenciaNombreUnidadMayor devuelve un mapa con cada NombreUnidadMayor y la cantidad de veces que aparece
func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMayor() (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		Conteo            int
	}

	// Realizar la consulta agrupada
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(*) as conteo").
		Group("nombre_unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener la frecuencia de NombreUnidadMayor: %w", err)
	}

	// Convertir los resultados en un mapa
	frecuencia := make(map[string]int)
	for _, resultado := range resultados {
		frecuencia[resultado.NombreUnidadMayor] = resultado.Conteo
	}

	return frecuencia, nil
}

//mostrar paf solo las activas y contrastarlas con la cantidad de profesores totales.

// ContarRegistrosExcluyendoEstados cuenta los registros en Pipelsoft donde cod_estado no sea "F1", "F9" o "A9"
// y calcula el porcentaje dividiéndolo por el total de profesores.
func (s *EstadisticasService) ContarRegistrosExcluyendoEstados() (int64, float64, error) {
	var count int64
	var totalProfesores int64

	// Contar los registros en Pipelsoft donde cod_estado no sea "F1", "F9", o "A9"
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}).
		Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar los registros excluyendo estados: %w", err)
	}

	// Contar el total de profesores
	if err := s.DB.Model(&models.ProfesorDB{}).
		Count(&totalProfesores).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar el total de profesores: %w", err)
	}

	// Calcular el porcentaje
	var porcentaje float64
	if totalProfesores > 0 {
		porcentaje = float64(count) / float64(totalProfesores) * 100
	}

	return count, porcentaje, nil
}

func (s *EstadisticasService) ObtenerEstadisticasPorUnidadMayor(unidadMayor string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Validar que el parámetro no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidad-mayor' es obligatorio")
	}

	// Contar los registros en pipelsofts que coincidan con la unidad mayor
	if err := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ?", unidadMayor).Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts para la unidad mayor %s: %w", unidadMayor, err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts para la unidad mayor
	if err := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ?", unidadMayor).
		Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts para la unidad mayor %s: %w", unidadMayor, err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso filtrando por unidad mayor
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"A1", "A2", "A3", "B1", "B9", "C1D", "C9D", "F1", "F9", "A9",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		if err := s.DB.Model(&models.Pipelsoft{}).
			Where("nombre_unidad_mayor = ? AND cod_estado = ?", unidadMayor, estado).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s para la unidad mayor %s: %w", estado, unidadMayor, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	return &resp, nil
}

func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayor(nombreUnidadMayor string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		Conteo            int
	}

	// Validar que se proporciona el parámetro
	if nombreUnidadMayor == "" {
		return nil, fmt.Errorf("el nombre de la unidad mayor es obligatorio")
	}

	// Realizar la consulta agrupada
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(*) as conteo").
		Where("nombre_unidad_mayor = ?", nombreUnidadMayor).
		Group("nombre_unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener la frecuencia de NombreUnidadMenor: %w", err)
	}

	// Convertir los resultados en un mapa
	frecuencia := make(map[string]int)
	for _, resultado := range resultados {
		frecuencia[resultado.NombreUnidadMenor] = resultado.Conteo
	}

	return frecuencia, nil
}
