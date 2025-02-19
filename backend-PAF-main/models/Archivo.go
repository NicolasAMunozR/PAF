package models

import "time"

type Archivo struct {
	ID                           uint      `gorm:"primaryKey;autoIncrement"` //ID
	UnidadMayor                  string    `gorm:"type:text"`                // sacar del sai o contratos
	UnidadMenor                  string    `gorm:"type:text"`                // misma idea
	NumeroCentroDeCostos         string    `gorm:"type:text"`                // preguntar
	CelulaIdentidad              string    `gorm:"type:text"`                // el Rut del docente
	LugarNacimiento              string    `gorm:"type:text"`                // preguntar
	FechaYHoraNacimiento         time.Time `gorm:"type:timestamp"`           // preguntar
	ApellidoP                    string    `gorm:"type:text"`                // preguntar o sacar desde pipelsoft
	ApellidoM                    string    `gorm:"type:text"`                // preguntar o sacar desde pipelsoft
	Nombres                      string    `gorm:"type:text"`                // sai o pipelsoft
	Nacionalidad                 string    `gorm:"type:text"`                //Preguntar
	Domicio                      string    `gorm:"type:text"`                //preguntar
	Correo                       string    `gorm:"type:text"`                //preguntar
	Titulo                       string    `gorm:"type:text"`                //preguntar
	Institucion                  string    `gorm:"type:text"`                //preguntar
	FechaObtencion               time.Time `gorm:"type:timestamp"`           //preguntar
	NumeroSemestre               string    `gorm:"type:text"`                //preguntar
	GradoAcademico               string    `gorm:"type:text"`                //Preguntar
	InstitucionGradoAcademico    string    `gorm:"type:text"`                //Preguntar
	FechaObtencionGradoAcademico time.Time `gorm:"type:timestamp"`           //preguntar
	TipoIngreso                  string    `gorm:"type:text"`
	// Identificacion de cargo
	Cargo       string `gorm:"type:text"` //preguntar
	Nivel       string `gorm:"type:text"` //Preguntar
	Grado       string `gorm:"type:text"` //Preguntar
	Rango       string `gorm:"type:text"` //Preguntar
	Funcion     string `gorm:"type:text"` //preguntar
	Jerarquia   string `gorm:"type:text"` //preguntar
	Asignatura  string `gorm:"type:text"` //preguntar o desde el sai
	NumeroHoras string `gorm:"type:text"` //SAI
	Categoria   string `gorm:"type:text"` //preguntar
	Calidad     string `gorm:"type:text"` //preguntar
	//cargo o actividad desempeñada en otro cargo publico
	//probablemente todos estos parametros sean opcionales
	LugarDesempeño         string    `gorm:"type:text"`      //preguntar
	CargoOtroPublico       string    `gorm:"type:text"`      //preguntar
	GradoOtroPublico       string    `gorm:"type:text"`      //preguntar
	NivelOtroPublico       string    `gorm:"type:text"`      //preguntar
	RangoOtroPublico       string    `gorm:"type:text"`      //preguntar
	NumeroHorasOtroPublico int       `gorm:"type:int"`       //preguntar
	CalidadOtroPublico     string    `gorm:"type:text"`      //preguntar
	FechaInicioContrato    time.Time `gorm:"type:timestamp"` //preguntar
	FechaFinContrato       time.Time `gorm:"type:timestamp"` //preguntar

	// Campo para almacenar el PDF
	ArchivoPDF []byte `gorm:"type:bytea"`

	Comentario string `gorm:"type:text"` //POSIBLE COMENTARIO DEL USUARIO

	// Relación con los archivos adjuntos
	ArchivosAdjuntos []ArchivoAdjunto `gorm:"foreignKey:ArchivoID"` //
}

// Nuevo modelo para almacenar múltiples archivos
type ArchivoAdjunto struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ArchivoID uint   `gorm:"index"`
	Nombre    string `gorm:"type:text"`
	Datos     []byte `gorm:"type:bytea"`
}
