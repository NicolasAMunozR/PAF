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

// Método para obtener un usuario por su Run
func (s *UsuariosService) GetUsuarioByRun(run string) (*models.Usuarios, error) {
	if run == "" {
		return nil, errors.New("el run no puede estar vacío")
	}

	var usuario models.Usuarios
	result := s.DB.Where("run = ?", run).First(&usuario)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, result.Error
	}

	return &usuario, nil
}
