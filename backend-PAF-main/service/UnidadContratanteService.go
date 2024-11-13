package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type UnidadContratanteService struct {
	DB *gorm.DB
}

// NewUnidadContratanteService crea una nueva instancia de UnidadContratanteService
func NewUnidadContratanteService(db *gorm.DB) *UnidadContratanteService {
	return &UnidadContratanteService{DB: db}
}

// ObtenerUnidadPorCodigo obtiene una unidad contratante por su código
func (s *UnidadContratanteService) ObtenerUnidadPorCodigo(codigo string) (*models.UnidadContratante, error) {
	var unidad models.UnidadContratante
	err := s.DB.Where("codigo = ?", codigo).First(&unidad).Error
	return &unidad, err
}

// ObtenerUnidadPorID obtiene una unidad contratante por su ID
func (s *UnidadContratanteService) ObtenerUnidadPorID(id uint) (*models.UnidadContratante, error) {
	var unidad models.UnidadContratante
	err := s.DB.First(&unidad, id).Error
	return &unidad, err
}

// ObtenerTodasUnidades obtiene todas las unidades contratantes
func (s *UnidadContratanteService) ObtenerTodasUnidades() ([]models.UnidadContratante, error) {
	var unidades []models.UnidadContratante
	err := s.DB.Find(&unidades).Error
	return unidades, err
}
