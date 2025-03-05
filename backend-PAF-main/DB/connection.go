package DB

import (
    "fmt"
    "log"
    "os"

    "github.com/NicolasAMunozR/PAF/backend-PAF/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    DBPersonal  *gorm.DB
    DBPersonal1 *gorm.DB
)

func InitDBConnections() {
    var err error

    // Obtén las credenciales de la base de datos desde las variables de entorno
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbPort := os.Getenv("DB_PORT")
    dbNamePersonal := os.Getenv("DB_NAME_PERSONAL")
    dbNamePersonal1 := os.Getenv("DB_NAME_PERSONAL1")

    // DSN para la base de datos personal
    DSNPersonal := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPassword, dbNamePersonal, dbPort)

    // DSN para la base de datos personal1
    DSNPersonal1 := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPassword, dbNamePersonal1, dbPort)

    // Conectar a la base de datos personal
    DBPersonal, err = gorm.Open(postgres.Open(DSNPersonal), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error al conectar a la base de datos personal: %v", err)
    }

    // Conectar a la base de datos personal1
    DBPersonal1, err = gorm.Open(postgres.Open(DSNPersonal1), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error al conectar a la base de datos personal1: %v", err)
    }

    // Migrar las tablas en ambas bases de datos
    err = DBPersonal.AutoMigrate(&models.HistorialPafAceptadas{}, &models.Horario{}, &models.Pipelsoft{}, &models.ProfesorDB{}, &models.Contrato{}, &models.Usuarios{}, &models.HistorialPasosPaf{}, &models.Auditoria{}, &models.Archivo{}, &models.ArchivoAdjunto{})
    if err != nil {
        log.Fatalf("Error al migrar la base de datos personal: %v", err)
    }

    err = DBPersonal1.AutoMigrate(&models.HistorialPafAceptadas{}, &models.Horario{}, &models.Pipelsoft{}, &models.ProfesorDB{}, &models.Contrato{}, &models.Usuarios{}, &models.HistorialPasosPaf{}, &models.Auditoria{}, &models.Archivo{}, &models.ArchivoAdjunto{})
    if err != nil {
        log.Fatalf("Error al migrar la base de datos personal1: %v", err)
    }

    fmt.Println("Conexión a las bases de datos exitosa.")
}