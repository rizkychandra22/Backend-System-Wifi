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
	config.DB.Model(&models.Attendance{}).Where("date = ? AND status = ?", dateStr, models.StatusHadir).Count(&count)
	if count == 0 {
		holiday := models.Attendance{
			Date:   dateStr,
			Status: models.StatusLibur,
		}
		config.DB.Create(&holiday)
		log.Println("Hari ini otomatis diset sebagai Hari Libur")
	}
}
