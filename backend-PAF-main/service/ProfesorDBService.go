package service

import (
	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Asegúrate de que la ruta del modelo sea correcta
	"gorm.io/gorm"
)

type ProfesorDBService struct {
	DBPersonal *gorm.DB // Conexión a la base de datos DBPersonal
}

func NewProfesorDBService(dbPersonal *gorm.DB) *ProfesorDBService {
	return &ProfesorDBService{DBPersonal: dbPersonal}
}

// ObtenerProfesoresPorCodigoAsignatura obtiene todos los profesores con el mismo código de asignatura en la DB DBPersonal
func (s *ProfesorDBService) ObtenerProfesoresPorCodigoAsignatura(codigoAsignatura string) ([]models.ProfesorDB, error) {
	var profesores []models.ProfesorDB

	// Realizar la consulta en la base de datos DBPersonal
	if err := s.DBPersonal.Where("cod_asignatura = ?", codigoAsignatura).
		Find(&profesores).Error; err != nil {
		return nil, err
	}

	return profesores, nil
}

// ObtenerListaProfesores obtiene todos los registros de profesores en la DB DBPersonal
func (s *ProfesorDBService) ObtenerListaProfesores() ([]models.ProfesorDB, error) {
	var profesores []models.ProfesorDB

	// Obtener todos los profesores desde DBPersonal
	if err := s.DBPersonal.Find(&profesores).Error; err != nil {
		return nil, err
	}

	return profesores, nil
}

// ObtenerProfesorPorRUT obtiene un profesor por su RUT desde la DB DBPersonal
func (s *ProfesorDBService) ObtenerProfesorPorRUT(run string) ([]models.ProfesorDB, error) {
	var profesor []models.ProfesorDB

	// Buscar profesor por RUT en DBPersonal
	if err := s.DBPersonal.Where("run = ?", run).
		Find(&profesor).Error; err != nil {
		return nil, err
	}

	return profesor, nil
}

func (s *ProfesorDBService) GetCountProfesoresNotInPipelsoft() (int, error) {
	var profesores []models.ProfesorDB
	var pipelsofts []models.Pipelsoft

	// Obtener todos los RUN de ProfesorDB y Pipelsoft
	err := s.DBPersonal.Find(&profesores).Error
	if err != nil {
		return 0, err
	}

	err = s.DBPersonal.Find(&pipelsofts).Error
	if err != nil {
		return 0, err
	}

	// Crear un mapa con los RUN de Pipelsoft para búsqueda rápida
	rutPipelsoft := make(map[string]bool)
	for _, pipel := range pipelsofts {
		rutPipelsoft[pipel.RunEmpleado] = true
	}

	// Contar los profesores cuyo RUN no está en Pipelsoft
	count := 0
	for _, profesor := range profesores {
		if _, exists := rutPipelsoft[profesor.RUN]; !exists {
			count++
		}
	}

	return count, nil
}
