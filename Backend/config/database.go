package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan")
	}


host := os.Getenv("DB_HOST")
user := os.Getenv("DB_USER")
password := os.Getenv("DB_PASSWORD")
dbname := os.Getenv("DB_NAME")
port := os.Getenv("DB_PORT")
sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port,
  sslmode)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Gagal Koneksi ke Db: %v\n", err)
			return
		}

		DB = database
		log.Println("Berhasil Terhubung ke Db")
}