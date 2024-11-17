package models

type DatosCombinados struct {
	PipelsoftData Pipelsoft  `json:"pipelsoft_data"`
	ProfesorData  ProfesorDB `json:"profesor_data"`
}
