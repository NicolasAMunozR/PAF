package models

type Contrato struct {
	// Identificador del docente, ejemplo: "12345678-9"
	RunDocente string `gorm:"column:run_docente;type:varchar(12);not null" json:"runDocente"`

	// Unidad mayor que contrata al empleado, ejemplo: "FACULTAD DE HUMANIDADES"
	UnidadMayor string `gorm:"column:unidad_mayor;type:varchar(100);not null" json:"unidadMayor"`

	// Unidad menor que contrata al empleado, ejemplo: "DETPO AGRICULTURA"
	UnidadMenor string `gorm:"column:unidad_menor;type:varchar(100);not null" json:"unidadMenor"`

	// Planta del docente, ejemplo: "PROFESOR HORAS CLASES"
	Planta string `gorm:"column:planta;type:varchar(50);not null" json:"planta"`

	// Horas de contrato, ejemplo: 40
	Horas int `gorm:"column:horas;type:int;not null" json:"horas"`
}
