package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"math"
	"net/http"
	"time"
)

var locWIB *time.Location

func init() {
	var err error
	locWIB, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		locWIB = time.FixedZone("WIB", 7*3600)
	}
}

// Haversine formula to calculate distance between two coordinates in meters
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Earth radius in meters

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	lat1Rad := lat1 * (math.Pi / 180.0)
	lat2Rad := lat2 * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func ClockIn(userID float64, lat, lng float64) (*models.Attendance, *utils.AppError) {
	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	// Cek apakah hari ini libur otomatis
	var checkHoliday models.Attendance
	if err := config.DB.Where("date = ? AND status = ?", dateStr, models.StatusLibur).First(&checkHoliday).Error; err == nil {
		return nil, utils.NewAppError(http.StatusForbidden, "Hari ini sudah dinyatakan libur karena tidak ada yang absen masuk hingga jam 12:30")
	}

	// Cek apakah sudah absen hari ini
	var existing models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&existing).Error; err == nil {
		return nil, utils.NewAppError(http.StatusConflict, "Anda sudah melakukan absen hari ini")
	}

	// Batas absen masuk adalah jam 12:30
	time1230 := time.Date(now.Year(), now.Month(), now.Day(), 12, 30, 0, 0, locWIB)
	if now.After(time1230) {
		return nil, utils.NewAppError(http.StatusForbidden, "Batas waktu absen masuk telah lewat (12:30)")
	}

	// Validasi Jarak (maksimal 100 meter dari kantor)
	officeLat := -7.033562
	officeLng := 106.949204
	// officeLat := -6.926958
	// officeLng := 106.908998
	distance := haversineDistance(officeLat, officeLng, lat, lng)
	if distance > 100 {
		return nil, utils.NewAppError(http.StatusForbidden, "Anda berada di luar area kantor (jarak > 100 meter)")
	}

	// Tentukan Grade
	var grade string
	time800 := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, locWIB)
	time805 := time.Date(now.Year(), now.Month(), now.Day(), 8, 5, 0, 0, locWIB)

	if now.Before(time800) {
		grade = "Disiplin"
	} else if now.Equal(time800) {
		grade = "Tepat Waktu"
	} else if now.Before(time805) || now.Equal(time805) {
		grade = "Toleransi Terlambat"
	} else {
		grade = "Terlambat"
	}

	uid := uint(userID)
	attendance := models.Attendance{
		UserID:     &uid,
		Date:       dateStr,
		ClockIn:    &now,
		Grade:      grade,
		ClockInLat: &lat,
		ClockInLng: &lng,
		Status:     models.StatusProses,
	}

	if err := config.DB.Create(&attendance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mencatat absen masuk")
	}

	return &attendance, nil
}

func ClockOut(userID float64, lat, lng float64) (*models.Attendance, *utils.AppError) {
	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	var attendance models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ? AND status = ?", userID, dateStr, models.StatusProses).First(&attendance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Anda belum melakukan absen masuk hari ini")
	}

	if attendance.ClockOut != nil {
		return nil, utils.NewAppError(http.StatusConflict, "Anda sudah melakukan absen keluar hari ini")
	}

	// Jam keluar harus antara 16:00 sampai 17:00
	time1600 := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, locWIB)
	time1700 := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, locWIB)

	if now.Before(time1600) {
		return nil, utils.NewAppError(http.StatusForbidden, "Absen keluar hanya dapat dilakukan mulai jam 16:00")
	}
	if now.After(time1700) {
		return nil, utils.NewAppError(http.StatusForbidden, "Batas waktu absen keluar (17:00) telah lewat. Sistem akan otomatis mencatat absen keluar Anda.")
	}

	// Validasi Jarak
	officeLat := -7.033562
	officeLng := 106.949204
	// officeLat := -6.926958
	// officeLng := 106.908998
	distance := haversineDistance(officeLat, officeLng, lat, lng)
	if distance > 100 {
		return nil, utils.NewAppError(http.StatusForbidden, "Jarak absen keluar berada di luar area kantor (> 100 meter)")
	}

	attendance.ClockOut = &now
	attendance.ClockOutLat = &lat
	attendance.ClockOutLng = &lng
	attendance.Status = models.StatusHadir

	if err := config.DB.Save(&attendance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mencatat absen keluar")
	}

	return &attendance, nil
}

func RequestIzin(userID float64, notes string) (*models.Attendance, *utils.AppError) {
	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	var existing models.Attendance
	err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&existing).Error

	uid := uint(userID)

	if err == nil {
		// Mid-day izin (Sudah absen masuk)
		if existing.Status != models.StatusProses {
			return nil, utils.NewAppError(http.StatusConflict, "Status absen tidak valid untuk mengajukan izin (sudah selesai atau libur)")
		}
		
		existing.Status = models.StatusIzin
		existing.ClockOut = &now
		if existing.Notes == nil || *existing.Notes == "" {
			existing.Notes = &notes
		} else {
			updatedNotes := *existing.Notes + " | Izin: " + notes
			existing.Notes = &updatedNotes
		}

		if err := config.DB.Save(&existing).Error; err != nil {
			return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengupdate status menjadi izin")
		}
		return &existing, nil
	} else {
		// Izin sebelum absen masuk (Harus sebelum 12:30)
		time1230 := time.Date(now.Year(), now.Month(), now.Day(), 12, 30, 0, 0, locWIB)
		if now.After(time1230) {
			return nil, utils.NewAppError(http.StatusForbidden, "Batas waktu pengajuan izin full-day (12:30) telah lewat. Jika sudah masuk, pastikan absen masuk terlebih dahulu.")
		}

		attendance := models.Attendance{
			UserID: &uid,
			Date:   dateStr,
			Status: models.StatusIzin,
			Notes:  &notes,
		}

		if err := config.DB.Create(&attendance).Error; err != nil {
			return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengajukan izin")
		}
		return &attendance, nil
	}
}

func GetTodayAttendance(userID float64) (*models.Attendance, *utils.AppError) {
	dateStr := time.Now().In(locWIB).Format("2006-01-02")

	var attendance models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&attendance).Error; err != nil {
		return nil, nil // Belum absen, tidak error tapi kosong
	}

	return &attendance, nil
}

func GetAttendanceHistory(userID float64) ([]models.Attendance, *utils.AppError) {
	var history []models.Attendance
	if err := config.DB.Where("user_id = ?", userID).Order("date DESC").Find(&history).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil riwayat absen")
	}
	return history, nil
}

func GetAllAttendance() ([]models.Attendance, *utils.AppError) {
	var records []models.Attendance
	if err := config.DB.Joins("JOIN users ON users.id = attendances.user_id").
		Where("users.role = ?", models.RoleEmployee).
		Preload("User").Order("date DESC").Find(&records).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data absen")
	}
	return records, nil
}
