package models

import (
	"time"

	"gorm.io/gorm"
)

// Pipelsoft representa todos los campos de Contrato, Persona, Proceso, y UnidadContratante en un solo modelo
type HistorialPafAceptadas struct {
	gorm.Model

	Run string `gorm:"type:text;not null"`

	CodigoPAF string `gorm:"type:text;not null"`

	FechaAceptacionPaf time.Time `gorm:"type:timestamp;not null"`
}
