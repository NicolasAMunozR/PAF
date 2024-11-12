package models

import (
	"gorm.io/gorm"
)

// Persona representa los datos de una persona
type Persona struct {
	gorm.Model             // Incluye ID, CreatedAt, UpdatedAt, DeletedAt
	Run             string `gorm:"type:text;not null"`
	Nombres         string `gorm:"type:text;not null"`
	PrimerApellido  string `gorm:"type:text;not null"`
	SegundoApellido string `gorm:"type:text;not null"`
	Correo          string `gorm:"type:text;not null;unique"`
}
