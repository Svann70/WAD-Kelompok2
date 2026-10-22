package config

import (
	"fmt"
	"log"
	"os"

	"jobtrack-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca environment variables sistem")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	var dsn string

	if databaseURL != "" {
		dsn = databaseURL
	} else {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbname := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		sslmode := os.Getenv("DB_SSLMODE")

		if sslmode == "" {
			if host != "" && host != "localhost" && host != "127.0.0.1" {
				sslmode = "require"
			} else {
				sslmode = "disable"
			}
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode)
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal Koneksi ke Database: %v\n", err)
	}

		DB = database
		log.Println("Berhasil Terhubung ke Db")

		database.AutoMigrate(&models.Job{}, &models.StatusHistory{})
}
	DB = database
	log.Println("Berhasil Terhubung ke Database")
}
