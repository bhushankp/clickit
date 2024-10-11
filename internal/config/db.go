package config

import (
	"clickit/internal/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes the MySQL connection
func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_DSN")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	DB = database

	// Migrate the schema
	DB.AutoMigrate(&models.Record{})
}
