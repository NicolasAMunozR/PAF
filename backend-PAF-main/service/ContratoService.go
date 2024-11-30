package service

import (
	"errors"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/gorm"
)

type ContratoService struct {
	DB *gorm.DB
}

// NewContratoService crea un nuevo servicio de contrato.
func NewContratoService(db *gorm.DB) *ContratoService {
	return &ContratoService{
		DB: db,
	}
}

// GetAllContratos devuelve todos los contratos registrados en la base de datos.
func (s *ContratoService) GetAllContratos() ([]models.Contrato, error) {
	var contratos []models.Contrato
	result := s.DB.Find(&contratos)
	if result.Error != nil {
		return nil, result.Error
	}
	return contratos, nil
}

// GetContratoByRun devuelve un contrato específico por el RUN del docente.
func (s *ContratoService) GetContratoByRun(run string) (*models.Contrato, error) {
	var contrato models.Contrato
	result := s.DB.Where("run_docente = ?", run).First(&contrato)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &contrato, nil
}

func (s *ContratoService) GetContratosByUnidadMayor(unidad string) ([]models.Contrato, int, error) {
	var contratos []models.Contrato
	result := s.DB.Where("unidad_mayor = ?", unidad).Find(&contratos)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Contar el número de elementos de cada unidad menor
	numElementos := len(contratos)

	return contratos, numElementos, nil
}

// CountContratosByUnidadMayor cuenta la cantidad de contratos por cada unidad mayor y los profesores en Pipelsoft.
func (s *ContratoService) CountContratosByUnidadMayor() (map[string]int, map[string]int, error) {
	var contratos []models.Contrato
	result := s.DB.Find(&contratos)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var pipelsofts []models.Pipelsoft
	result = s.DB.Find(&pipelsofts)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	// Crear un mapa para contar los contratos por unidad mayor
	contratoCounts := make(map[string]int)
	for _, contrato := range contratos {
		contratoCounts[contrato.UnidadMayor]++
	}

	// Crear un mapa para contar los profesores de pipelsofts por unidad mayor
	pipelsoftCounts := make(map[string]int)
	for _, contrato := range contratos {
		for _, pipelsoft := range pipelsofts {
			if contrato.RunDocente == pipelsoft.RunEmpleado {
				pipelsoftCounts[contrato.UnidadMayor]++
				break
			}
		}
	}

	return contratoCounts, pipelsoftCounts, nil
}
