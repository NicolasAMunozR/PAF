package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type HistorialPafAceptadasService struct {
	DB *gorm.DB // Conexión única a DBPersonal
}

func NewHistorialPafAceptadasService(db *gorm.DB) *HistorialPafAceptadasService {
	return &HistorialPafAceptadasService{
		DB: db,
	}
}

// CrearHistorial crea un nuevo registro en HistorialPafAceptadas en la DBPersonal
func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF string) (*models.HistorialPafAceptadas, error) {
	// Obtener los datos desde la tabla correspondiente en DBPersonal
	var pipelsoft models.Pipelsoft
	if err := s.DB.Where("codigo_paf = ?", codigoPAF).First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	var estado = pipelsoft.EstadoProceso
	estado = estado + 1

	// Crear el nuevo historial en DBPersonal
	historial := models.HistorialPafAceptadas{
		CodigoPAF:           pipelsoft.CodigoPAF,
		FechaInicioContrato: pipelsoft.FechaInicioContrato,
		FechaFinContrato:    pipelsoft.FechaFinContrato,
		CodigoAsignatura:    pipelsoft.CodigoAsignatura,
		NombreAsignatura:    pipelsoft.NombreAsignatura,
		CantidadHoras:       pipelsoft.CantidadHoras,
		Jerarquia:           pipelsoft.Jerarquia,
		Calidad:             pipelsoft.Calidad,
		EstadoProceso:       estado, // la idea es que cuando ellas aceptan la paf el estado ya deberia aumentar

	}

	// Insertar el historial en DBPersonal
	if err := s.DB.Create(&historial).Error; err != nil {
		return nil, err
	}

	return &historial, nil
}

// ObtenerHistorialPorID devuelve un registro de HistorialPafAceptadas por su ID desde DBPersonal
func (s *HistorialPafAceptadasService) ObtenerHistorialPorID(id uint) (*models.HistorialPafAceptadas, error) {
	var historial models.HistorialPafAceptadas
	// Buscar el historial en DBPersonal por su ID
	if err := s.DB.First(&historial, id).Error; err != nil {
		return nil, err
	}
	return &historial, nil
}

// EliminarHistorial elimina un registro de HistorialPafAceptadas por su CodigoPAF en DBPersonal
func (s *HistorialPafAceptadasService) EliminarHistorial(codigoPAF string) error {
	// Eliminar directamente usando la condición
	return s.DB.Where("codigo_paf = ?", codigoPAF).Delete(&models.HistorialPafAceptadas{}).Error
}

// ObtenerTodosLosHistoriales devuelve todos los registros de HistorialPafAceptadas desde DBPersonal
func (s *HistorialPafAceptadasService) ObtenerTodosLosHistoriales() ([]models.HistorialPafAceptadas, error) {
	var historiales []models.HistorialPafAceptadas
	// Obtener todos los historiales desde DBPersonal
	if err := s.DB.Find(&historiales).Error; err != nil {
		return nil, err
	}
	return historiales, nil
}
