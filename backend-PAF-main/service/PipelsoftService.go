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

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("created_at >= ?", sevenDaysAgo).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

// ObtenerContratosPorCodigoAsignatura obtiene todos los registros con el mismo código de asignatura
func (s *PipelsoftService) ObtenerContratosPorCodigoAsignatura(codigoAsignatura string) ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("codigo_asignatura = ?", codigoAsignatura).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

// ObtenerContratosUltimoMes obtiene todos los registros creados en el último mes
func (s *PipelsoftService) ObtenerContratosUltimoMes() ([]models.Pipelsoft, error) {
	oneMonthAgo := time.Now().AddDate(0, -1, 0)
	var pipelsofts []models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("fecha_inicio_contrato >= ?", oneMonthAgo).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

// ObtenerPersonaPorCorreo devuelve un registro de `Pipelsoft` por su correo
func (s *PipelsoftService) ObtenerPersonaPorCorreo(correo string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("correo = ?", correo).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	return &pipelsoft, nil
}

// ObtenerUnidadPorCodigo obtiene un registro por su código de unidad contratante
func (s *PipelsoftService) ObtenerUnidadPorCodigo(codigo string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("codigo_unidad_contratante = ?", codigo).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	return &pipelsoft, nil
}

func (s *PipelsoftService) ObtenerListaPersonas() ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

// ObtenerPersonaPorRUT devuelve un registro de `Pipelsoft` por su RUT
func (s *PipelsoftService) ObtenerPersonaPorRUT(run string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("run = ?", run).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	return &pipelsoft, nil
}

// ObtenerProcesoPorEstado obtiene todos los registros que estén en un estado específico
func (s *PipelsoftService) ObtenerProcesoPorEstado(estado string) ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("estado_proceso = ?", estado).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

func (s *PipelsoftService) ObtenerPersonasPorRUT(run string) ([]models.Pipelsoft, error) {
	var pipelsofts []models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("run = ?", run).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		Find(&pipelsofts).Error; err != nil {
		return nil, err
	}

	return pipelsofts, nil
}

// ObtenerPersonaPorPaf obtiene un registro de `Pipelsoft` por su código PAF
func (s *PipelsoftService) ObtenerPersonaPorPaf(codigoPaf string) (*models.Pipelsoft, error) {
	var pipelsoft models.Pipelsoft

	// Excluir registros cuyo CodigoPAF exista en la tabla HistorialPafAceptadas
	if err := s.DB.Where("codigo_paf = ?", codigoPaf).
		Not("codigo_paf IN (?)", s.DB.Table("historial_paf_aceptadas").Select("codigo_paf")).
		First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	return &pipelsoft, nil
}