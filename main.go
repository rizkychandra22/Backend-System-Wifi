package main

import (
	"backend-wifi/routes"
	"backend-wifi/config"
	"backend-wifi/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectDatabase()

	// Auto Migrate Schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	log.Println("Migrasi database berhasil")

	r := gin.Default()

	// Enable CORS untuk semua origin (agar Vite/React di port 5173 bisa menembak API)
	r.Use(cors.Default())

	// Serve static files from uploads directory
	r.Static("/uploads", "./uploads")

	// Setup Routes
	routes.SetupAuthRoutes(r)
	routes.SetupUserRoutes(r)

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Backend Absensi & Pembayaran WiFi Berjalan",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port " + port + "...")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
