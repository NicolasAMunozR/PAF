package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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

func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF int, profesor models.ProfesorDB, bloque []string, cod_asignatura_paf string, comentario string) (*models.HistorialPafAceptadas, error) {
	// Parsear los bloques en una lista de BloqueDTO
	bloquesDTO, err := parseBloques(bloque)
	if err != nil {
		return nil, fmt.Errorf("error al procesar el bloque: %w", err)
	}

	// Serializar los bloques a JSON
	bloquesJSON, err := json.Marshal(bloquesDTO)
	if err != nil {
		return nil, fmt.Errorf("error al serializar los bloques a JSON: %w", err)
	}

	// Extraer el codigo_asignatura del JSON de bloques
	var bloquesParsed []struct {
		CodigoAsignatura string `json:"codigoAsignatura"`
	}

	// Deserializar el JSON de bloques
	if err := json.Unmarshal(bloquesJSON, &bloquesParsed); err != nil {
		return nil, fmt.Errorf("error al deserializar bloquesJSON: %w", err)
	}

	// Verificar si hemos obtenido el codigo_asignatura
	if len(bloquesParsed) == 0 {
		return nil, fmt.Errorf("no se encontró el codigo_asignatura en los bloques")
	}

	// Iniciar una transacción para garantizar consistencia
	tx := s.DB.Begin()
	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	// Obtener los valores de Pipelsoft
	var pipelsoft models.Pipelsoft
	if err := tx.Where("id_paf = ? and codigo_asignatura = ?", codigoPAF, cod_asignatura_paf).First(&pipelsoft).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al obtener datos de Pipelsoft: %w", err)
	}

	// Verificar si ya existe un registro con el Código PAF
	var historialExistente models.HistorialPafAceptadas
	if err := tx.Where("llave = ?", pipelsoft.Llave).First(&historialExistente).Error; err == nil {
		// Si el registro existe, eliminarlo
		if err := tx.Delete(&historialExistente).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error al eliminar el historial existente: %w", err)
		}
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, fmt.Errorf("error al buscar historial existente: %w", err)
	}

	// Crear el nuevo registro de historial sin cambiar el estado
	historial := models.HistorialPafAceptadas{
		Run:                 profesor.RUN,
		IdPaf:               codigoPAF,
		FechaInicioContrato: pipelsoft.FechaInicioContrato,
		FechaFinContrato:    pipelsoft.FechaFinContrato,
		CodigoAsignatura:    pipelsoft.CodigoAsignatura,
		NombreAsignatura:    pipelsoft.NombreAsignatura,
		CantidadHoras:       pipelsoft.HorasAsignatura,
		DesEstado:           pipelsoft.DesEstado,
		SemestrePaf:         pipelsoft.Semestre,
		Jerarquia:           pipelsoft.Jerarquia,
		UltimaModificacion:  pipelsoft.UltimaModificacion,

		EstadoProceso:           pipelsoft.CodEstado, // Mantener el estado actual
		Llave:                   pipelsoft.Llave,
		CodigoModificacion:      0,
		BanderaModificacion:     0,
		DescripcionModificacion: nil,
		ProfesorRun:             profesor.RUN,
		Semestre:                profesor.Semestre,
		Tipo:                    profesor.Tipo,
		Seccion:                 profesor.Seccion,
		Cupo:                    profesor.Cupo,
		Bloque:                  bloquesJSON, // Bloques en JSON
		BanderaAceptacion:       0,
		Comentario:              comentario,
		// añadir el parametro semestre_inicio_paf cuando sepamos a lo que se refiere.
	}

	// Insertar el nuevo historial en la base de datos
	if err := tx.Create(&historial).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al crear el historial: %w", err)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	log.Println("Nuevo registro creado con éxito")
	return &historial, nil
}

