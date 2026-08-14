package seeder

import (
	"backend-wifi/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&models.User{}).Where("role = ?", "admin").Count(&count)
	if count == 0 {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		passwordStr := string(hashedPassword)
		admin := models.User{
			Name:     "Admin NetVerse",
			Phone:    os.Getenv("ADMIN_PHONE"),
			Role:     models.RoleAdmin,
			Password: &passwordStr,
		}
		db.Create(&admin)
		log.Println("Akun Admin default berhasil dibuat oleh sistem")
	}
}
