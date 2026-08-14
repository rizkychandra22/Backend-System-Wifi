package config

import (
	"log"
	"time"
)

func init() {
	// Setingan Jam Global Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Gagal memuat zona waktu Asia/Jakarta, menggunakan zona waktu lokal")
	} else {
		time.Local = loc
	}
}
