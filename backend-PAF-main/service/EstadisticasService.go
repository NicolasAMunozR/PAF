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
func (s *EstadisticasService) ObtenerEstadisticas(semestre string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Normalizar el formato del semestre
	var formato1, formato2 string
	if len(semestre) == 5 {
		if len(semestre) == 5 && semestre[2] == '-' { // Formato "1-23"
			mes := string(semestre[0]) + string(semestre[1])
			anio := "20" + semestre[3:]
			formato1 = semestre
			formato2 = fmt.Sprintf("%s-%s", anio, mes)
		}

	}
	if len(semestre) == 7 && semestre[4] == '-' { // Formato "2023-01"
		anio := semestre[:4]
		mes := semestre[5:]
		formato1 = fmt.Sprintf("%d-%s", removeLeadingZero(mes), anio[2:])
		formato2 = semestre
	} else if len(semestre) == 4 && semestre[1] == '-' { // Formato "1-23"
		mes := "0" + string(semestre[0])
		anio := "20" + semestre[2:]
		formato1 = semestre
		formato2 = fmt.Sprintf("%s-%s", anio, mes)
	}

	// Contar los RUN únicos en la tabla profesor_dbs con filtro por semestre
	if err := s.DB.Model(&models.ProfesorDB{}).
		Where("semestre = ? OR semestre = ? OR semestre = ?", formato1, formato2, semestre).
		Distinct("run").
		Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores únicos por RUN: %w", err)
	}

	// Contar todos los registros en la tabla pipelsofts, aplicando filtro por semestre
	query := s.DB.Model(&models.Pipelsoft{}).Where("semestre = ? OR semestre = ? OR semestre = ?", formato1, formato2, semestre)
	if err := query.Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts, aplicando filtro por semestre
	queryUnicos := s.DB.Model(&models.Pipelsoft{}).Distinct("run_empleado").
		Where("semestre = ? OR semestre = ? OR semestre = ?", formato1, formato2, semestre)
	if err := queryUnicos.Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
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

	// Contar los registros en pipelsofts por cada EstadoProceso
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador",
		"Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	// Contar registros por cada estado, aplicando filtro por semestre
	for _, estado := range estados {
		var count int64
		queryEstado := s.DB.Model(&models.Pipelsoft{}).
			Where("des_estado = ? AND (semestre = ? OR semestre = ? OR semestre = ?)", estado, formato1, formato2, semestre)
		if err := queryEstado.Count(&count).Error; err != nil {
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

// removeLeadingZero elimina un cero inicial de un mes si está presente
func removeLeadingZero(mes string) int {
	if mes[0] == '0' {
		return int(mes[1] - '0')
	}
	return int(mes[0]-'0')*10 + int(mes[1]-'0') // Caso raro, solo por seguridad
}

// ContarRegistrosPorNombreUnidadMayor cuenta los registros en Pipelsoft que coinciden con el nombre de la unidad Mayor
func (s *EstadisticasService) ContarRegistrosPorNombreUnidadMayor(nombreUnidadMayor string, semestre string) (int64, error) {
	var count int64
	fmt.Println("Valor de semestre:", semestre)

	// Iniciar la consulta
	query := s.DB.Model(&models.Pipelsoft{}).
		Where("nombre_unidad_mayor = ?", nombreUnidadMayor)

	// Aplicar el filtro por semestre si se proporciona
	if semestre != "" {
		query = query.Where("semestre = ?", semestre)
	}

	// Contar los registros que coinciden con el nombre de unidad Mayor
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("error al contar los registros por nombre de unidad Mayor: %w", err)
	}

	return count, nil
}

// ObtenerFrecuenciaNombreUnidadMayor devuelve un mapa con cada NombreUnidadMayor y la cantidad de veces que aparece
func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMayor(semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		Conteo            int
	}
	fmt.Println("Valor de semestre:", semestre)

	// Construir y ejecutar la consulta en una sola línea
	err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(*) as conteo").
		Group("nombre_unidad_mayor").
		Where("semestre = ?", semestre).
		Scan(&resultados).Error

	if err != nil {
		return nil, fmt.Errorf("error al obtener la frecuencia de NombreUnidadMayor: %w", err)
	}
	fmt.Println("Valor de semestre:", semestre)

	// Convertir los resultados a un mapa
	frecuencia := make(map[string]int)
	for _, resultado := range resultados {
		frecuencia[resultado.NombreUnidadMayor] = resultado.Conteo
	}

	return frecuencia, nil
}

