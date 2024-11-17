package models

// ProfesorDB representa la estructura de un profesor con su asignatura.
type ProfesorDB struct {
	RUN              string `json:"run"`               // RUN (identificador único)
	Semestre         string `json:"semestre"`          // Semestre de la asignatura
	CodigoAsignatura string `json:"codigo_asignatura"` // Código de la asignatura
	NombreAsignatura string `json:"nombre_asignatura"` // Nombre de la asignatura
	Seccion          string `json:"seccion"`           // Sección o grupo
	Cupo             int    `json:"cupo"`              // Capacidad del grupo
	Dia              string `json:"dia"`               // Día de la asignatura
	Bloque           string `json:"bloque"`            // Bloque horario
}
