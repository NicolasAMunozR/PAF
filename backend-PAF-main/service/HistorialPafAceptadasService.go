package service

import (
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Asegúrate de que los modelos estén correctos
	"gorm.io/gorm"
)

type HistorialPafAceptadasService struct {
	DBPersonal  *gorm.DB // Conexión a DBPersonal
	DBPipelsoft *gorm.DB // Conexión a DBPipelsoft
}

func NewHistorialPafAceptadasService(dbPersonal, dbPipelsoft *gorm.DB) *HistorialPafAceptadasService {
	return &HistorialPafAceptadasService{
		DBPersonal:  dbPersonal,
		DBPipelsoft: dbPipelsoft,
	}
}

// CrearHistorial crea un nuevo registro en HistorialPafAceptadas en la DBPersonal
func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF string) (*models.HistorialPafAceptadas, error) {
	// Primero, obtener los datos desde la base de datos de Pipelsoft (DBPipelsoft)
	var pipelsoft models.Pipelsoft
	if err := s.DBPipelsoft.Where("codigo_paf = ?", codigoPAF).First(&pipelsoft).Error; err != nil {
		return nil, err
	}

	// Crear el nuevo historial en DBPersonal
	historial := models.HistorialPafAceptadas{
		Run:                pipelsoft.Run,
		CodigoPAF:          codigoPAF,
		FechaAceptacionPaf: time.Now(),
	}

	// Insertar el historial en DBPersonal
	if err := s.DBPersonal.Create(&historial).Error; err != nil {
		return nil, err
	}

	return &historial, nil
}

// ObtenerHistorialPorID devuelve un registro de HistorialPafAceptadas por su ID desde DBPersonal
func (s *HistorialPafAceptadasService) ObtenerHistorialPorID(id uint) (*models.HistorialPafAceptadas, error) {
	var historial models.HistorialPafAceptadas
	// Buscar el historial en DBPersonal por su ID
	if err := s.DBPersonal.First(&historial, id).Error; err != nil {
		return nil, err
	}
	return &historial, nil
}

// EliminarHistorial elimina un registro de HistorialPafAceptadas por su CodigoPAF en DBPersonal
func (s *HistorialPafAceptadasService) EliminarHistorial(codigoPAF string) error {
	var historial models.HistorialPafAceptadas
	// Buscar el registro en DBPersonal por CodigoPAF
	if err := s.DBPersonal.Where("codigo_paf = ?", codigoPAF).First(&historial).Error; err != nil {
		return err
	}
	// Eliminar el registro en DBPersonal
	return s.DBPersonal.Delete(&historial).Error
}

// ObtenerTodosLosHistoriales devuelve todos los registros de HistorialPafAceptadas desde DBPersonal
func (s *HistorialPafAceptadasService) ObtenerTodosLosHistoriales() ([]models.HistorialPafAceptadas, error) {
	var historiales []models.HistorialPafAceptadas
	// Obtener todos los historiales desde DBPersonal
	if err := s.DBPersonal.Find(&historiales).Error; err != nil {
		return nil, err
	}
	return historiales, nil
}
