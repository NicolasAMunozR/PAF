package service

import (
	"errors"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type UsuariosService struct {
	DB *gorm.DB
}

// Nuevo servicio para usuarios
func NewUsuariosService(db *gorm.DB) *UsuariosService {
	return &UsuariosService{DB: db}
}

// Método para obtener una lista de usuarios por su Run
func (s *UsuariosService) GetUsuariosByRun(run string) ([]models.Usuarios, error) {
	if run == "" {
		return nil, errors.New("el run no puede estar vacío")
	}

	var usuarios []models.Usuarios
	result := s.DB.Where("run = ?", run).Find(&usuarios)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(usuarios) == 0 {
		return nil, errors.New("no se encontraron usuarios con el RUN proporcionado")
	}

	return usuarios, nil
}