//mostrar paf solo las activas y contrastarlas con la cantidad de profesores totales.

// ContarRegistrosExcluyendoEstados cuenta los registros en Pipelsoft donde cod_estado no sea "F1", "F9" o "A9"
// y calcula el porcentaje dividiéndolo por el total de profesores.
func (s *EstadisticasService) ContarRegistrosExcluyendoEstados(semestre string) (int64, float64, error) {
	var count int64
	var totalProfesores int64

	// Contar los registros en Pipelsoft donde cod_estado no sea "F1", "F9", o "A9", y aplicar el filtro por semestre si se proporciona
	queryPipelsoft := s.DB.Model(&models.Pipelsoft{}).
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"})
	if semestre != "" {
		queryPipelsoft = queryPipelsoft.Where("semestre = ?", semestre)
	}

	if err := queryPipelsoft.Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar los registros excluyendo estados: %w", err)
	}

	// Contar el total de RUN únicos en la tabla profesor_dbs (sin aplicar filtro por semestre)
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

func (s *EstadisticasService) ObtenerEstadisticasPorUnidadMayor(unidadMayor string, semestre string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Normalizar el formato del semestre
	var formato1, formato2 string
	if len(semestre) == 5 {
		if len(semestre) == 5 && semestre[2] == '-' { // Formato "1-23"
			mes := string(semestre[0]) + string(semestre[1])
			anio := "20" + semestre[3:]
			formato1 = semestre
			formato2 = fmt.Sprintf("%s-%s", anio, mes)
		}
	}

	if len(semestre) == 7 && semestre[4] == '-' { // Formato "2023-01"
		anio := semestre[:4]
		mes := semestre[5:]
		formato1 = fmt.Sprintf("%d-%s", removeLeadingZero(mes), anio[2:])
		formato2 = semestre
	} else if len(semestre) == 4 && semestre[1] == '-' { // Formato "1-23"
		mes := "0" + string(semestre[0])
		anio := "20" + semestre[2:]
		formato1 = semestre
		formato2 = fmt.Sprintf("%s-%s", anio, mes)
	}

	// Validar que el parámetro unidadMayor no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidad-mayor' es obligatorio")
	}

	// Construir una consulta base con los filtros dinámicos
	query := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ?", unidadMayor)
	if semestre != "" {
		query = query.Where("semestre = ? OR semestre = ? OR semestre = ?", semestre, formato1, formato2)
	}

	// Contar los registros en pipelsofts que coincidan con la unidad mayor y semestre
	if err := query.Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts para la unidad mayor %s y semestre %s: %w", unidadMayor, semestre, err)
	}

	// Contar los registros únicos de Run en la tabla pipelsofts para la unidad mayor y semestre
	if err := query.Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de Run en pipelsofts para la unidad mayor %s y semestre %s: %w", unidadMayor, semestre, err)
	}

	// Calcular el porcentaje de registros únicos en Pipelsoft
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar los registros en pipelsofts por cada EstadoProceso filtrando por unidad mayor y semestre
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)

	// Los códigos de estado definidos
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador", "Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	// Contar registros por cada estado
	for _, estado := range estados {
		var count int64
		// Modificar la consulta para filtrar por unidadMayor y semestre además de estado
		queryEstados := s.DB.Model(&models.Pipelsoft{}).Where("des_estado = ? AND nombre_unidad_mayor = ? AND semestre = ?", estado, unidadMayor, semestre)

		// Ejecutar la consulta y contar los registros
		if err := queryEstados.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros de pipelsofts con estado %s para la unidad mayor %s y semestre %s: %w", estado, unidadMayor, semestre, err)
		}

		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	// Contar la cantidad de profesores en la tabla contratos que tengan la misma unidad mayor y semestre
	// cambios
	queryContratos := s.DB.Model(&models.Contrato{}).Where("unidad_mayor = ?", unidadMayor)
	if err := queryContratos.Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores en la tabla contratos para la unidad mayor %s y semestre %s: %w", unidadMayor, semestre, err)
	}

	return &resp, nil
}

