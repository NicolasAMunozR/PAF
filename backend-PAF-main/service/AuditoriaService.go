package service

import (
	"fmt"
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type AuditoriaService struct {
	DB *gorm.DB
}

// NewAuditoriaService crea una nueva instancia del servicio de auditoría.
func NewAuditoriaService(db *gorm.DB) *AuditoriaService {
	return &AuditoriaService{DB: db}
}

// RegistrarAuditoria agrega un nuevo registro de auditoría.
func (s *AuditoriaService) RegistrarAuditoria(rut, tipoMod, descripcion string) error {
	auditoria := models.Auditoria{
		Rut:                    rut,
		TipoDeModificacion:     tipoMod,
		DescipcionModificacion: descripcion,
		FechaModificacion:      time.Now(),
	}

	if err := s.DB.Create(&auditoria).Error; err != nil {
		return fmt.Errorf("error al registrar auditoría: %w", err)
	}

	return nil
}
