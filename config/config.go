package config

import (
	"log"
	"os"	
	"invoicerator/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL is not set in the environment")
    }

    database, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Migrate the schema
    err = database.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    DB = database
}
