package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// BloqueDTO representa un bloque de información relacionado con horarios y asignaturas.
type BloqueDTO struct {
	Semestre         string `json:"semestre"`
	CodigoAsignatura string `json:"codigoAsignatura"` // Código de la asignatura
	Seccion          string `json:"seccion"`          // Sección
	Cupos            int    `json:"cupos"`            // Cupos disponibles
	Bloques          string `json:"bloques"`          // Información de bloques, e.g., "V7-J1-M5"

}

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
	SemestrePaf         string    `gorm:"type:text"`

	// Campos de Proceso
	EstadoProceso string `gorm:"type:text;not null"`

	// Indicadores de Modificación
	CodigoModificacion      int     `gorm:"type:int;null" json:"codigo_modificacion"`
	BanderaModificacion     int     `gorm:"type:int;null" json:"bandera_modificacion"`
	DescripcionModificacion *string `gorm:"type:text;null" json:"descripcion_modificacion"`

	// Información del Profesor
	ProfesorRun              string `json:"run"`               // RUN (identificador único)
	Semestre                 string `json:"semestre"`          // Semestre de la asignatura
	Tipo                     string `json:"tipo"`              // tipo
	ProfesorCodigoAsignatura string `json:"codigo_asignatura"` // codigo de la asignatura
	Seccion                  string `json:"seccion"`           // Sección o grupo
	Cupo                     int    `json:"cupo"`              // Capacidad del grupo

	// Bloques: Lista de objetos BloqueDTO almacenados como JSONB
	Bloque json.RawMessage `gorm:"type:jsonb"`

	BanderaAceptacion int    `gorm:"type:int;null"`
	Llave             string `gorm:"type:text"`
}
