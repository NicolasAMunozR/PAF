package service

import (
	"fmt"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Cambiar con el paquete donde se encuentran los modelos
	"gorm.io/gorm"
)

// Estructura para almacenar la respuesta
type EstadisticasResponse struct {
	TotalProfesores          int64
	TotalPipelsoft           int64
	ProfesoresNoEnPipelsoft  int
	EstadoProcesoCount       map[int]int // Mapa que almacena la cantidad de registros por EstadoProceso
	TotalPipelsoftUnicos     int64       // Nuevo campo para contar los Run únicos
	RegistrosPorNombreUnidad int64       // Nuevo campo para contar los registros por nombre de unidad contratante
}

// EstadisticasService define los métodos para obtener estadísticas de las tablas
type EstadisticasService struct {
	DB *gorm.DB
}

// Nueva instancia de EstadisticasService
func NewEstadisticasService(db *gorm.DB) *EstadisticasService {
	return &EstadisticasService{DB: db}
}

// ObtenerEstadisticas obtiene todas las estadísticas solicitadas
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
	if err := s.DB.Model(&models.Pipelsoft{}).Distinct("run").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts: %w", err)
	}

	// Contar los Run de los profesores que no existen en pipelsofts
	var profesoresNoEnPipelsoft int64
	if err := s.DB.Table("profesor_dbs").Where("run NOT IN (SELECT run FROM pipelsofts)").Count(&profesoresNoEnPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores que no están en pipelsofts: %w", err)
	}
	resp.ProfesoresNoEnPipelsoft = int(profesoresNoEnPipelsoft)

	// Contar los registros en pipelsofts por cada EstadoProceso (1-6)
	resp.EstadoProcesoCount = make(map[int]int)
	for i := 1; i <= 6; i++ {
		var count int64
		if err := s.DB.Model(&models.Pipelsoft{}).Where("estado_proceso = ?", i).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %d: %w", i, err)
		}
		resp.EstadoProcesoCount[i] = int(count)
	}

	return &resp, nil
}

// ContarRegistrosPorNombreUnidadContratante cuenta los registros en Pipelsoft que coinciden con el nombre de la unidad contratante
func (s *EstadisticasService) ContarRegistrosPorNombreUnidadContratante(nombreUnidadContratante string) (int64, error) {
	var count int64

	// Contar los registros que coinciden con el nombre de unidad contratante
	if err := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_contratante = ?", nombreUnidadContratante).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error al contar los registros por nombre de unidad contratante: %w", err)
	}

	return count, nil
}
