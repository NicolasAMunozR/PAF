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

// Método para crear un historial
func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF int, profesor models.ProfesorDB, bloque []string) (*models.HistorialPafAceptadas, error) {
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
	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	// Obtener los valores de jerarquía, calidad y estado desde la tabla Pipelsoft
	var pipelsoft models.Pipelsoft
	if err := tx.Where("id_paf = ?", codigoPAF).First(&pipelsoft).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al obtener datos de Pipelsoft: %w", err)
	}

	// Verificar si ya existe un registro con el Código PAF
	var historialExistente models.HistorialPafAceptadas
	if err := tx.Where("llave = ?", pipelsoft.Llave).First(&historialExistente).Error; err == nil {
		// Si el registro existe, eliminarlo
		if err := tx.Delete(&historialExistente).Error; err != nil {
			tx.Rollback() // Rollback si ocurre un error al eliminar
			return nil, fmt.Errorf("error al eliminar el historial existente: %w", err)
		}
		log.Println("Registro existente eliminado correctamente")
	} else if err != gorm.ErrRecordNotFound { // Si ocurre un error distinto de "registro no encontrado"
		tx.Rollback()
		return nil, fmt.Errorf("error al buscar historial existente: %w", err)
	}

	// Determinar el siguiente estado del proceso
	nuevoEstado, existe := estadoSiguiente[pipelsoft.CodEstado]
	if !existe {
		tx.Rollback()
		return nil, fmt.Errorf("estado desconocido: %s", pipelsoft.CodEstado)
	}

	// Crear el nuevo registro de historial
	historial := models.HistorialPafAceptadas{
		Run:                      profesor.RUN,
		IdPaf:                    codigoPAF,
		FechaInicioContrato:      pipelsoft.FechaInicioContrato,
		FechaFinContrato:         pipelsoft.FechaFinContrato,
		CodigoAsignatura:         pipelsoft.CodigoAsignatura,
		NombreAsignatura:         pipelsoft.NombreAsignatura,
		CantidadHoras:            profesor.Cupo,
		Jerarquia:                pipelsoft.Jerarquia,
		EstadoProceso:            nuevoEstado,
		Llave:                    pipelsoft.Llave,
		CodigoModificacion:       0,
		BanderaModificacion:      0,
		DescripcionModificacion:  nil,
		ProfesorRun:              profesor.RUN,
		Semestre:                 profesor.Semestre,
		Tipo:                     profesor.Tipo,
		ProfesorNombreAsignatura: profesor.NombreAsignatura,
		Seccion:                  profesor.Seccion,
		Cupo:                     profesor.Cupo,
		Bloque:                   bloquesJSON, // Usa el JSON serializado
		BanderaAceptacion:        0,
	}

	// Insertar el nuevo historial en la base de datos
	if err := tx.Create(&historial).Error; err != nil {
		tx.Rollback() // Rollback si no se puede crear el historial
		return nil, fmt.Errorf("error al crear el historial: %w", err)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	log.Println("Nuevo registro creado con éxito")
	return &historial, nil
}

// Función auxiliar para parsear bloques
func parseBloques(bloquesRaw []string) ([]models.BloqueDTO, error) {
	var bloques []models.BloqueDTO
	for _, bloque := range bloquesRaw {
		partes := strings.Fields(bloque) // Separar por espacios
		if len(partes) < 4 {
			return nil, fmt.Errorf("bloque mal formado: %s", bloque)
		}

		// Convertir el tercer elemento (cupos) a entero
		cupos, err := strconv.Atoi(partes[2])
		if err != nil {
			return nil, fmt.Errorf("error al convertir cupos a entero en bloque '%s': %w", bloque, err)
		}

		// Crear el objeto BloqueDTO
		bloques = append(bloques, models.BloqueDTO{
			CodigoAsignatura: partes[0],
			Seccion:          partes[1],
			Cupos:            cupos,                         // Valor convertido a int
			Bloques:          strings.Join(partes[3:], "-"), // Combina los bloques restantes
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
func (s *HistorialPafAceptadasService) EliminarHistorial(codigoPAF string) error {
	// Eliminar directamente usando la condición
	return s.DB.Where("id_paf = ?", codigoPAF).Delete(&models.HistorialPafAceptadas{}).Error
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

// ActualizarBanderaAceptacion actualiza la BanderaAceptacion de un historial en HistorialPafAceptadas a partir del codigoPAF
func (s *HistorialPafAceptadasService) ActualizarBanderaAceptacion(codigoPAF string, nuevaBanderaAceptacion int) error {
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

	// Guardar los cambios en la base de datos
	if err := tx.Save(&historial).Error; err != nil {
		tx.Rollback() // Rollback si no se puede guardar el historial
		return fmt.Errorf("error al actualizar BanderaAceptacion para codigoPAF %s: %w", codigoPAF, err)
	}

	// Confirmar la transacción
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	// Log exitoso
	log.Printf("BanderaAceptacion actualizada a %d para el historial con codigoPAF %s", nuevaBanderaAceptacion, codigoPAF)
	return nil
}
