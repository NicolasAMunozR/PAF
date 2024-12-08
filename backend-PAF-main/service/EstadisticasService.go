package service

import (
	"fmt"
	"strings"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models" // Cambiar con el paquete donde se encuentran los modelos
	"gorm.io/gorm"
)

type EstadisticasResponse struct {
	TotalPipelsoft             int64                   `json:"total_pipelsoft"`
	TotalPipelsoftUnicos       int64                   `json:"total_pipelsoft_unicos"`
	PorcentajeUnicos           float64                 `json:"porcentaje_unicos"`
	EstadoProcesoCount         map[string]int          `json:"estado_proceso_count"`
	EstadoProcesoPct           map[string]float64      `json:"estado_proceso_pct"`
	TotalProfesores            int64                   `json:"total_profesores"`
	ProfesoresNoEnPipelsoft    int64                   `json:"profesores_no_en_pipelsoft"`
	ProfesoresNoEnPipelsoftPct float64                 `json:"profesores_no_en_pipelsoft_pct"`
	UnidadesMenores            []string                `json:"unidades_menores"`        // Campo para las unidades menores
	UnidadMenorProfesores      []UnidadMenorProfesores `json:"unidad_menor_profesores"` // Campo para los profesores asociados a las unidades menores
}

type UnidadMenorProfesores struct {
	UnidadMenor string   `json:"unidad_menor"`
	Profesores  []string `json:"profesores"`
}

// EstadisticasService define los métodos para obtener estadísticas de las tablas
type EstadisticasService struct {
	DB *gorm.DB
}

// Constructor del servicio
func NewEstadisticasService(dbPersonal *gorm.DB) *EstadisticasService {
	return &EstadisticasService{
		DB: dbPersonal,
	}
}

// ObtenerEstadisticas obtiene las estadísticas generales de los profesores y Pipelsofts
func (s *EstadisticasService) ObtenerEstadisticas() (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Contar los RUN únicos en la tabla profesor_dbs
	if err := s.DB.Model(&models.ProfesorDB{}).
		Distinct("run").Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores únicos por RUN: %w", err)
	}
	// Contar todos los registros en la tabla pipelsofts
	if err := s.DB.Model(&models.Pipelsoft{}).Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts
	if err := s.DB.Model(&models.Pipelsoft{}).Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts: %w", err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los Run de los profesores que no existen en pipelsofts
	var profesoresNoEnPipelsoft int64
	if err := s.DB.Table("contratos").
		Where("run_docente NOT IN (SELECT run_empleado FROM pipelsofts)").
		Count(&profesoresNoEnPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores que no están en pipelsofts: %w", err)
	}

	// Asignar el valor a la estructura
	resp.ProfesoresNoEnPipelsoft = profesoresNoEnPipelsoft

	// Calcular porcentaje de profesores no en pipelsoft
	if resp.TotalProfesores > 0 {
		resp.ProfesoresNoEnPipelsoftPct = float64(profesoresNoEnPipelsoft) / float64(resp.TotalProfesores) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso (código de estado como string)
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador", "Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		if err := s.DB.Model(&models.Pipelsoft{}).Where("des_estado = ?", estado).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s: %w", estado, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	return &resp, nil
}

// ContarRegistrosPorNombreUnidadMayor cuenta los registros en Pipelsoft que coinciden con el nombre de la unidad Mayor
func (s *EstadisticasService) ContarRegistrosPorNombreUnidadMayor(nombreUnidadMayor string) (int64, error) {
	var count int64

	// Contar los registros que coinciden con el nombre de unidad Mayor
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_Mayor = ?", nombreUnidadMayor).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error al contar los registros por nombre de unidad Mayor: %w", err)
	}

	return count, nil
}

// ObtenerFrecuenciaNombreUnidadMayor devuelve un mapa con cada NombreUnidadMayor y la cantidad de veces que aparece
func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMayor() (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		Conteo            int
	}

	// Realizar la consulta agrupada
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(*) as conteo").
		Group("nombre_unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener la frecuencia de NombreUnidadMayor: %w", err)
	}

	// Convertir los resultados en un mapa
	frecuencia := make(map[string]int)
	for _, resultado := range resultados {
		frecuencia[resultado.NombreUnidadMayor] = resultado.Conteo
	}

	return frecuencia, nil
}

