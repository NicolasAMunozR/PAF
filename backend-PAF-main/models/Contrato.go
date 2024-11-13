package models

import (
	"gorm.io/gorm"
)

type Contrato struct {
	gorm.Model
	CodigoPAF        string `gorm:"type:text;not null"`
	FechaInicio      string `gorm:"type:date;not null"`
	FechaFin         string `gorm:"type:date;not null"`
	CodigoAsignatura string `gorm:"type:text;not null"`
	NombreAsignatura string `gorm:"type:text;not null"`
	CantidadHoras    int    `gorm:"type:int;not null"`
	Jerarquia        string `gorm:"type:text;not null"`
	Calidad          string `gorm:"type:text;not null"`
}
