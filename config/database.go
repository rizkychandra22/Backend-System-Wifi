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

func ConnectDatabase() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment variable dari sistem")
	}

	host     := os.Getenv("DB_HOST")
	user     := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname   := os.Getenv("DB_DATABASE")
	port     := os.Getenv("DB_PORT")
	sslmode  := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Gagal terkoneksi ke database PostgreSQL:", err)
	}

	log.Println("Berhasil terkoneksi ke database PostgreSQL!")
	DB = database
	return database
}
