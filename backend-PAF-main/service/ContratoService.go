package service

import (
	"errors"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type ContratoService struct {
	DB *gorm.DB
}

// NewContratoService crea un nuevo servicio de contrato.
func NewContratoService(db *gorm.DB) *ContratoService {
	return &ContratoService{
		DB: db,
	}
}

// GetAllContratos devuelve todos los contratos registrados en la base de datos.
func (s *ContratoService) GetAllContratos() ([]models.Contrato, error) {
	var contratos []models.Contrato
	result := s.DB.Find(&contratos)
	if result.Error != nil {
		return nil, result.Error
	}
	return contratos, nil
}

// GetContratoByRun devuelve un contrato específico por el RUN del docente.
func (s *ContratoService) GetContratoByRun(run string) (*models.Contrato, error) {
	var contrato models.Contrato
	result := s.DB.Where("run_docente = ?", run).First(&contrato)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &contrato, nil
}

// GetContratosByUnidadMayor devuelve los contratos asociados a una unidad mayor específica.
func (s *ContratoService) GetContratosByUnidadMayor(unidad string) ([]models.Contrato, error) {
	var contratos []models.Contrato
	result := s.DB.Where("unidad_mayor = ?", unidad).Find(&contratos)
	if result.Error != nil {
		return nil, result.Error
	}
	return contratos, nil
}
