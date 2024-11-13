package service

import (
	"errors"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type PersonaService struct {
	DB *gorm.DB
}

// NewPersonaService crea una instancia de PersonaService
func NewPersonaService(db *gorm.DB) *PersonaService {
	return &PersonaService{DB: db}
}

// ObtenerPersonaPorID devuelve una Persona por su ID
func (s *PersonaService) ObtenerPersonaPorID(id uint) (*models.Persona, error) {
	var persona models.Persona
	if err := s.DB.First(&persona, id).Error; err != nil {
		return nil, errors.New("Persona no encontrada")
	}
	return &persona, nil
}

// ObtenerTodasPersonas devuelve todas las Personas
func (s *PersonaService) ObtenerTodasPersonas() ([]models.Persona, error) {
	var personas []models.Persona
	if err := s.DB.Find(&personas).Error; err != nil {
		return nil, err
	}
	return personas, nil
}

// ObtenerPersonaPorCorreo devuelve una Persona por su correo
func (s *PersonaService) ObtenerPersonaPorCorreo(correo string) (*models.Persona, error) {
	var persona models.Persona
	if err := s.DB.Where("correo = ?", correo).First(&persona).Error; err != nil {
		return nil, errors.New("Persona no encontrada")
	}
	return &persona, nil
}

// ObtenerPersonaPorRUT devuelve una Persona por su RUT
func (s *PersonaService) ObtenerPersonaPorRUT(run string) (*models.Persona, error) {
	var persona models.Persona
	if err := s.DB.Where("run = ?", run).First(&persona).Error; err != nil {
		return nil, errors.New("Persona no encontrada con ese RUT")
	}
	return &persona, nil
}
