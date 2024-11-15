package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type HorarioService struct {
	DB *gorm.DB
}

func NewHorarioService(db *gorm.DB) *HorarioService {
	return &HorarioService{DB: db}
}

// ObtenerHorariosPorRun obtiene todos los horarios asociados a un Run
func (s *HorarioService) ObtenerHorariosPorRun(run string) ([]models.Horario, error) {
	var horarios []models.Horario
	if err := s.DB.Where("run = ?", run).Find(&horarios).Error; err != nil {
		return nil, err
	}
	return horarios, nil
}
