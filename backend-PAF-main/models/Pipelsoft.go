// models/pipelsoft.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Pipelsoft struct {
	gorm.Model
	// Campos de Persona
	Run             string `gorm:"type:text;not null"`
	Nombres         string `gorm:"type:text;not null"`
	PrimerApellido  string `gorm:"type:text;not null"`
	SegundoApellido string `gorm:"type:text;not null"`
	Correo          string `gorm:"type:text;not null;unique"`

	// Campos de UnidadContratante
	CodigoUnidadContratante string `gorm:"type:varchar(100);not null;unique"`
	NombreUnidadContratante string `gorm:"type:text;not null"`
	NombreUnidadMayor       string `gorm:"type:text"`

	// Campos de Contrato
	CodigoPAF           string `gorm:"type:text;not null"`
	FechaInicioContrato string `gorm:"type:date;not null"`
	FechaFinContrato    string `gorm:"type:date;not null"`
	CodigoAsignatura    string `gorm:"type:text;not null"`
	NombreAsignatura    string `gorm:"type:text;not null"`
	CantidadHoras       int    `gorm:"type:int;not null"`
	Jerarquia           string `gorm:"type:text;not null"`
	Calidad             string `gorm:"type:text;not null"`

	// Campos de Proceso
	EstadoProceso                  int       `gorm:"type:int;not null"`
	FechaUltimaModificacionProceso time.Time `gorm:"type:timestamp;not null"`
}