func (s *EstadisticasService) ObtenerFrecuenciaNombreUnidadMenorPorUnidadMayor(nombreUnidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		Conteo            int
	}

	fmt.Println("Valor de semestre:", semestre)

	// Validar que se proporcionan los parámetros
	if nombreUnidadMayor == "" || semestre == "" {
		return nil, fmt.Errorf("los parámetros nombreUnidadMayor y semestre son obligatorios")
	}

	// Realizar la consulta agrupada con filtro por semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(*) as conteo").
		Where("nombre_unidad_mayor = ?", nombreUnidadMayor).
		Where("semestre = ?", semestre).
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

func (s *EstadisticasService) ContarRegistrosPorUnidadMayorConRuns(unidadMayor, semestre string) (int64, int64, error) {
	var count int64     // Conteo de registros finales en Pipelsoft
	var totalRUNs int64 // Conteo de RUNs únicos

	// Paso 1: Obtener RUNs únicos de ProfesorDB
	var runsProfesorDB []string
	if err := s.DB.Model(&models.ProfesorDB{}).
		Distinct("RUN"). // Nombre del campo en el modelo de backend
		Pluck("RUN", &runsProfesorDB).Error; err != nil {
		return 0, 0, fmt.Errorf("error al obtener RUNs de ProfesorDB: %w", err)
	}
	fmt.Println("Valor de semestre:", semestre)
	fmt.Println("Valor de semestre:", runsProfesorDB)
	fmt.Println("Valor de unidad_mayor:", unidadMayor)

	// Paso 2: Filtrar Pipelsoft por RUNs coincidentes, unidad mayor y semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB). // Nombre del campo en la base de datos
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}).
		Where("nombre_unidad_mayor = ?", unidadMayor).
		Where("semestre = ?", semestre). // Filtro por semestre
		Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar registros filtrados de Pipelsoft: %w", err)
	}

	// Paso 3: Contar los RUNs únicos coincidentes entre ProfesorDB y Pipelsoft
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB). // Nombre del campo en la base de datos
		Where("semestre = ?", semestre).            // Filtro por semestre
		Distinct("run_empleado").                   // Nombre del campo en la base de datos
		Count(&totalRUNs).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar RUNs únicos en Pipelsoft: %w", err)
	}
	fmt.Println("Valor de count:", count)
	fmt.Println("Valor de total:", totalRUNs)

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

