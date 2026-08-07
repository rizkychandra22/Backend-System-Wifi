package main

import (
	"backend-wifi/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Koneksi ke Database PostgreSQL
	db := config.ConnectDatabase()
	_ = db // Temporary fix to prevent unused variable error

	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Backend Absensi & Pembayaran WiFi Berjalan",
		})
	})

	log.Println("Server is running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
