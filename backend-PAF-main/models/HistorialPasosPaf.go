package models

import "time"

type HistorialPasosPaf struct {
	IdPaf                string    `gorm:"type:text;not null"`
	RunDocente           string    `gorm:"type:text;not null"`
	EstadoNuevoPaf       string    `gorm:"type:text;not null"`
	CodigoEstadoPaf      string    `gorm:"type:text;not null"`
	FechaLlegadaPaf      time.Time `gorm:"type:timestamp;not null"`
	FechaModificacionPaf time.Time `gorm:"type:timestamp;not null"`
	// no se puede obtener el dato de la ultima modificacion
	// tal vez el procedimiento almacenado ayude, tenlo en cuenta
	// mandar data al front
}
