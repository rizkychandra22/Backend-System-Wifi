package helpers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"log"
	"time"
)

func AutoCheckout() {
	locWIB, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")
	autoCheckoutTime := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, locWIB)
	
	config.DB.Model(&models.Attendance{}).
		Where("date = ? AND status = ?", dateStr, models.StatusHadir).
		Updates(map[string]interface{}{
			"status":    models.StatusSelesai,
			"clock_out": autoCheckoutTime,
		})
	log.Println("Proses auto absen keluar pada jam 17:00 selesai")
}
