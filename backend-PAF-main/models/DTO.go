package models

import "time"

type PipelsoftDTO struct {
	RunEmpleado          string    `json:"run_empleado"`
	Nombres              string    `json:"nombres"`
	PrimerApp            string    `json:"primer_app"`
	SegundoApp           string    `json:"segundo_app"`
	NombreUnidadMayor    string    `json:"nombre_unidad_mayor"`
	NombreUnidadMenor    string    `json:"nombre_unidad_menor"`
	IdPaf                int       `json:"id_paf"`
	FechaInicioContrato  time.Time `json:"fecha_inicio_contrato"`
	FechaFinContrato     time.Time `json:"fecha_fin_contrato"`
	NombreAsignatura     string    `json:"nombre_asignatura"`
	HorasAsignatura      int       `json:"horas_asignatura"`
	CantidadHorasPaf     int       `json:"cantidad_horas_paf"`
	Jerarquia            string    `json:"jerarquia"`
	Semestre             string    `json:"semestre"`
	UltimaModificacion   time.Time `json:"ultima_modificacion"`
	Categoria            string    `json:"categoria"`
	CodEstado            string    `json:"cod_estado"`
	DesEstado            string    `json:"des_estado"`
	Llave                string    `json:"llave"`
	Veces                string    `json:"veces"`
	CodigoAsignaturaList []string  `json:"codigo_asignatura_list"` // Campo adicional
}