// 1
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesores(semestre string) (map[string]int, error) {
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

	// Paso 2: Filtrar RUNs presentes en la tabla Pipelsoft con el semestre
	var runFiltrados []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ? AND semestre = ?", runDocentes, semestre).
		Distinct("run_empleado").
		Pluck("run_empleado", &runFiltrados).Error; err != nil {
		return nil, fmt.Errorf("error al filtrar RUNs presentes en Pipelsoft: %w", err)
	}

	// Paso 3: Consultar unidades mayores y contar RUNs únicos en Pipelsoft
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ? AND semestre = ?", runFiltrados, semestre).
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
func (s *EstadisticasService) ObtenerUnidadesMayoresSinProfesoresEnPipelsoftPorSemestre(semestre string) (map[string]int, error) {
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

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft filtrados por semestre
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("semestre = ?", semestre).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft filtrados por semestre: %w", err)
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
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorSemestre(semestre string) (map[string]int, error) {
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

	// Paso 2: Filtrar registros en Pipelsoft donde CodEstado no sea "F1", "F9" o "A9",
	// y aplicar el filtro de semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).                 // RUNs coincidentes
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}). // CodEstado que no sean F1, F9 o A9 (inactivos)
		Where("semestre = ?", semestre).                          // Filtro por semestre
		Group("nombre_unidad_mayor").                             // Agrupar por unidad mayor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores con PAF activas: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesConPAFActivas := make(map[string]int)
	for _, resultado := range resultados {
		unidadesConPAFActivas[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesConPAFActivas, nil
}

// 4
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivosPorSemestre(semestre string) (map[string]int, error) {
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

	// Paso 2: Filtrar registros en Pipelsoft con CodEstado válido y aplicar el filtro por semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).             // RUNs coincidentes
		Where("cod_estado IN ?", []string{"F1", "F9", "A9"}). // CodEstado filtrado (activos)
		Where("semestre = ?", semestre).                      // Filtro por semestre
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
func (s *EstadisticasService) ObtenerUnidadesMayoresPorCodEstadoPAF(codEstadoPAF string, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Validar que codEstadoPAF y semestre no estén vacíos
	if codEstadoPAF == "" {
		return nil, fmt.Errorf("el parámetro codEstadoPAF no puede estar vacío")
	}
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro semestre no puede estar vacío")
	}

	// Consulta a la base de datos con filtros por codEstadoPAF y semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("des_estado = ?", codEstadoPAF).
		Where("semestre = ?", semestre).
		Group("nombre_unidad_mayor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores por codEstadoPAF y semestre: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMayores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMayores[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesMayores, nil
}

// 6
func (s *EstadisticasService) ObtenerEstadisticasPorUnidad(unidadMayor, unidadMenor, semestre string) (*EstadisticasResponse, error) {
	var resp EstadisticasResponse

	// Normalizar el formato del semestre
	var formato1, formato2 string
	if len(semestre) == 5 {
		if len(semestre) == 5 && semestre[2] == '-' { // Formato "1-23"
			mes := string(semestre[0]) + string(semestre[1])
			anio := "20" + semestre[3:]
			formato1 = semestre
			formato2 = fmt.Sprintf("%s-%s", anio, mes)
		}
	}

	if len(semestre) == 7 && semestre[4] == '-' { // Formato "2023-01"
		anio := semestre[:4]
		mes := semestre[5:]
		formato1 = fmt.Sprintf("%d-%s", removeLeadingZero(mes), anio[2:])
		formato2 = semestre
	} else if len(semestre) == 4 && semestre[1] == '-' { // Formato "1-23"
		mes := "0" + string(semestre[0])
		anio := "20" + semestre[2:]
		formato1 = semestre
		formato2 = fmt.Sprintf("%s-%s", anio, mes)
	}

	// Validar parámetros obligatorios
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' es obligatorio")
	}
	if unidadMenor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMenor' es obligatorio")
	}

	fmt.Println("Valor de semestre:", semestre)
	// Crear una consulta base con filtros dinámicos
	query := s.DB.Model(&models.Pipelsoft{}).Where("nombre_unidad_mayor = ? AND nombre_unidad_menor = ?", unidadMayor, unidadMenor)
	if semestre != "" {
		query = query.Where("semestre = ? OR semestre =? OR semestre = ?", semestre, formato1, formato2)
	}

	// Contar los registros totales
	if err := query.Count(&resp.TotalPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros en pipelsofts: %w", err)
	}

	// Contar los registros únicos de `run_empleado`
	if err := query.Distinct("run_empleado").Count(&resp.TotalPipelsoftUnicos).Error; err != nil {
		return nil, fmt.Errorf("error al contar los registros únicos de run_empleado en pipelsofts: %w", err)
	}

	// Calcular el porcentaje de registros únicos
	if resp.TotalPipelsoft > 0 {
		resp.PorcentajeUnicos = float64(resp.TotalPipelsoftUnicos) / float64(resp.TotalPipelsoft) * 100
	}

	// Contar registros por estado de proceso
	resp.EstadoProcesoCount = make(map[string]int)
	resp.EstadoProcesoPct = make(map[string]float64)
	estados := []string{
		"Sin Solicitar", "Enviada al Interesado", "Enviada al Validador", "Aprobada por Validador", "Rechazada por Validador",
		"Aprobada por Dir. Pregrado", "Rechazada por Dir. de Pregrado", "Aprobada por RRHH", "Rechazada por RRHH", "Anulada",
	}

	for _, estado := range estados {
		var count int64
		// Modificar la consulta para filtrar por des_estado, unidad_mayor y unidad_menor
		queryEstados := s.DB.Model(&models.Pipelsoft{}).Where("des_estado = ? AND nombre_unidad_mayor = ? AND nombre_unidad_menor = ? AND semestre = ?", estado, unidadMayor, unidadMenor, semestre)

		// Ejecutar la consulta y contar los registros
		if err := queryEstados.Count(&count).Error; err != nil {
			return nil, fmt.Errorf("error al contar los registros con estado '%s', unidad mayor '%s' y unidad menor '%s': %w", estado, unidadMayor, unidadMenor, err)
		}
		resp.EstadoProcesoCount[estado] = int(count)

		// Calcular el porcentaje de registros por cada estado
		if resp.TotalPipelsoft > 0 {
			resp.EstadoProcesoPct[estado] = float64(count) / float64(resp.TotalPipelsoft) * 100
		}
	}

	// Contar la cantidad de profesores en la tabla contratos
	queryContratos := s.DB.Model(&models.Contrato{}).Where("unidad_mayor = ? AND unidad_menor = ?", unidadMayor, unidadMenor)
	if err := queryContratos.Count(&resp.TotalProfesores).Error; err != nil {
		return nil, fmt.Errorf("error al contar los profesores en la tabla contratos: %w", err)
	}

	return &resp, nil
}

// 7
// ContarRegistrosPorUnidadMayorYUnidadMenor cuenta los registros en Pipelsoft filtrados por RUNs de ProfesorDB,
// y adicionalmente por unidad mayor y unidad menor.
func (s *EstadisticasService) ContarRegistrosPorUnidadMayorYUnidadMenor(unidadMayor, unidadMenor, semestre string) (int64, int64, error) {
	var count int64     // Conteo de registros finales en Pipelsoft
	var totalRUNs int64 // Conteo de RUNs únicos

	// Paso 1: Obtener RUNs únicos de ProfesorDB
	var runsProfesorDB []string
	if err := s.DB.Model(&models.ProfesorDB{}).
		Distinct("RUN"). // Nombre del campo en el modelo de backend
		Pluck("RUN", &runsProfesorDB).Error; err != nil {
		return 0, 0, fmt.Errorf("error al obtener RUNs de ProfesorDB: %w", err)
	}

	// Paso 2: Filtrar Pipelsoft por RUNs coincidentes, unidad mayor, unidad menor y semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB).               // Filtrar por RUNs coincidentes
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}). // Filtrar por estados
		Where("nombre_unidad_mayor = ?", unidadMayor).            // Filtrar por unidad mayor
		Where("nombre_unidad_menor = ?", unidadMenor).            // Filtrar por unidad menor
		Where("semestre = ?", semestre).                          // Filtrar por semestre
		Count(&count).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar registros filtrados de Pipelsoft: %w", err)
	}

	// Paso 3: Contar los RUNs únicos coincidentes entre ProfesorDB y Pipelsoft
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runsProfesorDB). // Filtrar por RUNs coincidentes
		Distinct("run_empleado").                   // Contar RUNs únicos
		Count(&totalRUNs).Error; err != nil {
		return 0, 0, fmt.Errorf("error al contar RUNs únicos en Pipelsoft: %w", err)
	}

	return count, totalRUNs, nil
}

