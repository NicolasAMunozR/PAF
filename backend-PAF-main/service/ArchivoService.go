package service

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

// EstadisticasService define los métodos para obtener estadísticas de las tablas
type ArchivoService struct {
	DB *gorm.DB
}

// NewContratoService crea un nuevo servicio de contrato.
func NewArchivoService(db *gorm.DB) *ArchivoService {
	return &ArchivoService{
		DB: db,
	}
}

func CrearPDF(db *gorm.DB, Run string) error {
	// Obtener datos del empleado
	var pipelsoft models.Pipelsoft
	if err := db.Where("run_empleado = ?", Run).First(&pipelsoft).Error; err != nil {
		return err
	}

	// Crear instancia de Archivo
	archivo := models.Archivo{
		UnidadMayor:         pipelsoft.NombreUnidadMayor,
		UnidadMenor:         pipelsoft.NombreUnidadMenor,
		CelulaIdentidad:     Run,
		ApellidoP:           pipelsoft.PrimerApp,
		ApellidoM:           pipelsoft.SegundoApp,
		Nombres:             pipelsoft.Nombres,
		FechaInicioContrato: pipelsoft.FechaInicioContrato,
		FechaFinContrato:    pipelsoft.FechaFinContrato,
	}

	// Crear PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Agregar logo
	logoPath := "assets/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, 160, 10, 30, 0, false, "", 0, "")
	}

	// ------------------------**Título Principal**--------------------------
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0) // Azul oscuro

	title := "PROPUESTA Y ASUNCION DE FUNCIONES"
	titleWidth := pdf.GetStringWidth(title)
	pageWidth := 210.0 // Ancho de la página A4 en mm

	pdf.SetX((pageWidth - titleWidth) / 2) // Centra el título
	pdf.CellFormat(titleWidth, 10, title, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	// ------------------------------------**Subtítulo (Ley)**-----------------------------------
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0) // Azul oscuro

	subtitle := "(Ley Numero 17.654, articulo 38)"
	subtitleWidth := pdf.GetStringWidth(subtitle)

	pdf.SetX((pageWidth - subtitleWidth) / 2) // Centra el subtítulo
	pdf.CellFormat(subtitleWidth, 10, subtitle, "", 0, "C", false, 0, "")
	pdf.Ln(9)

	pdf.SetFont("Arial", "B", 9) // setear fuente y tamaño de letra
	// Filas de datos
	pdf.SetFillColor(255, 255, 255)
	//--------------------------------- Primeras 3 celdas---------------------
	pdf.CellFormat(50, 6, "UnidadMayor", "1", 0, "L", true, 0, "")
	// Primera celda con fondo relleno
	pdf.CellFormat(130, 6, archivo.UnidadMayor, "1", 0, "L", true, 0, "")

	// Guardamos la posición X inicial antes de escribir el segundo dato
	x := pdf.GetX()

	// Retrocedemos X para escribir el número a la derecha dentro de la misma celda
	pdf.SetX(x - 20) // Movemos el cursor 20 mm hacia la izquierda (ajustar según el espacio necesario)
	pdf.CellFormat(20, 6, "50", "0", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(50, 6, "Numero Centro de Costos", "1", 0, "L", true, 0, "")
	pdf.CellFormat(130, 6, "SOY UN PLACEHOLDER", "1", 0, "L", true, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(50, 6, "UnidadMenor", "1", 0, "L", true, 0, "")
	pdf.CellFormat(130, 6, archivo.UnidadMenor, "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	pdf.CellFormat(30, 6, "Celula identidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, Run, "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "Fecha y lugar de nacimiento", "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 6, "SOY UN PLACEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 6, "P 01/12/2024", "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	// ---------------------------------------------------------------------------
	pdf.CellFormat(60, 6, archivo.ApellidoP, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, archivo.ApellidoM, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, archivo.Nombres, "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(60, 6, "Apellido paterno", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Apellido materno", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Nombres", "1", 0, "L", true, 0, "")

	//-------------------------------------------------------------------------
	pdf.Ln(9)
	pdf.CellFormat(40, 6, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 6, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(40, 6, "Nacionalidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 6, "Domicilio", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Correo", "1", 0, "L", true, 0, "")
	pdf.Ln(9)
	//---------------------------------------------------------------------------
	pdf.CellFormat(70, 6, "Titulo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "Institucion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "Fecha obtencion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "Nro semestres", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(70, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.Ln(9)
	//-------------------------------------------------------------------------
	pdf.CellFormat(70, 7, "Grado academico", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 7, "Institucion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Fecha obtencion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Nro semestres", "1", 0, "L", true, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(70, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	//-------------------------------------------------------
	pdf.CellFormat(180, 7, "Tipo de ingreso", "1", 0, "C", true, 0, "")
	pdf.Ln(6)
	// Guarda la posición inicial de la fila
	x2 := pdf.GetX() //esto se utiliza para poder modificar la posicion de los elementos, y tener mayor control
	y2 := pdf.GetY()

	pdf.MultiCell(30, 5, "Ingresa a la\nadministracion\npublica", "1", "C", true)

	// Posiciona el cursor manualmente después de cada celda
	pdf.SetXY(x2+30, y2)
	pdf.MultiCell(30, 5, "Pertenece a la\nadministracion\npublica", "1", "C", true)

	pdf.SetXY(x2+60, y2)
	pdf.MultiCell(30, 5, "Se reincorpora\ncon menos de\n6 meses", "1", "C", true)

	pdf.SetXY(x2+90, y2)
	pdf.MultiCell(30, 5, "Se reincorpora\ncon mas de\n6 meses", "1", "C", true)

	pdf.SetXY(x2+120, y2)
	pdf.MultiCell(30, 5, "Ya pertenece\n \na la universidad", "1", "C", true)

	pdf.SetXY(x2+150, y2)
	pdf.MultiCell(30, 5, "Ingreso a la\n \nUSACH", "1", "C", true)

	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	//-----------------------------------------------------------------------------

	//------------------------------IDENTIFICACION DEL CARGO------------------------
	pdf.Ln(9)
	pdf.CellFormat(180, 5, "Identificacion del cargo:", "1", 0, "C", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(42, 4, "Cargo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Nivel", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Grado", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Rango", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Calidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Funcion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Jerarquia", "1", 0, "L", true, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(42, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(100, 4, "Asignatura", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 4, "N Horas", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 4, "Categoria", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 4, "Calidad", "1", 0, "L", true, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(100, 5, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	//----------------------------------------------------------------------------

	//--------------------------------CARGO O ACTIVIDAD QUE DESEMPEÑAN-------------------------------
	pdf.CellFormat(180, 6, "Cargo o actividad que seguira desempeñando en el otro servicio publico:", "1", 0, "C", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(41, 5, "Luegar de desempeño", "1", 0, "L", true, 0, "")
	pdf.CellFormat(41, 5, "Cargo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 5, "Grado", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 5, "Nivel", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "Rango", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "N Horas", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "Calidad", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(41, 7, "Plceholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(41, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, "Placeholder", "1", 0, "L", true, 0, "")
	pdf.Ln(8)

	//-------------------------------------------------------------------------------------------------------------------

	//-------------------------TEXTO PLANO SIN CASILLAS (LAS CASILLAS ESTAN PERO NO SE PUEDEN VER)--------------------------

	// **Datos Personales**
	pdf.SetFont("Arial", "", 11)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 6, "Mediante el presente acto tomo conocimiento de la politica integral para el abordaje de la violencia de genero en la", "0", 0, "L", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "Universidad de Santiago de Chile la cual fue aprobada mediante Resolucion 9011 de fecha 10 de noviembre", "0", 0, "L", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "de 2023 y publicada en https://direcciondegenero.usach.cl/normativa. segun la ley numero 21.369", "0", 0, "L", true, 0, "")
	// Espacio antes del bloque de fechas
	pdf.Ln(7)

	//-----------------------------------------------------------------------------------------------------------------

	// Texto "Asumió funciones a contar"
	pdf.CellFormat(0, 6, "Asumio funciones a contar ", "", 0, "L", false, 0, "")
	pdf.Ln(5)

	// Guarda la posición inicial para alinear correctamente
	x3 := pdf.GetX()
	y := pdf.GetY()
	pdf.SetXY(x3+50, y-5)
	x4 := pdf.GetX()
	y2 = pdf.GetY()

	// Primera fecha (08 - 09 - 2024)
	pdf.CellFormat(12, 7, "08", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "09", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "2024", "1", 0, "C", false, 0, "")

	// Mueve la posición a la línea de abajo
	pdf.SetXY(x4, y2+7)
	pdf.CellFormat(12, 7, "Dia", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "Mes", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "Año", "1", 0, "C", false, 0, "")

	// Mueve la posición para el "Hasta"
	pdf.SetXY(x4+50, y2)
	pdf.CellFormat(0, 6, "Hasta", "", 0, "L", false, 0, "")

	// Segunda fecha (30 - 12 - 2024)
	pdf.SetXY(x4+70, y2)
	pdf.CellFormat(12, 7, "30", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "12", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "2024", "1", 0, "C", false, 0, "")

	// Mueve la posición a la línea de abajo
	pdf.SetXY(x4+70, y2+7)
	pdf.CellFormat(12, 7, "Dia", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "Mes", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "Año", "1", 0, "C", false, 0, "")
	pdf.Ln(9)
	pdf.CellFormat(0, 6, "O HASTA QUE SUS SERVICIOS SEAN NECESARIOS", "0", 0, "C", true, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 7)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 6, "Si al momento del cierre del periodo de instripcion de asignaturas no diera cumplimiento con la normativa que establece criterios para la planeacion docente de pregrado", "0", 0, "l", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "(Res. 8938/2023), la Universidad de Santiago de Chle podra poner termino al nombramiento resultante de la presente propuesta de asuncion de funciones (PAF)", "0", 0, "l", true, 0, "")
	pdf.Ln(8)
	pdf.CellFormat(90, 6, "Firmas:", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(90, 6, "Interesado(a)", "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 6, "Directivo(a) Superor:", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(90, 6, "Fecha:___/___/_____/", "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 6, "Fecha:___/___/_____/", "1", 0, "L", true, 0, "")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 4, "Adjuntar certificado de antiguedad del servicio publico donde se desempeño.", "0", 0, "l", true, 0, "")

	//		pdf.Ln(10)

	//	pdf.SetFont("Arial", "", 11)
	//	pdf.Cell(0, 8, fmt.Sprintf("Mediante el presente acto tomo conocimiento de la politica integral para el abordaje de la violencia de genero en la \nUniversidad de Santiago de Chile la cual fue aprobada mediante Resolucion 9011 de fecha 10 de noviembre de 2023 y \n publicada en https://direcciondegenero.usach.cl/normativa. segun la ley numero 21.369"))

	//	pdf.Cell(0, 8, fmt.Sprintf("RUT: %s", archivo.CelulaIdentidad))
	//	pdf.Ln(6)
	//	pdf.Cell(0, 8, "Correo: nikolas.salinas16@gmail.com")
	//	pdf.Ln(12)

	// **Información del Contrato**
	//	pdf.SetFont("Arial", "B", 12)
	//	pdf.SetFillColor(220, 220, 220)
	//	pdf.CellFormat(0, 8, "Informacion del Contrato", "1", 0, "C", true, 0, "")
	//	pdf.Ln(10)
	/*

		pdf.SetFont("Arial", "", 11)

		// Encabezados de la tabla
		pdf.SetFillColor(180, 180, 180)
		pdf.CellFormat(60, 8, "Campo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(100, 8, "Detalle", "1", 0, "C", true, 0, "")
		pdf.Ln(8)

		// Filas de datos
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(60, 8, "Cargo", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, "Nombre del Cargo", "1", 0, "L", true, 0, "")
		pdf.Ln(8)

		pdf.CellFormat(60, 8, "Fecha de Inicio", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, archivo.FechaInicioContrato.Format("02-01-2006"), "1", 0, "L", true, 0, "")
		pdf.Ln(8)

		pdf.CellFormat(60, 8, "Fecha de Fin", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, archivo.FechaFinContrato.Format("02-01-2006"), "1", 0, "L", true, 0, "")
		pdf.Ln(12)

		// **Espacio para Firma**
		pdf.Ln(10)
		pdf.Cell(0, 8, "_________________________________")
		pdf.Ln(6)
		pdf.Cell(0, 8, "Firma del Responsable")
	*/

	// Guardar en buffer
	var pdfBuffer bytes.Buffer
	err := pdf.Output(&pdfBuffer)
	if err != nil {
		return err
	}

	// Guardar en la BD
	archivo.ArchivoPDF = pdfBuffer.Bytes()
	if err := db.Create(&archivo).Error; err != nil {
		return err
	}

	return nil
}

// Obtener archivo PDF por ID
func (s *ArchivoService) GetArchivoPDF(id uint) ([]byte, error) {
	var archivo models.Archivo
	if err := s.DB.First(&archivo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("archivo no encontrado")
		}
		return nil, err
	}
	return archivo.ArchivoPDF, nil
}

// Obtiene ruts desde el SAI y los compara con pipelsoft, si un rut presente en el sai esta presente en pipelsoft se revisa la planta, si no es academico se crea un contrato

func CreaRContratoAutomaticamentePorSemestre(db *gorm.DB, semestre string) ([]models.ProfesorDB, []models.ProfesorDB, error) {
	var profesores []models.ProfesorDB
	var rutsPipelsoft []string
	plantaMap := make(map[string]string) // Mapa para almacenar la planta de cada RUN en Pipelsoft

	// Obtener RUNs y Planta de la tabla Pipelsoft
	rows, err := db.Table("pipelsofts").Select("run_empleado, planta").Rows()
	if err != nil {
		return nil, nil, fmt.Errorf("error obteniendo RUNs de Pipelsoft: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var run, planta string
		if err := rows.Scan(&run, &planta); err != nil {
			return nil, nil, fmt.Errorf("error escaneando datos de Pipelsoft: %v", err)
		}
		rutsPipelsoft = append(rutsPipelsoft, run)
		plantaMap[run] = planta // Guardamos la planta asociada al RUN
	}

	// Obtener todos los profesores del semestre solicitado
	if err := db.Where("semestre = ?", semestre).Find(&profesores).Error; err != nil {
		return nil, nil, fmt.Errorf("error obteniendo profesores de ProfesorDB para el semestre %s: %v", semestre, err)
	}

	// Convertir listas en conjuntos (map) para búsqueda rápida
	setPipelsoft := make(map[string]bool)
	for _, rut := range rutsPipelsoft {
		setPipelsoft[rut] = true
	}

	var rutsNoComunes []models.ProfesorDB
	var rutsAcademicos []models.ProfesorDB

	for _, profesor := range profesores {
		if setPipelsoft[profesor.RUN] {
			// Revisar si pertenece a la planta "ACADEMICO"
			if plantaMap[profesor.RUN] == "ACADEMICO" {
				rutsAcademicos = append(rutsAcademicos, profesor) // Agregar a lista de académicos
			} else {
				CrearPDF(db, profesor.RUN) // Si no es "ACADEMICO", generar contrato
			}
		} else {
			rutsNoComunes = append(rutsNoComunes, profesor)
		}
	}

	return rutsNoComunes, rutsAcademicos, nil
}

// AgregarComentario agrega un comentario a un archivo identificado por su ID
func (s *ArchivoService) AgregarComentario(id uint, comentario string) error {
	var archivo models.Archivo

	// Buscar el archivo por ID
	if err := s.DB.First(&archivo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("archivo no encontrado")
		}
		return err
	}

	// Agregar el comentario
	archivo.Comentario = comentario

	// Guardar los cambios en la base de datos
	if err := s.DB.Save(&archivo).Error; err != nil {
		return err
	}

	return nil
}

// SubirArchivo permite almacenar un nuevo archivo adjunto a un Archivo existente
func (s *ArchivoService) SubirArchivo(archivoID uint, nombre string, datos []byte) error {
	// Verificar si el Archivo existe
	var archivo models.Archivo
	if err := s.DB.First(&archivo, archivoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("archivo no encontrado")
		}
		return err
	}

	// Crear y guardar el nuevo archivo adjunto
	archivoAdjunto := models.ArchivoAdjunto{
		ArchivoID: archivoID,
		Nombre:    nombre,
		Datos:     datos,
	}

	if err := s.DB.Create(&archivoAdjunto).Error; err != nil {
		return err
	}

	return nil
}

// ObtenerListaArchivos obtiene todos los archivos adjuntos de un Archivo específico
func (s *ArchivoService) ObtenerListaArchivos(archivoID uint) ([]models.ArchivoAdjunto, error) {
	var archivos []models.ArchivoAdjunto
	if err := s.DB.Where("archivo_id = ?", archivoID).Find(&archivos).Error; err != nil {
		return nil, err
	}
	return archivos, nil
}

// ModificarArchivo permite actualizar un archivo adjunto existente
func (s *ArchivoService) ModificarArchivo(archivoAdjuntoID uint, nuevoNombre string, nuevosDatos []byte) error {
	var archivoAdjunto models.ArchivoAdjunto
	if err := s.DB.First(&archivoAdjunto, archivoAdjuntoID).Error; err != nil {
		return errors.New("archivo adjunto no encontrado")
	}

	// Actualizar los datos
	archivoAdjunto.Nombre = nuevoNombre
	archivoAdjunto.Datos = nuevosDatos

	if err := s.DB.Save(&archivoAdjunto).Error; err != nil {
		return err
	}
	return nil
}

// GetArchivosAdjuntosByRut obtiene los archivos adjuntos asociados a un RUT específico
func (s *ArchivoService) GetArchivosAdjuntosByRut(rut string) ([]models.ArchivoAdjunto, error) {
	var archivosAdjuntos []models.ArchivoAdjunto
	var archivos []models.Archivo

	// Buscar los IDs de los archivos asociados al RUT
	err := s.DB.Where("celula_identidad = ?", rut).Find(&archivos).Error
	if err != nil {
		return nil, err
	}

	// Extraer los IDs de los archivos
	var archivoIDs []uint
	for _, archivo := range archivos {
		archivoIDs = append(archivoIDs, archivo.ID)
	}

	// Buscar los archivos adjuntos relacionados con los archivos encontrados
	if len(archivoIDs) > 0 {
		err = s.DB.Where("archivo_id IN ?", archivoIDs).Find(&archivosAdjuntos).Error
		if err != nil {
			return nil, err
		}
	}

	return archivosAdjuntos, nil
}

// GetArchivoAdjuntoByID obtiene un archivo adjunto por su ID
func (s *ArchivoService) GetArchivoAdjuntoByID(id uint) (*models.ArchivoAdjunto, error) {
	var archivoAdjunto models.ArchivoAdjunto

	// Buscar el archivo en la base de datos
	err := s.DB.First(&archivoAdjunto, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("archivo no encontrado")
	} else if err != nil {
		return nil, err
	}

	// Verificar si el archivo tiene datos
	if len(archivoAdjunto.Datos) == 0 {
		return nil, errors.New("el archivo está vacío")
	}

	return &archivoAdjunto, nil
}

func CrearPDFSinData(db *gorm.DB, UnidadMayor string, UnidadMenor string, NumeroCentroDeCostos string, CelulaIdentidad string, LugarNacimiento string, FechaYHoraNacimiento string, ApellidoP string, ApellidoM string, Nombres string, Nacionalidad string, Domicio string, Correo string, Titulo string, Institucion string, FechaObtencion string, NumeroSemestre string, GradoAcademico string, InstitucionGradoAcademico string, FechaObtencionGradoAcademico string, TipoIngreso string, Cargo string, Nivel string, Grado string, Rango string, Funcion string, Jerarquia string, Asignatura string, NumeroHoras string, Categoria string, Calidad string, LugarDesempeño string, CargoOtroPublico string, GradoOtroPublico string, NivelOtroPublico string, RangoOtroPublico string, NumeroHorasOtroPublico int, CalidadOtroPublico string, FechaInicioContrato string, FechaFinContrato string) error {
	// Obtener datos del empleado

	archivo := models.Archivo{
		UnidadMayor:                  UnidadMayor,
		UnidadMenor:                  UnidadMenor,
		NumeroCentroDeCostos:         NumeroCentroDeCostos,
		CelulaIdentidad:              CelulaIdentidad, // Si "Run" viene de otro lado, asegúrate de definirlo antes
		LugarNacimiento:              LugarNacimiento,
		FechaYHoraNacimiento:         parseFechaS(FechaYHoraNacimiento), // Se debe convertir string a time.Time
		ApellidoP:                    ApellidoP,
		ApellidoM:                    ApellidoM,
		Nombres:                      Nombres,
		Nacionalidad:                 Nacionalidad,
		Domicio:                      Domicio,
		Correo:                       Correo,
		Titulo:                       Titulo,
		Institucion:                  Institucion,
		FechaObtencion:               parseFechaS(FechaObtencion), // Conversión
		NumeroSemestre:               NumeroSemestre,
		GradoAcademico:               GradoAcademico,
		InstitucionGradoAcademico:    InstitucionGradoAcademico,
		FechaObtencionGradoAcademico: parseFechaS(FechaObtencionGradoAcademico), // Conversión
		TipoIngreso:                  TipoIngreso,
		Cargo:                        Cargo,
		Nivel:                        Nivel,
		Grado:                        Grado,
		Rango:                        Rango,
		Funcion:                      Funcion,
		Jerarquia:                    Jerarquia,
		Asignatura:                   Asignatura,
		NumeroHoras:                  NumeroHoras,
		Categoria:                    Categoria,
		Calidad:                      Calidad,
		LugarDesempeño:               LugarDesempeño,
		CargoOtroPublico:             CargoOtroPublico,
		GradoOtroPublico:             GradoOtroPublico,
		NivelOtroPublico:             NivelOtroPublico,
		RangoOtroPublico:             RangoOtroPublico,
		NumeroHorasOtroPublico:       NumeroHorasOtroPublico,
		CalidadOtroPublico:           CalidadOtroPublico,
		FechaInicioContrato:          parseFechaS(FechaInicioContrato), // Conversión
		FechaFinContrato:             parseFechaS(FechaFinContrato),    // Conversión
	}
	// Crear PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Agregar logo
	logoPath := "assets/logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.Image(logoPath, 160, 10, 30, 0, false, "", 0, "")
	}

	// ------------------------**Título Principal**--------------------------
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0) // Azul oscuro

	title := "PROPUESTA Y ASUNCION DE FUNCIONES"
	titleWidth := pdf.GetStringWidth(title)
	pageWidth := 210.0 // Ancho de la página A4 en mm

	pdf.SetX((pageWidth - titleWidth) / 2) // Centra el título
	pdf.CellFormat(titleWidth, 10, title, "", 0, "C", false, 0, "")
	pdf.Ln(5)

	// ------------------------------------**Subtítulo (Ley)**-----------------------------------
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(0, 0, 0) // Azul oscuro

	subtitle := "(Ley Numero 17.654, articulo 38)"
	subtitleWidth := pdf.GetStringWidth(subtitle)

	pdf.SetX((pageWidth - subtitleWidth) / 2) // Centra el subtítulo
	pdf.CellFormat(subtitleWidth, 10, subtitle, "", 0, "C", false, 0, "")
	pdf.Ln(9)

	pdf.SetFont("Arial", "B", 9) // setear fuente y tamaño de letra
	// Filas de datos
	pdf.SetFillColor(255, 255, 255)
	//--------------------------------- Primeras 3 celdas---------------------
	pdf.CellFormat(50, 6, "UnidadMayor", "1", 0, "L", true, 0, "")
	// Primera celda con fondo relleno
	pdf.CellFormat(130, 6, archivo.UnidadMayor, "1", 0, "L", true, 0, "")

	// Guardamos la posición X inicial antes de escribir el segundo dato
	x := pdf.GetX()

	// Retrocedemos X para escribir el número a la derecha dentro de la misma celda
	pdf.SetX(x - 20) // Movemos el cursor 20 mm hacia la izquierda (ajustar según el espacio necesario)
	pdf.CellFormat(20, 6, "50", "0", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(50, 6, "Numero Centro de Costos", "1", 0, "L", true, 0, "")
	pdf.CellFormat(130, 6, "SOY UN PLACEHOLDER", "1", 0, "L", true, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(50, 6, "UnidadMenor", "1", 0, "L", true, 0, "")
	pdf.CellFormat(130, 6, archivo.UnidadMenor, "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	pdf.CellFormat(30, 6, "Celula identidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, archivo.CelulaIdentidad, "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "Fecha y lugar de nacimiento", "1", 0, "L", true, 0, "")
	pdf.CellFormat(45, 6, archivo.FechaYHoraNacimiento.GoString(), "1", 0, "L", true, 0, "")
	pdf.CellFormat(25, 6, "P 01/12/2024", "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	// ---------------------------------------------------------------------------
	pdf.CellFormat(60, 6, archivo.ApellidoP, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, archivo.ApellidoM, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, archivo.Nombres, "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(60, 6, "Apellido paterno", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Apellido materno", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Nombres", "1", 0, "L", true, 0, "")

	//-------------------------------------------------------------------------
	pdf.Ln(9)
	pdf.CellFormat(40, 6, archivo.Nacionalidad, "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 6, archivo.Domicio, "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, archivo.Correo, "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(40, 6, "Nacionalidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(80, 6, "Domicilio", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 6, "Correo", "1", 0, "L", true, 0, "")
	pdf.Ln(9)
	//---------------------------------------------------------------------------
	pdf.CellFormat(70, 6, "Titulo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, "Institucion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "Fecha obtencion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, "Nro semestres", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(70, 6, archivo.Titulo, "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, archivo.Institucion, "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, archivo.FechaObtencion.GoString(), "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, archivo.NumeroSemestre, "1", 0, "L", true, 0, "")
	pdf.Ln(9)
	//-------------------------------------------------------------------------
	pdf.CellFormat(70, 7, "Grado academico", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 7, "Institucion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Fecha obtencion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 7, "Nro semestres", "1", 0, "L", true, 0, "")
	pdf.Ln(6)

	pdf.CellFormat(70, 6, archivo.Grado, "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 6, archivo.InstitucionGradoAcademico, "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, archivo.FechaObtencionGradoAcademico.GoString(), "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 6, archivo.NumeroSemestre, "1", 0, "L", true, 0, "")
	pdf.Ln(9)

	//-------------------------------------------------------
	pdf.CellFormat(180, 7, "Tipo de ingreso", "1", 0, "C", true, 0, "")
	pdf.Ln(6)
	// Guarda la posición inicial de la fila
	x2 := pdf.GetX() //esto se utiliza para poder modificar la posicion de los elementos, y tener mayor control
	y2 := pdf.GetY()

	pdf.MultiCell(30, 5, "Ingresa a la\nadministracion\npublica", "1", "C", true)

	// Posiciona el cursor manualmente después de cada celda
	pdf.SetXY(x2+30, y2)
	pdf.MultiCell(30, 5, "Pertenece a la\nadministracion\npublica", "1", "C", true)

	pdf.SetXY(x2+60, y2)
	pdf.MultiCell(30, 5, "Se reincorpora\ncon menos de\n6 meses", "1", "C", true)

	pdf.SetXY(x2+90, y2)
	pdf.MultiCell(30, 5, "Se reincorpora\ncon mas de\n6 meses", "1", "C", true)

	pdf.SetXY(x2+120, y2)
	pdf.MultiCell(30, 5, "Ya pertenece\n \na la universidad", "1", "C", true)

	pdf.SetXY(x2+150, y2)
	pdf.MultiCell(30, 5, "Ingreso a la\n \nUSACH", "1", "C", true)

	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, "PALCEHOLDER", "1", 0, "L", true, 0, "")
	//-----------------------------------------------------------------------------

	//------------------------------IDENTIFICACION DEL CARGO------------------------
	pdf.Ln(9)
	pdf.CellFormat(180, 5, "Identificacion del cargo:", "1", 0, "C", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(42, 4, "Cargo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Nivel", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Grado", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Rango", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Calidad", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Funcion", "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, "Jerarquia", "1", 0, "L", true, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(42, 4, archivo.Cargo, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Nivel, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Grado, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Rango, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Calidad, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Funcion, "1", 0, "L", true, 0, "")
	pdf.CellFormat(23, 4, archivo.Jerarquia, "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(100, 4, "Asignatura", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 4, "N Horas", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 4, "Categoria", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 4, "Calidad", "1", 0, "L", true, 0, "")
	pdf.Ln(4)
	pdf.CellFormat(100, 5, archivo.Asignatura, "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, archivo.NumeroHoras, "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, archivo.Categoria, "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 5, archivo.Calidad, "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	//----------------------------------------------------------------------------

	//--------------------------------CARGO O ACTIVIDAD QUE DESEMPEÑAN-------------------------------
	pdf.CellFormat(180, 6, "Cargo o actividad que seguira desempeñando en el otro servicio publico:", "1", 0, "C", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(41, 5, "Luegar de desempeño", "1", 0, "L", true, 0, "")
	pdf.CellFormat(41, 5, "Cargo", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 5, "Grado", "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 5, "Nivel", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "Rango", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "N Horas", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 5, "Calidad", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(41, 7, archivo.LugarDesempeño, "1", 0, "L", true, 0, "")
	pdf.CellFormat(41, 7, archivo.CargoOtroPublico, "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 7, archivo.GradoOtroPublico, "1", 0, "L", true, 0, "")
	pdf.CellFormat(19, 7, archivo.NivelOtroPublico, "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, archivo.RangoOtroPublico, "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, archivo.NumeroHoras, "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 7, archivo.CalidadOtroPublico, "1", 0, "L", true, 0, "")
	pdf.Ln(8)

	//-------------------------------------------------------------------------------------------------------------------

	//-------------------------TEXTO PLANO SIN CASILLAS (LAS CASILLAS ESTAN PERO NO SE PUEDEN VER)--------------------------

	// **Datos Personales**
	pdf.SetFont("Arial", "", 11)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 6, "Mediante el presente acto tomo conocimiento de la politica integral para el abordaje de la violencia de genero en la", "0", 0, "L", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "Universidad de Santiago de Chile la cual fue aprobada mediante Resolucion 9011 de fecha 10 de noviembre", "0", 0, "L", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "de 2023 y publicada en https://direcciondegenero.usach.cl/normativa. segun la ley numero 21.369", "0", 0, "L", true, 0, "")
	// Espacio antes del bloque de fechas
	pdf.Ln(7)

	//-----------------------------------------------------------------------------------------------------------------

	// Texto "Asumió funciones a contar"
	pdf.CellFormat(0, 6, "Asumio funciones a contar ", "", 0, "L", false, 0, "")
	pdf.Ln(5)

	// Guarda la posición inicial para alinear correctamente
	x3 := pdf.GetX()
	y := pdf.GetY()
	pdf.SetXY(x3+50, y-5)
	x4 := pdf.GetX()
	y2 = pdf.GetY()

	// Primera fecha (08 - 09 - 2024)
	pdf.CellFormat(12, 7, "08", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "09", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "2024", "1", 0, "C", false, 0, "")

	// Mueve la posición a la línea de abajo
	pdf.SetXY(x4, y2+7)
	pdf.CellFormat(12, 7, "Dia", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "Mes", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "Año", "1", 0, "C", false, 0, "")

	// Mueve la posición para el "Hasta"
	pdf.SetXY(x4+50, y2)
	pdf.CellFormat(0, 6, "Hasta", "", 0, "L", false, 0, "")

	// Segunda fecha (30 - 12 - 2024)
	pdf.SetXY(x4+70, y2)
	pdf.CellFormat(12, 7, "30", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "12", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "2024", "1", 0, "C", false, 0, "")

	// Mueve la posición a la línea de abajo
	pdf.SetXY(x4+70, y2+7)
	pdf.CellFormat(12, 7, "Dia", "1", 0, "C", false, 0, "")
	pdf.CellFormat(12, 7, "Mes", "1", 0, "C", false, 0, "")
	pdf.CellFormat(13, 7, "Año", "1", 0, "C", false, 0, "")
	pdf.Ln(9)
	pdf.CellFormat(0, 6, "O HASTA QUE SUS SERVICIOS SEAN NECESARIOS", "0", 0, "C", true, 0, "")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 7)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 6, "Si al momento del cierre del periodo de instripcion de asignaturas no diera cumplimiento con la normativa que establece criterios para la planeacion docente de pregrado", "0", 0, "l", true, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 6, "(Res. 8938/2023), la Universidad de Santiago de Chle podra poner termino al nombramiento resultante de la presente propuesta de asuncion de funciones (PAF)", "0", 0, "l", true, 0, "")
	pdf.Ln(8)
	pdf.CellFormat(90, 6, "Firmas:", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(90, 6, "Interesado(a)", "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 6, "Directivo(a) Superor:", "1", 0, "L", true, 0, "")
	pdf.Ln(6)
	pdf.CellFormat(90, 6, "Fecha:___/___/_____/", "1", 0, "L", true, 0, "")
	pdf.CellFormat(90, 6, "Fecha:___/___/_____/", "1", 0, "L", true, 0, "")
	pdf.Ln(7)
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)
	pdf.CellFormat(0, 4, "Adjuntar certificado de antiguedad del servicio publico donde se desempeño.", "0", 0, "l", true, 0, "")

	//		pdf.Ln(10)

	//	pdf.SetFont("Arial", "", 11)
	//	pdf.Cell(0, 8, fmt.Sprintf("Mediante el presente acto tomo conocimiento de la politica integral para el abordaje de la violencia de genero en la \nUniversidad de Santiago de Chile la cual fue aprobada mediante Resolucion 9011 de fecha 10 de noviembre de 2023 y \n publicada en https://direcciondegenero.usach.cl/normativa. segun la ley numero 21.369"))

	//	pdf.Cell(0, 8, fmt.Sprintf("RUT: %s", archivo.CelulaIdentidad))
	//	pdf.Ln(6)
	//	pdf.Cell(0, 8, "Correo: nikolas.salinas16@gmail.com")
	//	pdf.Ln(12)

	// **Información del Contrato**
	//	pdf.SetFont("Arial", "B", 12)
	//	pdf.SetFillColor(220, 220, 220)
	//	pdf.CellFormat(0, 8, "Informacion del Contrato", "1", 0, "C", true, 0, "")
	//	pdf.Ln(10)
	/*

		pdf.SetFont("Arial", "", 11)

		// Encabezados de la tabla
		pdf.SetFillColor(180, 180, 180)
		pdf.CellFormat(60, 8, "Campo", "1", 0, "C", true, 0, "")
		pdf.CellFormat(100, 8, "Detalle", "1", 0, "C", true, 0, "")
		pdf.Ln(8)

		// Filas de datos
		pdf.SetFillColor(255, 255, 255)
		pdf.CellFormat(60, 8, "Cargo", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, "Nombre del Cargo", "1", 0, "L", true, 0, "")
		pdf.Ln(8)

		pdf.CellFormat(60, 8, "Fecha de Inicio", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, archivo.FechaInicioContrato.Format("02-01-2006"), "1", 0, "L", true, 0, "")
		pdf.Ln(8)

		pdf.CellFormat(60, 8, "Fecha de Fin", "1", 0, "L", true, 0, "")
		pdf.CellFormat(100, 8, archivo.FechaFinContrato.Format("02-01-2006"), "1", 0, "L", true, 0, "")
		pdf.Ln(12)

		// **Espacio para Firma**
		pdf.Ln(10)
		pdf.Cell(0, 8, "_________________________________")
		pdf.Ln(6)
		pdf.Cell(0, 8, "Firma del Responsable")
	*/

	// Guardar en buffer
	var pdfBuffer bytes.Buffer
	err := pdf.Output(&pdfBuffer)
	if err != nil {
		return err
	}

	// Guardar en la BD
	archivo.ArchivoPDF = pdfBuffer.Bytes()
	if err := db.Create(&archivo).Error; err != nil {
		return err
	}

	return nil
}

func parseFecha(FechaObtencion string) {
	panic("unimplemented")
}

func parseFechaS(fechaStr string) time.Time {
	layout := "2006-01-02" // Ajusta según el formato esperado
	t, err := time.Parse(layout, fechaStr)
	if err != nil {
		fmt.Println("Error al parsear fecha:", err)
		return time.Time{} // Retorna un valor vacío si falla la conversión
	}
	return t
}

//ObtenerProfesoresQueNoSePuedeGenerarContrato
//FALTA FILTRAR POR PLANTA

func ObtenerProfesoresQueNoSePuedeGenerarContrato(db *gorm.DB) ([]models.ProfesorDB, []models.ProfesorDB, []models.ProfesorDB, error) {
	var profesores []models.ProfesorDB
	var rutsPipelsoft []string
	var rutsArchivo []string

	// Obtener RUNs de la tabla Pipelsoft
	rows, err := db.Table("pipelsofts").Select("run_empleado").Rows()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error obteniendo RUNs de Pipelsoft: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var run string
		if err := rows.Scan(&run); err != nil {
			return nil, nil, nil, fmt.Errorf("error escaneando datos de Pipelsoft: %v", err)
		}
		rutsPipelsoft = append(rutsPipelsoft, run)
	}

	// Obtener todos los profesores (sin filtro por unidad mayor)
	if err := db.Find(&profesores).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("error obteniendo profesores de ProfesorDB: %v", err)
	}

	// Obtener todos los RUTs de la tabla archivo
	archivoRows, err := db.Table("archivos").Select("celula_identidad").Rows()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error obteniendo RUTs de la tabla archivo: %v", err)
	}
	defer archivoRows.Close()

	for archivoRows.Next() {
		var rut string
		if err := archivoRows.Scan(&rut); err != nil {
			return nil, nil, nil, fmt.Errorf("error escaneando datos de la tabla archivo: %v", err)
		}
		rutsArchivo = append(rutsArchivo, rut)
	}

	// Convertir listas en conjuntos (map) para búsqueda rápida
	setPipelsoft := make(map[string]bool)
	for _, rut := range rutsPipelsoft {
		setPipelsoft[rut] = true
	}

	setArchivo := make(map[string]bool)
	for _, rut := range rutsArchivo {
		setArchivo[rut] = true
	}

	var rutsNoComunes []models.ProfesorDB
	var rutsContratables []models.ProfesorDB
	var rutsConContrato []models.ProfesorDB

	for _, profesor := range profesores {
		// Si el profesor ya tiene contrato, lo agregamos a rutsConContrato y no debe estar en otras listas
		if setArchivo[profesor.RUN] {
			rutsConContrato = append(rutsConContrato, profesor)
			continue
		}

		if setPipelsoft[profesor.RUN] {
			// Agregar a los que pueden generar contrato
			rutsContratables = append(rutsContratables, profesor)
		} else {
			rutsNoComunes = append(rutsNoComunes, profesor)
		}
	}

	// Retornar las listas
	return rutsNoComunes, rutsContratables, rutsConContrato, nil
}

// Función para verificar si un RUT está en la lista de rutsConContrato
func containsRut(rutsConContrato []models.ProfesorDB, rut string) bool {
	for _, profesor := range rutsConContrato {
		if profesor.RUN == rut {
			return true
		}
	}
	return false
}
