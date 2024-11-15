package DB

import (
	"fmt"
	"log"

	"github.com/NicolasAMunozR/PAF/backend-PAF/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = "host=localhost user=postgres password=alonsoreyes104 dbname=PAF port=5432 sslmode=disable"
var DB *gorm.DB

// DBconnection establece la conexión a la base de datos y la asigna a la variable global DB.
func DBconnection() {
	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	err = DB.AutoMigrate(models.Pipelsoft{})
	if err != nil {
		panic("failed to migrate the database")
	}
	err = DB.AutoMigrate(models.HistorialPafAceptadas{})
	if err != nil {
		panic("failed to migrate the database")
	}

	err = DB.AutoMigrate(models.Horario{})
	if err != nil {
		panic("failed to migrate the database")
	}

	fmt.Println("Conexión a la base de datos exitosa.")

}
