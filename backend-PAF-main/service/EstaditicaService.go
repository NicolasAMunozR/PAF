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
	if err := s.DB.Model(&models.ProfesorDB{}).Count(&resp.TotalProfesores).Error; err != nil {
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
	if err := s.DB.Table("profesor_dbs").
		Where("run NOT IN (SELECT run_empleado FROM pipelsofts)").
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

// ContarRegistrosPorNombreUnidadContratante cuenta los registros en Pipelsoft que coinciden con el nombre de la unidad contratante
func (s *EstadisticasService) ContarRegistrosPorNombreUnidadContratante(nombreUnidadContratante string) (int64, error) {
	var count int64

	// Contar los registros que coinciden con el nombre de unidad contratante
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_contratante = ?", nombreUnidadContratante).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error al contar los registros por nombre de unidad contratante: %w", err)
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