//mostrar paf solo las activas y contrastarlas con la cantidad de profesores totales.

// ContarRegistrosExcluyendoEstados cuenta los registros en Pipelsoft donde cod_estado no sea "F1", "F9" o "A9"
// y calcula el porcentaje dividiéndolo por el total de profesores.
func (s *EstadisticasService) ContarRegistrosExcluyendoEstados() (int64, float64, error) {
	var count int64
	var totalProfesores int64

	// Contar los registros en Pipelsoft donde cod_estado no sea "F1", "F9", o "A9"
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}).
		Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar los registros excluyendo estados: %w", err)
	}
	// Contar el total de RUN únicos en la tabla profesor_dbs
	if err := s.DB.Model(&models.ProfesorDB{}).
		Distinct("run").Count(&totalProfesores).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar el total de profesores únicos por RUN: %w", err)
	}

	// Calcular el porcentaje
	var porcentaje float64
	if totalProfesores > 0 {
		porcentaje = float64(count) / float64(totalProfesores) * 100
	}

	return count, porcentaje, nil
}

func (s *EstadisticasService) ObtenerEstadisticasPorUnidadMayor(unidadMayor string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Validar que el parámetro no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidad-mayor' es obligatorio")
	}

	// Contar los registros en pipelsofts que coincidan con la unidad mayor
	if err := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ?", unidadMayor).Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts para la unidad mayor %s: %w", unidadMayor, err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts para la unidad mayor
	if err := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ?", unidadMayor).
		Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts para la unidad mayor %s: %w", unidadMayor, err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso filtrando por unidad mayor
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador", "Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		if err := s.DB.Model(&models.Pipelsoft{}).
			Where("nombre_unidad_mayor = ? AND des_estado = ?", unidadMayor, estado).
			Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s para la unidad mayor %s: %w", estado, unidadMayor, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	// Contar la cantidad de profesores en la tabla contratos que tengan la misma unidad mayor
	if err := s.DB.Model(&models.Contrato{}).Where("unidad_mayor = ?", unidadMayor).Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores en la tabla contratos para la unidad mayor %s: %w", unidadMayor, err)
	}

	return &resp, nil
}

func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayor(nombreUnidadMayor string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		Conteo            int
	}

	// Validar que se proporciona el parámetro
	if nombreUnidadMayor == "" {
		return nil, fmt.Errorf("el nombre de la unidad mayor es obligatorio")
	}

	// Realizar la consulta agrupada
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(*) as conteo").
		Where("nombre_unidad_mayor = ?", nombreUnidadMayor).
		Group("nombre_unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener la frecuencia de NombreUnidadMenor: %w", err)
	}

	// Convertir los resultados en un mapa
	frecuencia := make(map[string]int)
	for _, resultado := range resultados {
		frecuencia[resultado.NombreUnidadMenor] = resultado.Conteo
	}

	return frecuencia, nil
}

func (s *EstadisticasService) ContarRegistrosPorUnidadMayorConRuns(unidadMayor string) (int64, int64, error) {
	var count int64     // Conteo de registros finales en Pipelsoft
	var totalRUNs int64 // Conteo de RUNs únicos

	// Paso 1: Obtener RUNs únicos de ProfesorDB
	var runsProfesorDB []string
	if err := s.DB.Model(&models.ProfesorDB{}).
		Distinct("RUN"). // Nombre del campo en el modelo de backend
		Pluck("RUN", &runsProfesorDB).Error; err != nil {
		return 0, 0, fmt.Errorf("error al obtener RUNs de ProfesorDB: %w", err)
	}

	// Paso 2: Filtrar Pipelsoft por RUNs coincidentes
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB). // Nombre del campo en la base de datos
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}).
		Where("nombre_unidad_mayor = ?", unidadMayor).
		Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar registros filtrados de Pipelsoft: %w", err)
	}

	// Paso 3: Contar los RUNs únicos coincidentes entre ProfesorDB y Pipelsoft
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB). // Nombre del campo en la base de datos
		Distinct("run_empleado").                   // Nombre del campo en la base de datos
		Count(&totalRUNs).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar RUNs únicos en Pipelsoft: %w", err)
	}

	return count, totalRUNs, nil
}

