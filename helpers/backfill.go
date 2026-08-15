package helpers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"log"
)

func BackfillRegisteredBy() {
	var admin models.User
	if err := config.DB.Where("role = ?", models.RoleAdmin).First(&admin).Error; err != nil {
		log.Println("Tidak menemukan admin, melewati backfill")
		return
	}

	result := config.DB.Model(&models.User{}).
		Where("role = ? AND registered_by_id IS NULL", models.RoleCustomer).
		Update("registered_by_id", admin.ID)
	
	if result.Error != nil {
		log.Println("Error backfilling registered_by_id:", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Berhasil mengisi registered_by_id untuk %d pelanggan\n", result.RowsAffected)
	}
}