type UnidadMenorConProfesores struct {
	NombreUnidadMenor string `json:"nombre_unidad_menor"`
	TotalProfesores   int    `json:"total_profesores"`
}

// 8.1
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor(unidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que unidadMayor no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro unidadMayor no puede estar vacío")
	}

	// Validar que semestre no esté vacío
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro semestre no puede estar vacío")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato para la unidad mayor proporcionada
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Where("unidad_mayor = ?", unidadMayor).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Filtrar los RUNs que estén presentes en la tabla Pipelsoft, considerando el semestre
	var runFiltrados []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runDocentes).
		Where("semestre = ?", semestre). // Filtrar por semestre en Pipelsoft
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

// 8.3
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayor(unidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Validar que el parámetro 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Validar que el semestre no esté vacío
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro 'semestre' debe ser proporcionado")
	}

	// Paso 1: Obtener los RUNs únicos de la tabla Contrato
	var runsContrato []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runsContrato).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar registros en Pipelsoft donde CodEstado no sea "F1", "F9" o "A9" y los RUNs obtenidos
	// También se filtra por 'unidadMayor' y 'semestre'
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).                 // RUNs coincidentes
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}). // CodEstado que no sean F1, F9 o A9 (inactivos)
		Where("nombre_unidad_mayor = ?", unidadMayor).            // Filtrar por 'unidadMayor'
		Where("semestre = ?", semestre).                          // Filtrar por semestre en Pipelsoft
		Group("nombre_unidad_mayor").                             // Agrupar por unidad mayor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores con PAF inactivos filtrados por unidad mayor: %w", err)
	}

	// Convertir resultados a un mapa
	unidadesConPAFInactivos := make(map[string]int)
	for _, resultado := range resultados {
		unidadesConPAFInactivos[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesConPAFInactivos, nil
}

// 8.2
// 8.2
func (s *EstadisticasService) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_8_3(unidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		UnidadMenor     string
		TotalProfesores int
	}

	// Validar que el parámetro 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato para el semestre específico
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft para el semestre específico
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("semestre = ?", semestre).
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Contar RUNs únicos por unidad menor en Contrato, filtrando por 'unidadMayor' y 'semestre'
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_menor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Where("unidad_mayor = ?", unidadMayor). // Filtro por 'unidadMayor'
		Group("unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores sin profesores en Pipelsoft: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesSinProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesSinProfesores[resultado.UnidadMenor] = resultado.TotalProfesores
	}

	// Devolver el mapa con los resultados
	return unidadesSinProfesores, nil
}