// NormalizarRun elimina el guion y el número posterior, además de manejar el prefijo "0".
func NormalizarRun(run string) string {
	// Dividir en partes por el guion
	parts := strings.Split(run, "-")
	// Tomar la primera parte y eliminar ceros iniciales
	if len(parts) > 0 {
		return strings.TrimLeft(parts[0], "0")
	}
	return ""
}

// ObtenerRUNUnicosExcluidos obtiene los RUN únicos de ProfesorDB que no están en Pipelsoft.
func (s *EstadisticasService) ObtenerRUNUnicosExcluidos() ([]string, []string, error) {
	// Obtener los RUN únicos de ProfesorDB
	var profesorRuns []string
	if err := s.DB.Model(&models.ProfesorDB{}).Distinct("run").Pluck("run", &profesorRuns).Error; err != nil {
		return nil, nil, err
	}

	// Obtener los RunEmpleado de Pipelsoft
	var pipelsoftEmpleados []string
	if err := s.DB.Model(&models.Pipelsoft{}).Pluck("run_empleado", &pipelsoftEmpleados).Error; err != nil {
		return nil, nil, err
	}

	// Normalizar los RunEmpleado
	normalizados := make(map[string]bool)
	for _, run := range pipelsoftEmpleados {
		normalizados[NormalizarRun(run)] = true
	}

	// Filtrar los RUN únicos de ProfesorDB que no están en Pipelsoft
	var excluidos []string
	for _, run := range profesorRuns {
		if !normalizados[NormalizarRun(run)] {
			excluidos = append(excluidos, run)
		}
	}

	return profesorRuns, excluidos, nil
}

// CompararRuns compara los RUN excluidos con los RunDocente en Contrato.
func (s *EstadisticasService) CompararRuns(runsExcluidos []string) ([]string, int, error) {
	// Obtener todos los RunDocente de Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).Distinct("run_docente").Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, 0, err
	}

	// Crear un mapa para búsqueda rápida de RunDocente
	runDocenteMap := make(map[string]bool)
	for _, run := range runDocentes {
		runDocenteMap[run] = true
	}

	// Filtrar los RUN excluidos que no están en Contrato
	var noEncontrados []string
	for _, run := range runsExcluidos {
		if !runDocenteMap[run] {
			noEncontrados = append(noEncontrados, run)
		}
	}

	return noEncontrados, len(noEncontrados), nil
}

type UnidadMayorConProfesores struct {
	NombreUnidadMayor string
	TotalProfesores   int64
}

func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesores() (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar RUNs presentes en la tabla Pipelsoft
	var runFiltrados []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runDocentes).
		Distinct("run_empleado").
		Pluck("run_empleado", &runFiltrados).Error; err != nil {
		return nil, fmt.Errorf("error al filtrar RUNs presentes en Pipelsoft: %w", err)
	}

	// Paso 3: Consultar unidades mayores y contar RUNs únicos
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runFiltrados).
		Group("nombre_unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores con profesores: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesConProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesConProfesores[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesConProfesores, nil
}

// diferenciaDeSlices devuelve los elementos presentes en el primer slice pero no en el segundo
func diferenciaDeSlices(a, b []string) []string {
	existe := make(map[string]bool)
	for _, item := range b {
		existe[item] = true
	}

	var diferencia []string
	for _, item := range a {
		if !existe[item] {
			diferencia = append(diferencia, item)
		}
	}

	return diferencia
}

