package models

import "gorm.io/gorm"


type InfoPasada struct {
	SemestrePasado string `json:"semestre"`
	CupoPasado  int    `json:"cupo_pasado"` 
}

type ProfesorDB struct {
	gorm.Model
	RUN              string       `json:"run"`               // RUN (identificador único)
	Semestre         string       `json:"semestre"`          // Semestre de la asignatura
	Tipo             string       `json:"tipo"`              // Código de la asignatura
	NombreAsignatura string       `json:"nombre_asignatura"` // Nombre de la asignatura
	Seccion          string       `json:"seccion"`           // Sección o grupo
	Cupo             int          `gorm:"type:int"`          // Capacidad del grupo
	Bloque           string       `json:"bloque"`            // Bloque horario
	CodigoAsignatura string       `json:"codigo_asignatura"`
	CantidadHoras    int          `gorm:"type:int"`
	SemestrePasado  string        `json:"semestre_pasado"`
	CupoPasado 		int           `gorm:"type:int"` 
}