func parseBloques(bloquesRaw []string) ([]models.BloqueDTO, error) {
	var bloques []models.BloqueDTO
	for _, bloque := range bloquesRaw {
		partes := strings.Fields(bloque) // Separar por espacios
		if len(partes) < 5 {             // Asegúrate de que haya al menos 5 elementos
			return nil, fmt.Errorf("bloque mal formado: %s", bloque)
		}

		// Convertir el cuarto elemento (cupos) a entero
		cupos, err := strconv.Atoi(partes[3])
		if err != nil {
			return nil, fmt.Errorf("error al convertir cupos a entero en bloque '%s': %w", bloque, err)
		}

		// El semestre es el primer elemento
		semestre := partes[0]

		// Crear el objeto BloqueDTO
		bloques = append(bloques, models.BloqueDTO{
			Semestre:         semestre,  // Asignar el semestre al inicio
			CodigoAsignatura: partes[1], // Asignar el código de la asignatura
			Seccion:          partes[2], // Asignar la sección
			Cupos:            cupos,     // Asignar los cupos
			Bloques:          partes[4], // Asignar los bloques (último elemento)
		})
	}
	return bloques, nil
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
func (s *HistorialPafAceptadasService) EliminarHistorial(id int64) error {
	// Eliminar directamente usando la condición
	return s.DB.Where("id = ?", id).Delete(&models.HistorialPafAceptadas{}).Error
}

func (s *HistorialPafAceptadasService) ObtenerTodosLosHistoriales() ([]models.HistorialPafAceptadas, error) {
	var historiales []models.HistorialPafAceptadas
	// Obtener todos los historiales con 'bandera_aceptacion' igual a 1 desde DBPersonal
	if err := s.DB.Where("bandera_aceptacion = ?", 1).Find(&historiales).Error; err != nil {
		return nil, err
	}
	return historiales, nil
}

func (s *HistorialPafAceptadasService) ActualizarBanderaAceptacion(codigoPAF string, nuevaBanderaAceptacion int) error {
	// Definir el orden de los estados
	estadoSiguiente := map[string]string{
		"A1":  "A2",
		"A2":  "A3",
		"A3":  "B1",
		"B1":  "C1D",
		"B9":  "C9D",
		"C1D": "F1",
		"C9D": "A9",
		"F1":  "A9", // Fin del proceso, no hay un estado siguiente válido
		"F9":  "A9",
		"A9":  "A9", // Estado terminal
	}

	// Iniciar una transacción para garantizar consistencia
	tx := s.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", tx.Error)
	}

	// Buscar el historial correspondiente al codigoPAF
	var historial models.HistorialPafAceptadas
	err := tx.Where("id_paf = ?", codigoPAF).First(&historial).Error
	if err != nil {
		tx.Rollback() // Rollback si ocurre un error al buscar el historial
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("historial con codigoPAF %s no encontrado", codigoPAF)
		}
		return fmt.Errorf("error al buscar historial con codigoPAF %s: %w", codigoPAF, err)
	}

	// Actualizar el valor de BanderaAceptacion
	historial.BanderaAceptacion = nuevaBanderaAceptacion

	// Determinar el siguiente estado del proceso
	estadoActual := historial.EstadoProceso
	if nuevoEstado, existe := estadoSiguiente[estadoActual]; existe {
		historial.EstadoProceso = nuevoEstado
	} else {
		// Rollback si el estado actual no tiene un mapeo definido
		tx.Rollback()
		return fmt.Errorf("estado actual %s no tiene un estado siguiente definido", estadoActual)
	}

	// Guardar los cambios en la base de datos
	if err := tx.Save(&historial).Error; err != nil {
		tx.Rollback() // Rollback si no se puede guardar el historial
		return fmt.Errorf("error al actualizar historial para codigoPAF %s: %w", codigoPAF, err)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	// Log exitoso
	log.Printf("BanderaAceptacion actualizada a %d y EstadoProceso cambiado a %s para el historial con codigoPAF %s", nuevaBanderaAceptacion, historial.EstadoProceso, codigoPAF)
	return nil
}
