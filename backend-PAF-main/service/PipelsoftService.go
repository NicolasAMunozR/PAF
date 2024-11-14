// service/pipelsoft_service.go
package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type PipelsoftService struct {
	DB *gorm.DB
}

func NewPipelsoftService(db *gorm.DB) *PipelsoftService {
	return &PipelsoftService{DB: db}
}

// ObtenerContratosUltimos7Dias obtiene todos los registros creados en los últimos 7 días
func (s *PipelsoftService) ObtenerContratosUltimos7Dias() ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if err := s.DB.Where("created_at >= ?", sevenDaysAgo).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}
	return pipelsofts, nil
}

// ObtenerContratosPorCodigoAsignatura obtiene todos los registros con el mismo código de asignatura
func (s *PipelsoftService) ObtenerContratosPorCodigoAsignatura(codigoAsignatura string) ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft
	if err := s.DB.Where("codigo_asignatura = ?", codigoAsignatura).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}
	return pipelsofts, nil
}

// ObtenerContratosUltimoMes obtiene todos los registros creados en el último mes
func (s *PipelsoftService) ObtenerContratosUltimoMes() ([]models.Pipelsoft, error) {
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	var pipelsofts []models.Pipelsoft
	if err := s.DB.Where("fecha_inicio_contrato >= ?", oneMonthAgo).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}
	return pipelsofts, nil
}

// ObtenerPersonaPorCorreo devuelve un registro de `Pipelsoft` por su correo
func (s *PipelsoftService) ObtenerPersonaPorCorreo(correo string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft
	if err := s.DB.Where("correo = ?", correo).First(&pipelsoft).Error; err != nil {
		return nil, err
	}
	return &pipelsoft, nil
}

// ObtenerPersonaPorRUT devuelve un registro de `Pipelsoft` por su RUT
func (s *PipelsoftService) ObtenerPersonaPorRUT(run string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft
	if err := s.DB.Where("run = ?", run).First(&pipelsoft).Error; err != nil {
		return nil, err
	}
	return &pipelsoft, nil
}

// ObtenerProcesoPorEstado obtiene todos los registros que estén en un estado específico
func (s *PipelsoftService) ObtenerProcesoPorEstado(estado string) ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft
	if err := s.DB.Where("estado_proceso = ?", estado).Find(&pipelsofts).Error; err != nil {
		return nil, err
	}
	return pipelsofts, nil
}

// ObtenerUnidadPorCodigo obtiene un registro por su código de unidad contratante
func (s *PipelsoftService) ObtenerUnidadPorCodigo(codigo string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft
	if err := s.DB.Where("codigo_unidad_contratante = ?", codigo).First(&pipelsoft).Error; err != nil {
		return nil, err
	}
	return &pipelsoft, nil
}