// 2
func (s *EstadisticasService) ObtenerUnidadesMayoresSinProfesoresEnPipelsoft() (map[string]int, error) {
	var resultados []struct {
		UnidadMayor     string
		TotalProfesores int
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Contar RUNs únicos por unidad mayor en Contrato
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_mayor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Group("unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores sin profesores en Pipelsoft: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesSinProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesSinProfesores[resultado.UnidadMayor] = resultado.TotalProfesores
	}

	return unidadesSinProfesores, nil
}

// 3
func (s *EstadisticasService) ObtenerUnidadesMayoresSinProfesoresEnPipelsoft_3() (map[string]int, error) {
	var resultados []struct {
		UnidadMayor     string
		TotalProfesores int
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Contar RUNs únicos por unidad mayor en Contrato
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_mayor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Group("unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores sin profesores en Pipelsoft: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesSinProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesSinProfesores[resultado.UnidadMayor] = resultado.TotalProfesores
	}

	return unidadesSinProfesores, nil
}

// 4
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivos() (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Paso 1: Obtener los RUNs únicos de la tabla Contrato
	var runsContrato []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runsContrato).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar registros en Pipelsoft con CodEstado válido y los RUNs obtenidos
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).             // RUNs coincidentes
		Where("cod_estado IN ?", []string{"F1", "F9", "A9"}). // CodEstado filtrado (activos)
		Group("nombre_unidad_mayor").                         // Agrupar por unidad mayor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores con PAF activos: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesConPAFActivos := make(map[string]int)
	for _, resultado := range resultados {
		unidadesConPAFActivos[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesConPAFActivos, nil
}

// 5
func (s *EstadisticasService) ObtenerUnidadesMayoresPorCodEstadoPAF(codEstadoPAF string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Validar que codEstadoPAF no esté vacío
	if codEstadoPAF == "" {
		return nil, fmt.Errorf("el parámetro codEstadoPAF no puede estar vacío")
	}

	// Consulta a la base de datos
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("des_estado = ?", codEstadoPAF).
		Group("nombre_unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores por codEstadoPAF: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMayores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMayores[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesMayores, nil
}

func (s *EstadisticasService) ObtenerEstadisticasPorUnidad(unidadMayor, unidadMenor string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Validar que al menos 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Crear un query base para los filtros
	query := s.DB.Model(&models.Pipelsoft{})

	// Filtrar por 'unidadMayor'
	query = query.Where("nombre_unidad_mayor = ?", unidadMayor)

	// Si 'unidadMenor' está vacío, obtenemos las unidades menores asociadas a la 'unidadMayor'
	if unidadMenor != "" {
		query = query.Where("nombre_unidad_menor = ?", unidadMenor)
	}

	// Contar los registros en pipelsofts
	if err := query.Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts
	if err := query.Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts: %w", err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador",
		"Rechazada por Validador", "Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado",
		"Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		if err := query.Where("des_estado = ?", estado).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s: %w", estado, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	// Obtener las unidades menores asociadas a la 'unidadMayor' filtrada
	var unidadesMenores []string
	if err := query.Distinct("nombre_unidad_menor").Pluck("nombre_unidad_menor", &unidadesMenores).Error; err != nil {
		return nil, fmt.Errorf("error al obtener las unidades menores filtradas: %w", err)
	}

	// Incluir las unidades menores en la respuesta
	resp.UnidadesMenores = unidadesMenores

	// Obtener los profesores asociados a cada unidad menor, filtrando por las mismas condiciones
	unidadMenorProfesores := []UnidadMenorProfesores{}
	for _, unidad := range unidadesMenores {
		var profesores []string
		if err := s.DB.Model(&models.Contrato{}).
			Where("unidad_mayor = ?", unidadMayor).
			Where("unidad_menor = ?", unidad).
			Distinct("run_docente").
			Pluck("run_docente", &profesores).Error; err != nil {
			return nil, fmt.Errorf("error al obtener los profesores para la unidad menor %s: %w", unidad, err)
		}
		unidadMenorProfesores = append(unidadMenorProfesores, UnidadMenorProfesores{
			UnidadMenor: unidad,
			Profesores:  profesores,
		})
	}

	// Incluir los profesores asociados a las unidades menores en la respuesta
	resp.UnidadMenorProfesores = unidadMenorProfesores

	// Contar los profesores en la tabla contratos que coincidan con los filtros
	queryContratos := s.DB.Model(&models.Contrato{}).
		Where("unidad_mayor = ?", unidadMayor)

	if unidadMenor != "" {
		queryContratos = queryContratos.Where("unidad_menor = ?", unidadMenor)
	}

	if err := queryContratos.Distinct("run_docente").Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores en la tabla contratos: %w", err)
	}

	// Contar los profesores que no están en Pipelsoft
	var profesoresNoEnPipelsoft int64
	if err := queryContratos.Where("run_docente NOT IN (SELECT run_empleado FROM pipelsofts)").Count(&profesoresNoEnPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores que no están en pipelsofts: %w", err)
	}
	resp.ProfesoresNoEnPipelsoft = profesoresNoEnPipelsoft

	// Calcular porcentaje de profesores no en pipelsoft
	if resp.TotalProfesores > 0 {
		resp.ProfesoresNoEnPipelsoftPct = float64(profesoresNoEnPipelsoft) / float64(resp.TotalProfesores) * 100
	}

	return &resp, nil
}

// 7
func (s *EstadisticasService) ObtenerEstadisticasPorUnidadTOTO(unidadMayor, unidadMenor string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Validar que al menos 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Crear un query base para los filtros
	query := s.DB.Model(&models.Pipelsoft{})
	// Filtrar por 'unidadMayor' y 'cod_estado'
	query = query.Where("nombre_unidad_mayor = ?", unidadMayor)
	query = query.Where("cod_estado NOT IN (?, ?, ?)", "F1", "F9", "A9")

	// Si 'unidadMenor' está especificado, agregar el filtro
	if unidadMenor != "" {
		query = query.Where("nombre_unidad_menor = ?", unidadMenor)
	}

	// Obtener estadísticas generales de Pipelsoft
	if err := query.Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Obtener registros únicos de 'run_empleado' de Pipelsoft
	if err := query.Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts: %w", err)
	}

	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Obtener los estados de proceso
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador",
		"Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	for _, estado := range estados {
		var count int64
		// Filtrar por estado también
		if err := query.Where("des_estado = ?", estado).Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s: %w", estado, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	// Obtener las unidades menores asociadas a la 'unidadMayor'
	var unidadesMenores []string
	// Aplicar los mismos filtros al obtener las unidades menores
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_mayor = ?", unidadMayor).
		Where("cod_estado NOT IN (?, ?, ?)", "F1", "F9", "A9").
		Distinct("nombre_unidad_menor").
		Pluck("nombre_unidad_menor", &unidadesMenores).Error; err != nil {
		return nil, fmt.Errorf("error al obtener las unidades menores: %w", err)
	}

	resp.UnidadesMenores = unidadesMenores

	// Obtener los profesores asociados a cada unidad menor
	unidadMenorProfesores := []UnidadMenorProfesores{}
	for _, unidad := range unidadesMenores {
		var profesores []string
		// Filtrar los profesores según la unidad menor
		if err := s.DB.Model(&models.Contrato{}).
			Where("unidad_menor = ?", unidad).
			Distinct("run_docente").
			Pluck("run_docente", &profesores).Error; err != nil {
			return nil, fmt.Errorf("error al obtener los profesores para la unidad menor %s: %w", unidad, err)
		}
		unidadMenorProfesores = append(unidadMenorProfesores, UnidadMenorProfesores{
			UnidadMenor: unidad,
			Profesores:  profesores,
		})
	}

	resp.UnidadMenorProfesores = unidadMenorProfesores

	// Contar los profesores en la tabla contratos que coincidan con los filtros
	queryContratos := s.DB.Model(&models.Contrato{})
	// Filtrar por 'unidadMayor' y 'unidadMenor' en los contratos también
	if unidadMayor != "" {
		queryContratos = queryContratos.Where("unidad_mayor = ?", unidadMayor)
	}
	if unidadMenor != "" {
		queryContratos = queryContratos.Where("unidad_menor = ?", unidadMenor)
	}

	// Contar los profesores en contratos
	if err := queryContratos.Distinct("run_docente").Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores en la tabla contratos: %w", err)
	}

	// Contar los profesores que no están en Pipelsoft
	var profesoresNoEnPipelsoft int64
	// Asegurarse de filtrar los profesores que no estén en los estados especificados
	if err := queryContratos.Where("run_docente NOT IN (SELECT run_empleado FROM pipelsofts WHERE cod_estado NOT IN (?, ?, ?))", "F1", "F9", "A9").Count(&profesoresNoEnPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores que no están en pipelsofts: %w", err)
	}
	resp.ProfesoresNoEnPipelsoft = profesoresNoEnPipelsoft

	if resp.TotalProfesores > 0 {
		resp.ProfesoresNoEnPipelsoftPct = float64(profesoresNoEnPipelsoft) / float64(resp.TotalProfesores) * 100
	}

	return &resp, nil
}

