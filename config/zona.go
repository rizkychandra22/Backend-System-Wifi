package config

import (
	"log"
	"time"
)

var ZonaWaktu *time.Location

func init() {
	// Setingan Jam Global Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Gagal memuat zona waktu Asia/Jakarta, menggunakan zona waktu lokal")
		ZonaWaktu = time.Local
	} else {
		time.Local = loc
		ZonaWaktu = loc
	}
}
