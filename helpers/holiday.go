package helpers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"log"
	"time"
)

func CheckAutoHoliday() {
	locWIB, _ := time.LoadLocation("Asia/Jakarta")
	dateStr := time.Now().In(locWIB).Format("2006-01-02")
	
	var count int64
	config.DB.Model(&models.Attendance{}).Where("date = ? AND status != ?", dateStr, models.StatusLibur).Count(&count)
	if count == 0 {
		var users []models.User
		config.DB.Where("role = ?", models.RoleEmployee).Find(&users)

		for _, u := range users {
			uid := u.ID

			var existingCount int64
			config.DB.Model(&models.Attendance{}).Where("user_id = ? AND date = ?", uid, dateStr).Count(&existingCount)

			if existingCount == 0 {
				holiday := models.Attendance{
					UserID: &uid,
					Date:   dateStr,
					Status: models.StatusLibur,
					Grade:  "-",
				}
				config.DB.Create(&holiday)
			}
		}
		log.Println("Hari ini dinyatakan Libur untuk semua karyawan oleh Sistem")
	}
}