type UnidadMenorConProfesores struct {
	NombreUnidadMenor string `json:"nombre_unidad_menor"`
	TotalProfesores   int    `json:"total_profesores"`
}

// 8.1
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(unidadMayor string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que unidadMayor no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro unidadMayor no puede estar vacío")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato para la unidad mayor proporcionada
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Where("unidad_mayor = ?", unidadMayor).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Filtrar los RUNs que estén presentes en la tabla Pipelsoft
	var runFiltrados []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runDocentes).
		Distinct("run_empleado").
		Pluck("run_empleado", &runFiltrados).Error; err != nil {
		return nil, fmt.Errorf("error al filtrar RUNs presentes en Pipelsoft: %w", err)
	}

	// Paso 3: Consultar las unidades menores y contar RUNs únicos de los filtrados
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runFiltrados).
		Group("nombre_unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores con profesores: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenoresConProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenoresConProfesores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenoresConProfesores, nil
}

// 8.2
func (s *EstadisticasService) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft() ([]UnidadMenorConProfesores, error) {
	var resultado []UnidadMenorConProfesores

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Contar RUNs únicos por unidad menor en Contrato
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_menor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Group("unidad_menor").
		Scan(&resultado).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores sin profesores en Pipelsoft: %w", err)
	}

	return resultado, nil
}

