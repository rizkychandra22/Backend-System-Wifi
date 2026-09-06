package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"math"
	"net/http"
	"strings"
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
		return nil, utils.NewAppError(http.StatusForbidden, "Hari ini sudah dinyatakan libur karena tidak ada yang absen masuk hingga jam 12:00")
	}

	// Cek apakah sudah absen hari ini
	var existing models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&existing).Error; err == nil {
		return nil, utils.NewAppError(http.StatusConflict, "Anda sudah melakukan absen hari ini")
	}

	// Batas absen masuk adalah jam 12:00
	time1200 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, locWIB)
	if now.After(time1200) {
		return nil, utils.NewAppError(http.StatusForbidden, "Batas waktu absen masuk telah lewat (12:00)")
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
	time750 := time.Date(now.Year(), now.Month(), now.Day(), 7, 50, 0, 0, locWIB)
	time800 := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, locWIB)
	time810 := time.Date(now.Year(), now.Month(), now.Day(), 8, 10, 0, 0, locWIB)

	if now.Before(time750) || now.Equal(time750) {
		grade = "Disiplin"
	} else if (now.After(time750) && now.Before(time800)) || now.Equal(time800) {
		grade = "Tepat Waktu"
	} else if (now.After(time800) && now.Before(time810)) || now.Equal(time810) {
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

		time1200 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, locWIB)
		if now.Before(time1200) {
			return nil, utils.NewAppError(http.StatusForbidden, "Izin setengah hari (Halfday) hanya dapat diajukan setelah batas jam 12:00 siang.")
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
		// Izin sebelum absen masuk (Harus sebelum 12:00)
		time1200 := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, locWIB)
		if now.After(time1200) {
			return nil, utils.NewAppError(http.StatusForbidden, "Batas waktu pengajuan izin full-day (12:00) telah lewat. Jika sudah masuk, pastikan absen masuk terlebih dahulu.")
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

type UpdateAttendanceInput struct {
	ClockIn  *string `json:"clock_in"`
	ClockOut *string `json:"clock_out"`
	Grade    *string `json:"grade"`
	Status   *string `json:"status"`
	Notes    *string `json:"notes"`
}

func UpdateAttendance(id string, input UpdateAttendanceInput) (*models.Attendance, *utils.AppError) {
	var attendance models.Attendance
	if err := config.DB.Preload("User").First(&attendance, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Data absensi tidak ditemukan")
	}

	attendanceDate, err := time.ParseInLocation("2006-01-02", attendance.Date, locWIB)
	if err != nil {
		attendanceDate = time.Now().In(locWIB)
	}

	// Update ClockIn
	if input.ClockIn != nil {
		str := strings.TrimSpace(*input.ClockIn)
		if str == "" {
			attendance.ClockIn = nil
		} else {
			if t, err := time.ParseInLocation("15:04", str, locWIB); err == nil {
				parsed := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(), t.Hour(), t.Minute(), 0, 0, locWIB)
				attendance.ClockIn = &parsed
			} else if t, err := time.ParseInLocation("15:04:05", str, locWIB); err == nil {
				parsed := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(), t.Hour(), t.Minute(), t.Second(), 0, locWIB)
				attendance.ClockIn = &parsed
			} else if t, err := time.Parse(time.RFC3339, str); err == nil {
				tInLoc := t.In(locWIB)
				attendance.ClockIn = &tInLoc
			} else {
				return nil, utils.NewAppError(http.StatusBadRequest, "Format jam masuk tidak valid (gunakan format HH:mm)")
			}
		}
	}

	// Update ClockOut
	if input.ClockOut != nil {
		str := strings.TrimSpace(*input.ClockOut)
		if str == "" {
			attendance.ClockOut = nil
		} else {
			if t, err := time.ParseInLocation("15:04", str, locWIB); err == nil {
				parsed := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(), t.Hour(), t.Minute(), 0, 0, locWIB)
				attendance.ClockOut = &parsed
			} else if t, err := time.ParseInLocation("15:04:05", str, locWIB); err == nil {
				parsed := time.Date(attendanceDate.Year(), attendanceDate.Month(), attendanceDate.Day(), t.Hour(), t.Minute(), t.Second(), 0, locWIB)
				attendance.ClockOut = &parsed
			} else if t, err := time.Parse(time.RFC3339, str); err == nil {
				tInLoc := t.In(locWIB)
				attendance.ClockOut = &tInLoc
			} else {
				return nil, utils.NewAppError(http.StatusBadRequest, "Format jam keluar tidak valid (gunakan format HH:mm)")
			}
		}
	}

	// Update Grade
	if input.Grade != nil && strings.TrimSpace(*input.Grade) != "" {
		attendance.Grade = strings.TrimSpace(*input.Grade)
	} else if attendance.ClockIn != nil {
		cIn := *attendance.ClockIn
		time750 := time.Date(cIn.Year(), cIn.Month(), cIn.Day(), 7, 50, 0, 0, locWIB)
		time800 := time.Date(cIn.Year(), cIn.Month(), cIn.Day(), 8, 0, 0, 0, locWIB)
		time810 := time.Date(cIn.Year(), cIn.Month(), cIn.Day(), 8, 10, 0, 0, locWIB)

		if cIn.Before(time750) || cIn.Equal(time750) {
			attendance.Grade = "Disiplin"
		} else if (cIn.After(time750) && cIn.Before(time800)) || cIn.Equal(time800) {
			attendance.Grade = "Tepat Waktu"
		} else if (cIn.After(time800) && cIn.Before(time810)) || cIn.Equal(time810) {
			attendance.Grade = "Toleransi Terlambat"
		} else {
			attendance.Grade = "Terlambat"
		}
	}

	// Update Status
	if input.Status != nil && strings.TrimSpace(*input.Status) != "" {
		attendance.Status = models.AttendanceStatus(strings.TrimSpace(*input.Status))
	} else {
		if attendance.ClockIn != nil && attendance.ClockOut != nil {
			attendance.Status = models.StatusHadir
		} else if attendance.ClockIn != nil && attendance.ClockOut == nil {
			attendance.Status = models.StatusProses
		}
	}

	// Update Notes
	if input.Notes != nil {
		attendance.Notes = input.Notes
	}

	if err := config.DB.Save(&attendance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui data absensi")
	}

	return &attendance, nil
}
