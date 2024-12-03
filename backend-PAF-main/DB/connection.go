package DB

import (
	"fmt"
	"log"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBPipelsoft *gorm.DB
	DBProfesor  *gorm.DB
	DBPersonal  *gorm.DB
)
var DSN = "host=localhost user=postgres password=alonsoreyes104 dbname=PAF port=5432 sslmode=disable"
var DB *gorm.DB

// Conexiones para múltiples bases de datos
var (
	DSNProfesor  = "host=localhost user=postgres password=alonsoreyes104 dbname=Profesor port=5432 sslmode=disable"
	DSNPipelsoft = "host=localhost user=postgres password=alonsoreyes104 dbname=Pipelsoft port=5432 sslmode=disable"
	DSNTPersonal = "host=localhost user=postgres password=alonsoreyes104 dbname=personal port=5432 sslmode=disable"
)

// InitDBConnections inicializa las conexiones a todas las bases de datos.
func InitDBConnections() {
	var err error

	// Conectar a la base de datos Pipelsoft
	DBPipelsoft, err = gorm.Open(postgres.Open(DSNPipelsoft), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos Pipelsoft: %v", err)
	}
	err = DBPipelsoft.AutoMigrate(models.Pipelsoft{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos Pipelsoft: %v", err)
	}

	fmt.Println("Conexión a la base de datos Pipelsoft exitosa.")

	// Conectar a la base de datos ProfesorDB
	DBProfesor, err = gorm.Open(postgres.Open(DSNProfesor), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos ProfesorDB: %v", err)
	}
	err = DBProfesor.AutoMigrate(models.ProfesorDB{})
	if err != nil {
		log.Fatalf("Error al migrar la base de datos ProfesorDB: %v", err)
	}
	fmt.Println("Conexión a la base de datos ProfesorDB exitosa.")

	// Conectar a la base de datos TerceraDB
	DBPersonal, err = gorm.Open(postgres.Open(DSNTPersonal), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos TerceraDB: %v", err)
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

	fmt.Println("Conexión a la base de datos TerceraDB exitosa.")
}
