package seeder

import (
	"backend-wifi/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) {
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
}
