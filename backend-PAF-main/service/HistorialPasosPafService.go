package service

import (
	"fmt"
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type HistorialPasosPafService struct {
	DB *gorm.DB
}

func NewHistorialPasosPafService(db *gorm.DB) *HistorialPasosPafService {
	return &HistorialPasosPafService{DB: db}
}

func (s *HistorialPasosPafService) ObtenerYCalcularPorIdYRun(idPaf, runDocente string) ([]models.HistorialPasosPaf, []time.Duration, error) {
	var historiales []models.HistorialPasosPaf

	// Filtrar por id_paf y run_docente, ordenado por FechaLlegadaPaf
	result := s.DB.Where("id_paf = ? AND run_docente = ?", idPaf, runDocente).
		Order("fecha_llegada_paf ASC").
		Find(&historiales)

	if result.Error != nil {
		return nil, nil, fmt.Errorf("error al obtener los datos: %w", result.Error)
	}

	// Calcular el tiempo entre FechaLlegadaPaf y FechaModificacionPaf
	var duraciones []time.Duration
	for _, historial := range historiales {
		duracion := historial.FechaModificacionPaf.Sub(historial.FechaLlegadaPaf)
		duraciones = append(duraciones, duracion)
	}

	return historiales, duraciones, nil
}
