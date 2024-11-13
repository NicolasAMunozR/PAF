package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type ContratoService struct {
	DB *gorm.DB
}

// NewContratoService crea una nueva instancia de ContratoService
func NewContratoService(db *gorm.DB) *ContratoService {
	return &ContratoService{
		DB: db,
	}
}

// ObtenerContratoPorID obtiene un contrato por su ID
func (s *ContratoService) ObtenerContratoPorID(id uint) (*models.Contrato, error) {
	var contrato models.Contrato
	if err := s.DB.First(&contrato, id).Error; err != nil {
		return nil, err
	}
	return &contrato, nil
}

// ObtenerTodosContratos obtiene todos los contratos
func (s *ContratoService) ObtenerTodosContratos() ([]models.Contrato, error) {
	var contratos []models.Contrato
	if err := s.DB.Find(&contratos).Error; err != nil {
		return nil, err
	}
	return contratos, nil
}

// ObtenerContratosUltimos7Dias obtiene todos los contratos creados en los últimos 7 días
func (s *ContratoService) ObtenerContratosUltimos7Dias() ([]models.Contrato, error) {
	var contratos []models.Contrato
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	if err := s.DB.Where("created_at >= ?", sevenDaysAgo).Find(&contratos).Error; err != nil {
		return nil, err
	}
	return contratos, nil
}

// ObtenerContratosPorCodigoAsignatura obtiene todos los contratos con el mismo código de asignatura
func (s *ContratoService) ObtenerContratosPorCodigoAsignatura(codigoAsignatura string) ([]models.Contrato, error) {
	var contratos []models.Contrato
	if err := s.DB.Where("codigo_asignatura = ?", codigoAsignatura).Find(&contratos).Error; err != nil {
		return nil, err
	}
	return contratos, nil
}

// ObtenerContratosUltimoMes obtiene todos los contratos creados en el último mes
func (s *ContratoService) ObtenerContratosUltimoMes() ([]models.Contrato, error) {
	// Obtenemos la fecha actual
	currentDate := time.Now()

	// Calculamos la fecha de hace un mes
	oneMonthAgo := currentDate.AddDate(0, -1, 0)

	var contratos []models.Contrato
	// Filtramos los contratos cuyo fecha_inicio esté dentro del último mes
	err := s.DB.Where("fecha_inicio >= ?", oneMonthAgo).Find(&contratos).Error
	if err != nil {
		return nil, err
	}
	return contratos, nil
}