// 8.4
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivos(unidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que el parámetro 'unidadMayor' esté presente
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}

	// Validar que el parámetro 'semestre' esté presente
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro 'semestre' debe ser proporcionado")
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
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).             // RUNs coincidentes
		Where("cod_estado IN ?", []string{"F1", "F9", "A9"}). // CodEstado filtrado (activos)
		Where("nombre_unidad_mayor = ?", unidadMayor).        // Filtro por 'unidadMayor'
		Where("semestre = ?", semestre).                      // Filtro por 'semestre'
		Group("nombre_unidad_menor").                         // Agrupar por unidad menor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores con PAF activos: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenores, nil
}

// 8.5
func (s *EstadisticasService) ObtenerUnidadesMenoresPorCodEstadoPAF(codEstadoPAF, unidadMayor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que codEstadoPAF no esté vacío
	if codEstadoPAF == "" {
		return nil, fmt.Errorf("el parámetro codEstadoPAF no puede estar vacío")
	}

	// Validar que unidadMayor no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro unidadMayor no puede estar vacío")
	}

	// Validar que semestre no esté vacío
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro semestre no puede estar vacío")
	}

	// Consulta a la base de datos filtrando por codEstadoPAF, unidadMayor y semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("des_estado = ?", codEstadoPAF).
		Where("nombre_unidad_mayor = ?", unidadMayor).
		Where("semestre = ?", semestre). // Filtro por semestre
		Group("nombre_unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores por codEstadoPAF, unidadMayor y semestre: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenores, nil
}

// 9.1
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresPorUnidadMayor9_1(unidadMayor, unidadMenor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que unidadMayor no esté vacío
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro unidadMayor no puede estar vacío")
	}

	// Validar que semestre no esté vacío
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro semestre no puede estar vacío")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato para la unidad mayor proporcionada
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Where("unidad_mayor = ?", unidadMayor).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Filtrar los RUNs que estén presentes en la tabla Pipelsoft con el filtro de semestre
	var runFiltrados []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("run_empleado IN ?", runDocentes).
		Where("semestre = ?", semestre). // Filtro adicional por semestre solo en Pipelsoft
		Distinct("run_empleado").
		Pluck("run_empleado", &runFiltrados).Error; err != nil {
		return nil, fmt.Errorf("error al filtrar RUNs presentes en Pipelsoft: %w", err)
	}

	// Paso 3: Consultar las unidades menores y contar RUNs únicos de los filtrados en Pipelsoft
	query := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runFiltrados).
		Group("nombre_unidad_menor")

	// Si se proporciona el parámetro unidadMenor, agregarlo a la consulta
	if unidadMenor != "" {
		query = query.Where("nombre_unidad_menor = ?", unidadMenor)
	}

	// Ejecutar la consulta
	if err := query.Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores con profesores: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenoresConProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenoresConProfesores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenoresConProfesores, nil
}

