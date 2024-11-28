package models

import "gorm.io/gorm"

// ProfesorDB representa la estructura de un profesor con su asignatura.
type ProfesorDB struct {
	gorm.Model
	RUN              string `json:"run"`               // RUN (identificador único)
	Semestre         string `json:"semestre"`          // Semestre de la asignatura
	Tipo             string `json:"tipo"`              // Código de la asignatura
	NombreAsignatura string `json:"nombre_asignatura"` // Nombre de la asignatura
	Seccion          string `json:"seccion"`           // Sección o grupo
	Cupo             int    `gorm:"type:int"`          // Capacidad del grupo
	Bloque           string `json:"bloque"`            // Bloque horario
	CodigoAsignatura string `json:"codigo_asignatura"`
}
