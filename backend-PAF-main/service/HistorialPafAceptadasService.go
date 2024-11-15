// service/historial_paf_aceptadas_service.go
package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type HistorialPafAceptadasService struct {
	DB *gorm.DB
}

func NewHistorialPafAceptadasService(db *gorm.DB) *HistorialPafAceptadasService {
	return &HistorialPafAceptadasService{DB: db}
}

// CrearHistorial crea un nuevo registro en HistorialPafAceptadas
func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF string) (*models.HistorialPafAceptadas, error) {
	var pipelsoft models.Pipelsoft
	if err := s.DB.Where("codigo_paf = ?", codigoPAF).First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	historial := models.HistorialPafAceptadas{
		Run:                pipelsoft.Run,
		CodigoPAF:          codigoPAF,
		FechaAceptacionPaf: time.Now(),
	}

	if err := s.DB.Create(&historial).Error; err != nil {
		return nil, err
	}

	return &historial, nil
}

// ObtenerHistorialPorID devuelve un registro de HistorialPafAceptadas por su ID
func (s *HistorialPafAceptadasService) ObtenerHistorialPorID(id uint) (*models.HistorialPafAceptadas, error) {
	var historial models.HistorialPafAceptadas
	if err := s.DB.First(&historial, id).Error; err != nil {
		return nil, err
	}
	return &historial, nil
}

// EliminarHistorial elimina un registro de HistorialPafAceptadas por su CodigoPAF
func (s *HistorialPafAceptadasService) EliminarHistorial(codigoPAF string) error {
	var historial models.HistorialPafAceptadas
	// Buscar el registro por CodigoPAF
	if err := s.DB.Where("codigo_paf = ?", codigoPAF).First(&historial).Error; err != nil {
		return err
	}
	// Eliminar el registro encontrado
	return s.DB.Delete(&historial).Error
}

// ObtenerTodosLosHistoriales devuelve todos los registros de HistorialPafAceptadas
func (s *HistorialPafAceptadasService) ObtenerTodosLosHistoriales() ([]models.HistorialPafAceptadas, error) {
	var historiales []models.HistorialPafAceptadas
	if err := s.DB.Find(&historiales).Error; err != nil {
		return nil, err
	}
	return historiales, nil
}
