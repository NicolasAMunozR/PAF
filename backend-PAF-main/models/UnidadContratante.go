package models

import "gorm.io/gorm"

// UnidadContratante representa los datos de una unidad contratante
type UnidadContratante struct {
	gorm.Model
	Codigo            string `gorm:"type:varchar(100);not null;unique"`
	Nombre            string `gorm:"type:text;not null"`
	NombreUnidadMayor string `gorm:"type:text"` // Campo opcional
}
