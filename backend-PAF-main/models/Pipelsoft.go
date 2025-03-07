// models/pipelsoft.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Pipelsoft struct {
	gorm.Model
	// Campos de Persona
	RunEmpleado string `gorm:"type:text"`
	Nombres     string `gorm:"type:text"`
	PrimerApp   string `gorm:"type:text"`
	SegundoApp  string `gorm:"type:text"`

	// Campos de UnidadContratante
	NombreUnidadMayor string `gorm:"type:text"`
	NombreUnidadMenor string `gorm:"type:text"`

	// Campos de Contrato
	IdPaf               int       `gorm:"type:int"`
	FechaInicioContrato time.Time `gorm:"type:timestamp"`
	FechaFinContrato    time.Time `gorm:"type:timestamp"`
	CodigoAsignatura    string    `gorm:"type:text"`
	NombreAsignatura    string    `gorm:"type:text"`
	HorasAsignatura     int       `gorm:"type:int"`
	CantidadHorasPaf    int       `gorm:"type:int"`
	Jerarquia           string    `gorm:"type:text"`
	Semestre            string    `gorm:"type:text"`
	UltimaModificacion  time.Time `gorm:"type:timestamp"`

	Categoria string `gorm:"type:text"`

	// Campos de Proceso
	CodEstado                    string `gorm:"type:text"`
	DesEstado                    string `gorm:"type:text"`
	Llave                        string `gorm:"type:text"`
	Veces                        string `gorm:"type:text"`
	Planta                       string `gorm:"type:text"`
	NumeroCentroDeCostos         string `gorm:"type:text"`
	LugarNacimiento              string `gorm:"type:text"`
	Nacionalidad                 string `gorm:"type:text"`
	Domicilio                    string `gorm:"type:text"`
	Correo                       string `gorm:"type:text"`
	Titulo                       string `gorm:"type:text"`
	Institucion                  string `gorm:"type:text"`
	FechaObtencion               string `gorm:"type:text"`
	NumeroSemestre               string `gorm:"type:text"`
	GradoAcademico               string `gorm:"type:text"`
	InstitucionGradoAcademico    string `gorm:"type:text"`
	FechaObtencionGradoAcademico string `gorm:"type:text"`
	TipoIngreso                  string `gorm:"type:text"`
	Cargo                        string `gorm:"type:text"`
	Nivel                        string `gorm:"type:text"`
	Grado                        string `gorm:"type:text"`
	Rango                        string `gorm:"type:text"`
	Funcion                      string `gorm:"type:text"`
	Asignatura                   string `gorm:"type:text"`
	NumeroHoras                  string `gorm:"type:text"`
	Calidad                      string `gorm:"type:text"`
	LugarDesempeño               string `gorm:"type:text"`
	CargoOtroPublico             string `gorm:"type:text"`
	GradoOtroPublico             string `gorm:"type:text"`
	NivelOtroPublico             string `gorm:"type:text"`
	RangoOtroPublico             string `gorm:"type:text"`
	NumeroHorasOtroPublico       int    `gorm:"type:int"`
	CalidadOtroPublico           string `gorm:"type:text"`
	FechaYHoraNacimiento         string `gorm:"type:text"`

	//añadir dos campos solicitados
	// y posiblemente añadir mas logica de negocio

}
