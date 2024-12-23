package models

type Usuarios struct {
	Run    string `gorm:"type:text"`
	Nombre string `gorm:"type:text"`
	Rol    string `gorm:"type:text"`
	//se añadio este parametro
	UnidadMayor string `gorm:"type:text"`
	//se añadio este parametro
	UnidadMenor string `gorm:"type:text"`
	  // Se agregan los siguientes campos
    Cod_umayor int gorm:"type:integer"
    Cod_umenor int gorm:"type:integer"
    acceso int gorm:"type:integer"
    vista_universidad int gorm:"type:integer"
    vista_facultad int gorm:"type:integer"
}
