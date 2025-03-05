package models

import (
	"encoding/json"
	"time"

	"github.com/lib/pq"
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
	DeletedAt gorm.DeletedAt `gorm:"index"` // Si no quieres soft delete, quita este campo

	// Campos de Contrato
	Run                 string    `gorm:"type:text"`
	IdPaf               int       `gorm:"type:numeric"`
	FechaInicioContrato time.Time `gorm:"type:timestamp"`
	FechaFinContrato    time.Time `gorm:"type:timestamp"`
	CodigoAsignatura    string    `gorm:"type:text"`
	NombreAsignatura    string    `gorm:"type:text"`
	CantidadHoras       int       `gorm:"type:int"`
	Jerarquia           string    `gorm:"type:text"`
	Calidad             string    `gorm:"type:text"`
	SemestrePaf         string    `gorm:"type:text"`
	DesEstado           string    `gorm:"type:text"`

	// Campos de Proceso
	EstadoProceso string `gorm:"type:text"`

	// Indicadores de Modificación
	CodigoModificacion      int     `gorm:"type:int;null" json:"codigo_modificacion"`
	BanderaModificacion     int     `gorm:"type:int;null" json:"bandera_modificacion"`
	DescripcionModificacion *string `gorm:"type:text;null" json:"descripcion_modificacion"`

	// Información del Profesor
	ProfesorRun              string    `json:"run"`               // RUN (identificador único)
	Semestre                 string    `json:"semestre"`          // Semestre de la asignatura
	Tipo                     string    `json:"tipo"`              // tipo
	ProfesorCodigoAsignatura string    `json:"codigo_asignatura"` // codigo de la asignatura
	Seccion                  string    `json:"seccion"`           // Sección o grupo
	Cupo                     int       `json:"cupo"`              // Capacidad del grupo
	UltimaModificacion       time.Time `gorm:"type:timestamp"`
	Comentario               string    `gorm:"type:text"`

	// Bloques: Lista de objetos BloqueDTO almacenados como JSONB
	Bloque json.RawMessage `gorm:"type:jsonb"`

	BanderaAceptacion int    `gorm:"type:int;null"`
	Llave             string `gorm:"type:text"`
	SemestreInicioPaf string `gorm:"type:text"`
	NombreAsignaturasExtras pq.StringArray `gorm:"type:text[]" json:"comentarios_extras"`
}
