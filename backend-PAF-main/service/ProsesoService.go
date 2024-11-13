package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"

	"gorm.io/gorm"
)

// ProcesoService estructura el servicio de procesos
type ProcesoService struct {
	DB *gorm.DB
}

// NewProcesoService crea una nueva instancia de ProcesoService
func NewProcesoService(db *gorm.DB) *ProcesoService {
	return &ProcesoService{
		DB: db,
	}
}

// ObtenerProceso por ID obtiene un proceso
func (s *ProcesoService) ObtenerProcesoPorID(id uint) (models.Proceso, error) {
	var proceso models.Proceso
	// Obtener el proceso por ID
	err := s.DB.First(&proceso, id).Error
	if err != nil {
		return models.Proceso{}, err
	}
	return proceso, nil
}

// ObtenerTodosLosProcesos obtiene todos los procesos
func (s *ProcesoService) ObtenerTodosLosProcesos() ([]models.Proceso, error) {
	var procesos []models.Proceso
	// Obtener todos los procesos
	err := s.DB.Find(&procesos).Error
	if err != nil {
		return nil, err
	}
	return procesos, nil
}

// ObtenerProcesoPorEstado obtiene todos los procesos que estén en un estado específico
func (s *ProcesoService) ObtenerProcesoPorEstado(estado string) ([]models.Proceso, error) {
	var procesos []models.Proceso
	// Filtrar procesos por estado
	err := s.DB.Where("estado = ?", estado).Find(&procesos).Error
	if err != nil {
		return nil, err
	}
	return procesos, nil
}
