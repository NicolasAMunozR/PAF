package models

import (
	"time"
)

type Auditoria struct {
	Rut                    string    `gorm:"type:text"`
	TipoDeModificacion     string    `gorm:"type:text"`
	DescipcionModificacion string    `gorm:"type:text"`
	FechaModificacion      time.Time `gorm:"type:timestamp"`
}
