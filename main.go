package main

import (
	"backend-wifi/config"
	"backend-wifi/helpers"
	"backend-wifi/models"
	"backend-wifi/models/seeder"
	"backend-wifi/routes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	db := config.ConnectDatabase()

	// Auto Migrate Schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.DeviceLockout{},
		&models.Attendance{},
		&models.WifiPackage{},
		&models.Payment{},
		&models.Overtime{},
	); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	log.Println("Migrasi database berhasil")

	seeder.SeedAdminUser(db)
	helpers.BackfillRegisteredBy()
	r := gin.Default()

	// Enable CORS untuk semua origin dan izinkan header Authorization
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Setup Routes
	routes.SetupAuthRoutes(r)
	routes.SetupUserRoutes(r)
	routes.SetupAttendanceRoutes(r)
	routes.SetupCustomerRoutes(r)
	routes.SetupWifiPackageRoutes(r)
	routes.SetupPaymentRoutes(r)
	routes.SetupOvertimeRoutes(r)

	// Setup Scheduler for Attendance
	locWIB, _ := time.LoadLocation("Asia/Jakarta")
	c := cron.New(cron.WithLocation(locWIB))

	// Holiday check and auto checkout
	c.AddFunc("31 12 * * *", helpers.CheckAutoHoliday)
	c.AddFunc("1 17 * * *", helpers.AutoCheckout)

	c.Start()
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
