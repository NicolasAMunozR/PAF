package DB

import (
	"fmt"
	"log"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBPersonal *gorm.DB
)
var DSN = "host=localhost user=postgres password=alonsoreyes104 dbname=PAF port=5432 sslmode=disable"
var DB *gorm.DB

// Conexiones para múltiples bases de datos
var (
	DSNTPersonal = "host=localhost user=postgres password=alonsoreyes104 dbname=personal1 port=5432 sslmode=disable"
)

// InitDBConnections inicializa las conexiones a todas las bases de datos.
func InitDBConnections() {
	var err error

	// Conectar a la base de datos personal
	DBPersonal, err = gorm.Open(postgres.Open(DSNTPersonal), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos personal: %v", err)
	}

	err = DBPersonal.AutoMigrate(models.HistorialPafAceptadas{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}

	err = DBPersonal.AutoMigrate(models.Horario{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.Pipelsoft{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.ProfesorDB{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}

	err = DBPersonal.AutoMigrate(models.Contrato{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}

	err = DBPersonal.AutoMigrate(models.Usuarios{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.HistorialPasosPaf{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.Auditoria{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.Archivo{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	err = DBPersonal.AutoMigrate(models.ArchivoAdjunto{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}

	fmt.Println("Conexión a la base de datos TerceraDB exitosa.")
}
