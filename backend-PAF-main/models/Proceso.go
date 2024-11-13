package models

import (
	"time"
)

// Proceso estructura para almacenar los datos del proceso PAF
type Proceso struct {
	ID                      uint      `gorm:"primary_key"`
	Estado                  string    `gorm:"type:text;not null"`
	FechaUltimaModificacion time.Time `gorm:"type:timestamp;not null"`
}
