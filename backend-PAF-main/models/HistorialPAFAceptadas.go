package models

import (
	"gorm.io/gorm"
)

// HistorialPafAceptadas representa el historial de aceptaciones de PAFs con información relevante.
type HistorialPafAceptadas struct {
	gorm.Model

	// Campos de Contrato
	Run                 string `gorm:"type:text;not null"`
	CodigoPAF           string `gorm:"type:text;not null"`
	FechaInicioContrato string `gorm:"type:date;not null"`
	FechaFinContrato    string `gorm:"type:date;not null"`
	CodigoAsignatura    string `gorm:"type:text;not null"`
	NombreAsignatura    string `gorm:"type:text;not null"`
	CantidadHoras       int    `gorm:"type:int;not null"`
	Jerarquia           string `gorm:"type:text;not null"`
	Calidad             string `gorm:"type:text;not null"`

	// Campos de Proceso
	EstadoProceso int `gorm:"type:int;not null"`

	//solo puede ser 0 o 1 e indica cuando se detecta que se modifico, si es un cero no se modifico, si es un 1 se modifico
	CodigoModificacion int `gorm:"type:int;null" json:"codigo_modificacion"`

	// Código que indica el tipo de modificación (puede ser nulo)
	// 0 = todo bien, 1 = modificación, 2 = eliminada, 3 = rechazada
	BanderaModificacion int `gorm:"type:int;null" json:"bandera_modificacion"`

	// Descripción detallada de la modificación (puede ser nulo)
	DescripcionModificacion *string `gorm:"type:text;null" json:"descripcion_modificacion"`
}
