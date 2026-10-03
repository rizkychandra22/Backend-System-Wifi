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

	// Bersihkan duplikasi absensi lama jika ada sebelum membuat unique index
	cleanDuplicateAttendances := `
		DELETE FROM attendances a
		USING attendances b
		WHERE a.id > b.id
		  AND a.user_id = b.user_id
		  AND a.date = b.date;
	`
	if err := db.Exec(cleanDuplicateAttendances).Error; err != nil {
		log.Printf("Catatan: Pembersihan awal absensi (bisa dilewati jika tabel baru): %v", err)
	}

	// Auto Migrate Schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.DeviceLockout{},
		&models.Attendance{},
		&models.WifiPackage{},
		&models.Payment{},
		&models.Overtime{},
		&models.Subscription{},
		&models.Allowance{},
	); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	log.Println("Migrasi database berhasil")

	seeder.SeedAdminUser(db)
	helpers.BackfillRegisteredBy()

	// Bebaskan device_id untuk role selain employee (customer dan admin)
	db.Model(&models.User{}).Where("role != ?", models.RoleEmployee).Update("device_id", nil)
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
	routes.SetupSubscriptionRoutes(r)
	routes.SetupAllowanceRoutes(r)

	// Setup Scheduler for Attendance
	locWIB, _ := time.LoadLocation("Asia/Jakarta")
	c := cron.New(cron.WithLocation(locWIB))

	// Holiday check and auto checkout
	c.AddFunc("1 12 * * *", helpers.CheckAutoHoliday)
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
