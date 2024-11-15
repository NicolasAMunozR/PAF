package models

import (
	"gorm.io/gorm"
)

type Horario struct {
	gorm.Model
	CodigoPAF string   `json:"codigo_paf" gorm:"type:varchar(100);not null"`
	Horarios  []string `json:"horarios" gorm:"type:text[]"` // Lista de horarios
	Run       string   `json:"run" gorm:"type:varchar(20);not null"`
}
