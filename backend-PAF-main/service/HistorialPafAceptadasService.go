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
		CodigoModificacion:  0,
		BanderaModificacion: 0,
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

// Nueva función para actualizar la descripción de las modificaciones en HistorialPafAceptadas
func (service *HistorialPafAceptadasService) ActualizarModificaciones() error {
	var historial []models.HistorialPafAceptadas
	var pipelsoft models.Pipelsoft

	// Obtener todos los registros de HistorialPafAceptadas donde CodigoModificacion != 0
	if err := service.DB.Where("codigo_modificacion != 0").Find(&historial).Error; err != nil {
		return fmt.Errorf("error al obtener historial de paf aceptadas: %w", err)
	}

	// Iterar sobre cada registro de historial y actualizar la descripción de la modificación
	for _, h := range historial {
		// Buscar el registro correspondiente en Pipelsoft usando CodigoPAF
		if err := service.DB.Where("codigo_paf = ?", h.CodigoPAF).First(&pipelsoft).Error; err != nil {
			// Si no se encuentra el registro de Pipelsoft, continuamos con el siguiente
			if err == gorm.ErrRecordNotFound {
				log.Printf("No se encontró Pipelsoft para el CodigoPAF %s", h.CodigoPAF)
				continue
			}
			return fmt.Errorf("error al buscar Pipelsoft: %w", err)
		}

		// Aquí evaluamos qué campos cambiaron. Suponiendo que solo se cambian ciertos campos, como:
		var cambios []string

		// Comparar los valores y detectar qué cambió, agregar a los cambios
		if h.FechaInicioContrato != pipelsoft.FechaInicioContrato {
			cambios = append(cambios, fmt.Sprintf("FechaInicioContrato cambiado de %s a %s", h.FechaInicioContrato, pipelsoft.FechaInicioContrato))
		}
		if h.FechaFinContrato != pipelsoft.FechaFinContrato {
			cambios = append(cambios, fmt.Sprintf("FechaFinContrato cambiado de %s a %s", h.FechaFinContrato, pipelsoft.FechaFinContrato))
		}
		if h.CodigoAsignatura != pipelsoft.CodigoAsignatura {
			cambios = append(cambios, fmt.Sprintf("CodigoAsignatura cambiado de %s a %s", h.CodigoAsignatura, pipelsoft.CodigoAsignatura))
		}
		if h.NombreAsignatura != pipelsoft.NombreAsignatura {
			cambios = append(cambios, fmt.Sprintf("NombreAsignatura cambiado de %s a %s", h.NombreAsignatura, pipelsoft.NombreAsignatura))
		}
		if h.CantidadHoras != pipelsoft.CantidadHoras {
			cambios = append(cambios, fmt.Sprintf("CantidadHoras cambiadas de %d a %d", h.CantidadHoras, pipelsoft.CantidadHoras))
		}
		if h.Jerarquia != pipelsoft.Jerarquia {
			cambios = append(cambios, fmt.Sprintf("Jerarquia cambiada de %s a %s", h.Jerarquia, pipelsoft.Jerarquia))
		}
		if h.Calidad != pipelsoft.Calidad {
			cambios = append(cambios, fmt.Sprintf("Calidad cambiada de %s a %s", h.Calidad, pipelsoft.Calidad))
		}

		// Si hay cambios, agregamos la descripción de los cambios
		if len(cambios) > 0 {
			descripcion := fmt.Sprintf("Modificaciones realizadas: %s", fmt.Sprint(cambios))
			h.DescripcionModificacion = &descripcion

			// Actualizar el registro de HistorialPafAceptadas con la nueva descripción
			if err := service.DB.Save(&h).Error; err != nil {
				return fmt.Errorf("error al actualizar HistorialPafAceptadas: %w", err)
			}
		}
	}

	return nil
}