// 9.2
func (s *EstadisticasService) ObtenerUnidadesMenoresSinProfesoresEnPipelsoft_9_2(unidadMayor, unidadMenor, semestre string) (map[string]int, error) {
	var resultados []struct {
		UnidadMenor     string
		TotalProfesores int
	}

	// Validar que los parámetros 'unidadMayor' y 'unidadMenor' estén presentes
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}
	if unidadMenor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMenor' debe ser proporcionado")
	}

	// Validar que el semestre no esté vacío
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro semestre debe ser proporcionado")
	}

	// Paso 1: Obtener todos los RUNs únicos de la tabla Contrato
	var runDocentes []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runDocentes).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Contrato: %w", err)
	}

	// Paso 2: Obtener RUNs únicos de la tabla Pipelsoft con el filtro por semestre
	var runPipelsoft []string
	if err := s.DB.Model(&models.Pipelsoft{}).
		Where("semestre = ?", semestre). // Filtro por semestre solo en Pipelsoft
		Distinct("run_empleado").
		Pluck("run_empleado", &runPipelsoft).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de la tabla Pipelsoft: %w", err)
	}

	// Paso 3: Filtrar RUNs que están en Contrato pero no en Pipelsoft
	runsUnicosSinPipelsoft := diferenciaDeSlices(runDocentes, runPipelsoft)

	// Paso 4: Consultar unidades menores sin profesores en Pipelsoft
	if err := s.DB.Model(&models.Contrato{}).
		Select("unidad_menor, COUNT(DISTINCT run_docente) as total_profesores").
		Where("run_docente IN ?", runsUnicosSinPipelsoft).
		Where("unidad_mayor = ?", unidadMayor). // Filtro por 'unidadMayor'
		Where("unidad_menor = ?", unidadMenor). // Filtro por 'unidadMenor'
		Group("unidad_menor").
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores sin profesores en Pipelsoft: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesSinProfesores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesSinProfesores[resultado.UnidadMenor] = resultado.TotalProfesores
	}

	return unidadesSinProfesores, nil
}

// 9.3
func (s *EstadisticasService) ObtenerUnidadesMayoresConProfesoresFiltradosPAFActivasPorUnidadMayorYUnidadMenor9_3(unidadMayor, unidadMenor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMayor string
		TotalProfesores   int
	}

	// Validar que los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' estén presentes
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}
	if unidadMenor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMenor' debe ser proporcionado")
	}
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro 'semestre' debe ser proporcionado")
	}

	// Paso 1: Obtener los RUNs únicos de la tabla Contrato
	var runsContrato []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runsContrato).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar registros en Pipelsoft donde CodEstado no sea "F1", "F9" o "A9" y los RUNs obtenidos
	// También se filtra por 'unidadMayor', 'unidadMenor' y 'semestre'
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_mayor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).                 // RUNs coincidentes
		Where("cod_estado NOT IN ?", []string{"F1", "F9", "A9"}). // CodEstado que no sean F1, F9 o A9 (inactivos)
		Where("nombre_unidad_mayor = ?", unidadMayor).            // Filtrar por 'unidadMayor'
		Where("nombre_unidad_menor = ?", unidadMenor).            // Filtrar por 'unidadMenor'
		Where("semestre = ?", semestre).                          // Filtrar por semestre
		Group("nombre_unidad_mayor").                             // Agrupar por unidad mayor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades mayores con PAF inactivos filtrados por unidad mayor, unidad menor y semestre: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesConPAFInactivos := make(map[string]int)
	for _, resultado := range resultados {
		unidadesConPAFInactivos[resultado.NombreUnidadMayor] = resultado.TotalProfesores
	}

	return unidadesConPAFInactivos, nil
}

