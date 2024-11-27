package service

import (
	"fmt"
	"log"

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

func (s *HistorialPafAceptadasService) CrearHistorial(codigoPAF int, profesor models.ProfesorDB, bloque []string) (*models.HistorialPafAceptadas, error) {
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

	// Verificar si ya existe un registro con el Código PAF
	var historialExistente models.HistorialPafAceptadas
	if err := tx.Where("codigo_paf = ?", codigoPAF).First(&historialExistente).Error; err == nil {
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

	// Obtener los valores de jerarquía, calidad y estado desde la tabla Pipelsoft
	var pipelsoft models.Pipelsoft
	if err := tx.Where("id_paf = ?", codigoPAF).First(&pipelsoft).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error al obtener datos de Pipelsoft: %w", err)
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
		CodigoAsignatura:         profesor.CodigoAsignatura,
		NombreAsignatura:         profesor.NombreAsignatura,
		CantidadHoras:            profesor.Cupo,
		Jerarquia:                pipelsoft.Jerarquia, // Obtenido desde Pipelsoft
		Calidad:                  pipelsoft.Calidad,   // Obtenido desde Pipelsoft
		EstadoProceso:            nuevoEstado,         // Nuevo estado calculado
		CodigoModificacion:       0,
		BanderaModificacion:      0,
		DescripcionModificacion:  nil,
		ProfesorRun:              profesor.RUN,
		Semestre:                 profesor.Semestre,
		ProfesorCodigoAsignatura: profesor.CodigoAsignatura,
		ProfesorNombreAsignatura: profesor.NombreAsignatura,
		Seccion:                  profesor.Seccion,
		Cupo:                     profesor.Cupo,
		Bloque:                   bloque, // Bloque recibido como parámetro
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

// ActualizarBanderaAceptacion actualiza la BanderaAceptacion de un historial en HistorialPafAceptadas a partir del codigoPAF
func (s *HistorialPafAceptadasService) ActualizarBanderaAceptacion(codigoPAF string, nuevaBanderaAceptacion int) error {
	// Iniciar una transacción para garantizar consistencia
	tx := s.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", tx.Error)
	}

	// Buscar el historial correspondiente al codigoPAF
	var historial models.HistorialPafAceptadas
	err := tx.Where("codigo_paf = ?", codigoPAF).First(&historial).Error
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
