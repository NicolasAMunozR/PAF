package models

import (
	"time"

	"gorm.io/gorm"
)

// HistorialPafAceptadas representa el historial de aceptaciones de PAFs con información relevante.
type HistorialPafAceptadas struct {
	gorm.Model

	// Campos de Contrato
	Run                 string    `gorm:"type:text;not null"`
	IdPaf               int       `gorm:"type:numeric;not null"`
	FechaInicioContrato time.Time `gorm:"type:timestamp;not null"`
	FechaFinContrato    time.Time `gorm:"type:timestamp;not null"`
	CodigoAsignatura    string    `gorm:"type:text;not null"`
	NombreAsignatura    string    `gorm:"type:text;not null"`
	CantidadHoras       int       `gorm:"type:int;not null"`
	Jerarquia           string    `gorm:"type:text;not null"`
	Calidad             string    `gorm:"type:text;not null"`

	// Campos de Proceso
	EstadoProceso string `gorm:"type:text;not null"`

	//solo puede ser 0 o 1 e indica cuando se detecta que se modifico, si es un cero no se modifico, si es un 1 se modifico
	CodigoModificacion int `gorm:"type:int;null" json:"codigo_modificacion"`

	// Código que indica el tipo de modificación (puede ser nulo)
	// 0 = todo bien, 1 = modificación, 2 = eliminada, 3 = rechazada
	BanderaModificacion int `gorm:"type:int;null" json:"bandera_modificacion"`

	// Descripción detallada de la modificación (puede ser nulo)
	DescripcionModificacion *string `gorm:"type:text;null" json:"descripcion_modificacion"`

	ProfesorRun              string `json:"run"`                      // RUN (identificador único)
	Semestre                 string `json:"semestre"`                 // Semestre de la asignatura
	Tipo                     string `json:"tipo"`        // Código de la asignatura
	ProfesorCodigoAsignatura string `json:"codigo_asignatura"`        // Nombre de la asignatura
	Seccion                  string `json:"seccion"`                  // Sección o grupo
	Cupo                     int    `json:"cupo"`                     // Capacidad del grupo
	Bloque                   string `json:"bloque" gorm:"type:jsonb"` // Bloque horario

	BanderaAceptacion int `gorm:"type:int;null"`
}
