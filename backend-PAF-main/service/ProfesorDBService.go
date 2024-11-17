package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Asegúrate de que la ruta del modelo sea correcta
	"gorm.io/gorm"
)

type ProfesorDBService struct {
	DBProfesor *gorm.DB // Conexión a la base de datos ProfesorDB
}

func NewProfesorDBService(dbProfesor *gorm.DB) *ProfesorDBService {
	return &ProfesorDBService{DBProfesor: dbProfesor}
}

// ObtenerProfesoresPorCodigoAsignatura obtiene todos los profesores con el mismo código de asignatura en la DB ProfesorDB
func (s *ProfesorDBService) ObtenerProfesoresPorCodigoAsignatura(codigoAsignatura string) ([]models.ProfesorDB, error) {
	var profesores []models.ProfesorDB

	// Realizar la consulta en la base de datos ProfesorDB
	if err := s.DBProfesor.Where("cod_asignatura = ?", codigoAsignatura).
		Find(&profesores).Error; err != nil {
		return nil, err
	}

	return profesores, nil
}

// ObtenerProfesorPorCorreo obtiene un profesor por su correo electrónico desde la DB ProfesorDB
//func (s *ProfesorDBService) ObtenerProfesorPorCorreo(correo string) (*models.ProfesorDB, error) {
//	var profesor models.ProfesorDB

// Buscar profesor por correo en ProfesorDB
//	if err := s.DBProfesor.Where("correo = ?", correo).
//		First(&profesor).Error; err != nil {
//		return nil, err
//	}

//	return &profesor, nil
//}

// ObtenerListaProfesores obtiene todos los registros de profesores en la DB ProfesorDB
func (s *ProfesorDBService) ObtenerListaProfesores() ([]models.ProfesorDB, error) {
	var profesores []models.ProfesorDB

	// Obtener todos los profesores desde ProfesorDB
	if err := s.DBProfesor.Find(&profesores).Error; err != nil {
		return nil, err
	}

	return profesores, nil
}

// ObtenerProfesorPorRUT obtiene un profesor por su RUT desde la DB ProfesorDB
func (s *ProfesorDBService) ObtenerProfesorPorRUT(run string) (*models.ProfesorDB, error) {
	var profesor models.ProfesorDB

	// Buscar profesor por RUT en ProfesorDB
	if err := s.DBProfesor.Where("run = ?", run).
		First(&profesor).Error; err != nil {
		return nil, err
	}

	return &profesor, nil
}