// 9.4
func (s *EstadisticasService) ObtenerUnidadesMenoresConProfesoresFiltradosPAFActivosPorUnidadMayorYUnidadMenor9_4(unidadMayor, unidadMenor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que los parámetros 'unidadMayor', 'unidadMenor' y 'semestre' estén presentes
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' debe ser proporcionado")
	}
	if unidadMenor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMenor' debe ser proporcionado")
	}
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro 'semestre' debe ser proporcionado")
	}

	// Paso 1: Obtener los RUNs únicos de la tabla Contrato
	var runsContrato []string
	if err := s.DB.Model(&models.Contrato{}).
		Distinct("run_docente").
		Pluck("run_docente", &runsContrato).Error; err != nil {
		return nil, fmt.Errorf("error al obtener RUNs únicos de Contrato: %w", err)
	}

	// Paso 2: Filtrar registros en Pipelsoft con CodEstado válido y los RUNs obtenidos
	// También se filtra por 'unidadMayor', 'unidadMenor' y 'semestre'
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("run_empleado IN ?", runsContrato).             // RUNs coincidentes
		Where("cod_estado IN ?", []string{"F1", "F9", "A9"}). // CodEstado filtrado (activos)
		Where("nombre_unidad_mayor = ?", unidadMayor).        // Filtro por 'unidadMayor'
		Where("nombre_unidad_menor = ?", unidadMenor).        // Filtro por 'unidadMenor'
		Where("semestre = ?", semestre).                      // Filtro por 'semestre'
		Group("nombre_unidad_menor").                         // Agrupar por unidad menor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores con PAF activos filtrados por unidad mayor, unidad menor y semestre: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenores, nil
}

// 9.5
func (s *EstadisticasService) ObtenerUnidadesMenoresPorCodEstadoPAFPorCodEstadoYUnidadMayorYUnidadMenor9_5(codEstadoPAF, unidadMayor, unidadMenor, semestre string) (map[string]int, error) {
	var resultados []struct {
		NombreUnidadMenor string
		TotalProfesores   int
	}

	// Validar que los parámetros 'codEstadoPAF', 'unidadMayor', 'unidadMenor' y 'semestre' no estén vacíos
	if codEstadoPAF == "" {
		return nil, fmt.Errorf("el parámetro 'codEstadoPAF' no puede estar vacío")
	}
	if unidadMayor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMayor' no puede estar vacío")
	}
	if unidadMenor == "" {
		return nil, fmt.Errorf("el parámetro 'unidadMenor' no puede estar vacío")
	}
	if semestre == "" {
		return nil, fmt.Errorf("el parámetro 'semestre' no puede estar vacío")
	}

	// Consulta a la base de datos filtrando por codEstadoPAF, unidadMayor, unidadMenor y semestre
	if err := s.DB.Model(&models.Pipelsoft{}).
		Select("nombre_unidad_menor, COUNT(DISTINCT run_empleado) as total_profesores").
		Where("des_estado = ?", codEstadoPAF).         // Filtro por 'codEstadoPAF'
		Where("nombre_unidad_mayor = ?", unidadMayor). // Filtro por 'unidadMayor'
		Where("nombre_unidad_menor = ?", unidadMenor). // Filtro adicional por 'unidadMenor'
		Where("semestre = ?", semestre).               // Filtro por 'semestre'
		Group("nombre_unidad_menor").                  // Agrupar por unidad menor
		Scan(&resultados).Error; err != nil {
		return nil, fmt.Errorf("error al obtener unidades menores por codEstadoPAF, unidadMayor, unidadMenor y semestre: %w", err)
	}

	// Convertir los resultados a un mapa
	unidadesMenores := make(map[string]int)
	for _, resultado := range resultados {
		unidadesMenores[resultado.NombreUnidadMenor] = resultado.TotalProfesores
	}

	return unidadesMenores, nil
}
