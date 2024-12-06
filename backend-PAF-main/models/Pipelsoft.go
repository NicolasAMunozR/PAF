// models/pipelsoft.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Pipelsoft struct {
	gorm.Model
	// Campos de Persona
	RunEmpleado string `gorm:"type:text"`
	Nombres     string `gorm:"type:text"`
	PrimerApp   string `gorm:"type:text"`
	SegundoApp  string `gorm:"type:text"`

	// Campos de UnidadContratante
	NombreUnidadMayor string `gorm:"type:text"`
	NombreUnidadMenor string `gorm:"type:text"`

	// Campos de Contrato
	IdPaf               int       `gorm:"type:int"`
	FechaInicioContrato time.Time `gorm:"type:timestamp"`
	FechaFinContrato    time.Time `gorm:"type:timestamp"`
	CodigoAsignatura    string    `gorm:"type:text"`
	NombreAsignatura    string    `gorm:"type:text"`
	HorasAsignatura     int       `gorm:"type:int"`
	CantidadHorasPaf    int       `gorm:"type:int"`
	Jerarquia           string    `gorm:"type:text"`
	Semestre            string    `gorm:"type:text"`
	UltimaModificacion  time.Time `gorm:"type:Date"`

	Categoria string `gorm:"type:text"`

	// Campos de Proceso
	CodEstado string `gorm:"type:text"`
	DesEstado string `gorm:"type:text"`
	Llave     string `gorm:"type:text"`
	Veces     string `gorm:"type:text"`
}
