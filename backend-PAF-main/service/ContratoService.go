package service

import (
	"errors"
	"fmt"
	"log"

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

func (s *ContratoService) GetContratoByRun(run string) (*models.Contrato, error) {
	var contrato models.Contrato
	result := s.DB.Where("run_docente = ?", run).First(&contrato)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Es mejor devolver un error específico en lugar de retornar nil
		return nil, fmt.Errorf("contrato no encontrado para el RUN %s", run)
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

func (s *ContratoService) ProfesorUnidadMayorYPaf() (map[string]int, map[string]int, error) {
	var contratos []models.Contrato
	result := s.DB.Find(&contratos)
	if result.Error != nil {
		// Devolver mapas vacíos si ocurre un error al obtener los contratos
		return make(map[string]int), make(map[string]int), result.Error
	}

	// Crear un mapa de Pipelsoft para mejorar la búsqueda
	var pipelsofts []models.Pipelsoft
	result = s.DB.Find(&pipelsofts)
	if result.Error != nil {
		// Devolver mapas vacíos si ocurre un error al obtener los pipelsofts
		return make(map[string]int), make(map[string]int), result.Error
	}

	// Crear un mapa de RunEmpleado de Pipelsoft para búsqueda rápida
	pipelsoftMap := make(map[string]struct{}) // Mapa de RunEmpleado
	for _, pipelsoft := range pipelsofts {
		pipelsoftMap[pipelsoft.RunEmpleado] = struct{}{}
	}

	// Crear un mapa para contar los contratos por unidad mayor
	contratoCounts := make(map[string]int)
	for _, contrato := range contratos {
		contratoCounts[contrato.UnidadMayor]++
	}

	// Crear un mapa para contar los profesores de pipelsofts por unidad mayor
	pipelsoftCounts := make(map[string]int)
	for _, contrato := range contratos {
		// Verificar si el RunDocente del contrato está en el mapa de pipelsofts
		if _, encontrado := pipelsoftMap[contrato.RunDocente]; encontrado {
			// Si hay coincidencia, sumar 1 en pipelsoftCounts
			pipelsoftCounts[contrato.UnidadMayor]++
		}
		// Si no se encuentra coincidencia, no hacer nada (seguir con el siguiente contrato)
	}

	// Devolver los resultados
	return contratoCounts, pipelsoftCounts, nil
}

// ProfesorUnidadMayorNOPaf cuenta la cantidad de contratos por unidad mayor
// y los profesores en Pipelsoft, sumando 1 cuando no hay coincidencias.
func (s *ContratoService) ProfesorUnidadMayorNOPaf() (map[string]int, map[string]int, error) {
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
		// Inicializar la variable `found` como `false`
		found := false
		for _, pipelsoft := range pipelsofts {
			// Si el RunDocente del contrato coincide con RunEmpleado del pipelsoft
			if contrato.RunDocente == pipelsoft.RunEmpleado {
				found = true
				break // No hace falta hacer nada, si hay coincidencia, se detiene la búsqueda
			}
		}

		// Si no se encontró una coincidencia, se incrementa el contador para esa UnidadMayor
		if !found {
			pipelsoftCounts[contrato.UnidadMayor]++
		}
	}

	return contratoCounts, pipelsoftCounts, nil
}

func (s *ContratoService) GetPafByUnidadMayor(nombreUnidadMayor string) (map[string][]models.Pipelsoft, error) {
	// Definir los estados posibles
	estados := []string{
		"A1", "A2", "A3", "B1", "B9", "C1D", "C9D", "F1", "F9", "A9",
	}

	// Buscar las PAF correspondientes a la unidad mayor
	var pipelsofts []models.Pipelsoft
	result := s.DB.Where("nombre_unidad_mayor = ?", nombreUnidadMayor).Find(&pipelsofts)
	if result.Error != nil {
		log.Println("Error al buscar las PAF:", result.Error)
		return nil, result.Error
	}

	// Log para verificar cuántos registros se obtuvieron
	log.Printf("Cantidad de PAF encontradas para '%s': %d\n", nombreUnidadMayor, len(pipelsofts))
	if len(pipelsofts) == 0 {
		log.Println("No se encontraron PAF para esta unidad mayor.")
	}

	// Crear un mapa para almacenar las PAF divididas por estado
	pafPorEstado := make(map[string][]models.Pipelsoft)

	// Inicializar las listas de PAF por estado, incluso si están vacías
	for _, estado := range estados {
		pafPorEstado[estado] = []models.Pipelsoft{}
	}

	// Iterar sobre las PAF encontradas y clasificarlas por estado
	for _, pipelsoft := range pipelsofts {
		// Log para verificar cada registro de PAF
		log.Printf("PAF encontrada: %+v\n", pipelsoft)

		// Verificar si el estado de la PAF es uno de los estados definidos
		if _, exists := pafPorEstado[pipelsoft.CodEstado]; exists {
			// Añadir la PAF al estado correspondiente
			pafPorEstado[pipelsoft.CodEstado] = append(pafPorEstado[pipelsoft.CodEstado], pipelsoft)
			log.Printf("PAF añadida al estado %s: %+v\n", pipelsoft.CodEstado, pipelsoft)
		} else {
			log.Printf("Estado no reconocido: %s\n", pipelsoft.CodEstado)
		}
	}

	// Log para verificar el mapa resultante
	log.Println("Mapa de PAF por estado:", pafPorEstado)

	// Devolver el mapa con las PAF por estado
	return pafPorEstado, nil
}
