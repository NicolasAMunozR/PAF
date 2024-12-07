package models

type Usuarios struct {
	Run    string `gorm:"type:text"`
	Nombre string `gorm:"type:text"`
	Rol    string `gorm:"type:text"`
	//se añadio este parametro
	UnidadMayor string `gorm:"type:text"`
}
