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
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db := config.ConnectDatabase()

	// Auto Migrate Schema
	if err := db.AutoMigrate(&models.User{}, &models.IPLockout{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	log.Println("Migrasi database berhasil")

	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		passwordStr := string(hashedPassword)
		admin := models.User{
			Name:     "Admin Utama",
			Phone:    "081234567890",
			Role:     models.RoleAdmin,
			Password: &passwordStr,
		}
		db.Create(&admin)
		log.Println("Akun Admin default berhasil dibuat (Phone: 081234567890)")
	}

	r := gin.Default()

	// Enable CORS untuk semua origin dan izinkan header Authorization
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

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
