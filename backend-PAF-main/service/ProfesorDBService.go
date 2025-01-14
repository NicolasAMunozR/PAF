package service

import (
	"fmt"

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

// no deben estar desde la tabla de contratos
// deben ser datos que no estan en contratos pero si en el sai
// se debe contar y enviar la info cuando la data no esta en contrato
// o cuando el parametro planta es distinto de ACADEMICO con tilde en la E
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

// ObtenerProfesoresSinContratoYNoAcademico obtiene:
// 1. Profesores en `ProfesorDB` que no tienen contrato.
// 2. Profesores en `ProfesorDB` asociados a contratos donde la `Planta` es distinta de "ACADÉMICO".
func (s *ProfesorDBService) ObtenerProfesoresSinContratoYNoAcademico() (map[string]interface{}, error) {
	// 1. Profesores que no tienen contrato
	var profesoresSinContrato []models.ProfesorDB
	err := s.DBPersonal.Raw(`
		SELECT * 
		FROM profesor_dbs 
		WHERE run NOT IN (SELECT run_docente FROM contratos)
	`).Scan(&profesoresSinContrato).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener profesores sin contrato: %w", err)
	}

	// 2. Profesores asociados a contratos donde `Planta` es distinta de "ACADÉMICO"
	var profesoresNoAcademico []models.ProfesorDB
	err = s.DBPersonal.Raw(`
		SELECT p.*
		FROM profesor_dbs p
		JOIN contratos c ON p.run = c.run_docente
		WHERE c.planta != 'ACADÉMICO'
	`).Scan(&profesoresNoAcademico).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener profesores no académicos: %w", err)
	}

	// Consolidar resultados en un mapa
	resultado := map[string]interface{}{
		"profesores_sin_contrato": profesoresSinContrato,
		"profesores_no_academico": profesoresNoAcademico,
	}

	return resultado, nil
}