// 8.3
// 8.3
func (s *EstadisticasService) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(unidadMayor string) ([]UnidadMenorConProfesores, error) {
	var resultado []UnidadMenorConProfesores

	// Validar que el parámetro 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Contar RUNs únicos por unidad menor en Contrato, filtrando por 'unidadMayor'
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_menor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Where("unidad_mayor = ?", unidadMayor). // Filtro por 'unidadMayor'
		Group("unidad_menor").
		Scan(&resultado).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores sin profesores en Pipelsoft: %w", err)
	}

	return resultado, nil
}

// 8.4
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos() ([]UnidadMenorConProfesores, error) {
	var resultado []UnidadMenorConProfesores

	// Paso 1: Obtener los RUNs únicos de la tabla Contrato
	var runsContrato []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runsContrato).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar registros en Pipelsoft con CodEstado válido y los RUNs obtenidos
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).             // RUNs coincidentes
		Where("cod_estado IN ?", []string{"F1", "F9", "A9"}). // CodEstado filtrado (activos)
		Group("nombre_unidad_menor").                         // Agrupar por unidad menor
		Scan(&resultado).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores con PAF activos: %w", err)
	}

	return resultado, nil
}

// 8.5
func (s *EstadisticasService) ObtenerUnidadesMenoresPorCodEstadoPAF(codEstadoPAF string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que codEstadoPAF no esté vacío
	if codEstadoPAF == "" {
		return nil, fmt.Errorf("el parámetro codEstadoPAF no puede estar vacío")
	}

	// Consulta a la base de datos
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("cod_estado = ?", codEstadoPAF).
		Group("nombre_unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores por codEstadoPAF: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenores, nil
}
