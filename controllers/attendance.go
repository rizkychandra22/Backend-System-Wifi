package controllers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

// ClockIn handler
func ClockIn(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Lat float64 `json:"lat" binding:"required"`
		Lng float64 `json:"lng" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lat and Lng are required"})
		return
	}

	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	// Cek apakah hari ini libur otomatis
	var checkHoliday models.Attendance
	if err := config.DB.Where("date = ? AND status = ?", dateStr, models.StatusLibur).First(&checkHoliday).Error; err == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hari ini sudah dinyatakan libur karena tidak ada yang absen masuk hingga jam 08:30"})
		return
	}

	// Cek apakah sudah absen hari ini
	var existing models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah melakukan absen hari ini"})
		return
	}

	// Batas absen masuk adalah jam 08:30 (jika absen lebih dari 08:30 dianggap tidak bisa absen/sudah libur otomatis)
	// Kita toleransi jika mungkin scheduler telat, tetap tolak jika jam > 8:30
	time830 := time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, locWIB)
	if now.After(time830) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Batas waktu absen masuk telah lewat (08:30)"})
		return
	}

	// Validasi Jarak (maksimal 100 meter dari kantor)
	officeLat := -7.033562
	officeLng := 106.949204
	distance := haversineDistance(officeLat, officeLng, input.Lat, input.Lng)
	if distance > 100 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda berada di luar area kantor (jarak > 100 meter)"})
		return
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

	uid := uint(userID.(float64))
	attendance := models.Attendance{
		UserID:     &uid,
		Date:       dateStr,
		ClockIn:    &now,
		Grade:      grade,
		ClockInLat: &input.Lat,
		ClockInLng: &input.Lng,
		Status:     models.StatusHadir,
	}

	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat absen masuk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Absen masuk berhasil",
		"data":    attendance,
	})
}

// ClockOut handler
func ClockOut(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Lat float64 `json:"lat" binding:"required"`
		Lng float64 `json:"lng" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lat and Lng are required"})
		return
	}

	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	var attendance models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&attendance).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anda belum melakukan absen masuk hari ini"})
		return
	}

	if attendance.ClockOut != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah melakukan absen keluar hari ini"})
		return
	}

	// Jam keluar harus antara 16:00 sampai 18:00
	time1600 := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, locWIB)
	time1800 := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, locWIB)

	if now.Before(time1600) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Absen keluar hanya dapat dilakukan mulai jam 16:00"})
		return
	}
	if now.After(time1800) {
		// Logika ini mungkin disetel dari scheduler (otomatis set selesai jika telat)
		// Namun jika mereka panggil manual via API dan sudah telat:
		c.JSON(http.StatusForbidden, gin.H{"error": "Batas waktu absen keluar (18:00) telah lewat"})
		return
	}

	// Validasi Jarak (maksimal 100 meter dari kantor)
	officeLat := -7.033562
	officeLng := 106.949204
	distance := haversineDistance(officeLat, officeLng, input.Lat, input.Lng)
	if distance > 100 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Jarak absen keluar berada di luar area kantor (> 100 meter)"})
		return
	}

	attendance.ClockOut = &now
	attendance.ClockOutLat = &input.Lat
	attendance.ClockOutLng = &input.Lng
	attendance.Status = models.StatusSelesai

	if err := config.DB.Save(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencatat absen keluar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Absen keluar berhasil",
		"data":    attendance,
	})
}

// RequestIzin handler
func RequestIzin(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Notes string `json:"notes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Notes is required for izin"})
		return
	}

	now := time.Now().In(locWIB)
	dateStr := now.Format("2006-01-02")

	var existing models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah memiliki catatan absen hari ini"})
		return
	}

	uid := uint(userID.(float64))
	attendance := models.Attendance{
		UserID: &uid,
		Date:   dateStr,
		Status: models.StatusIzin,
		Notes:  &input.Notes,
	}

	if err := config.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengajukan izin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Izin berhasil dicatat",
		"data":    attendance,
	})
}

// GetTodayAttendance handler for employee
func GetTodayAttendance(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	dateStr := time.Now().In(locWIB).Format("2006-01-02")

	var attendance models.Attendance
	if err := config.DB.Where("user_id = ? AND date = ?", userID, dateStr).First(&attendance).Error; err != nil {
		// Belum absen
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendance})
}

// GetAttendanceHistory handler for employee
func GetAttendanceHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var history []models.Attendance
	if err := config.DB.Where("user_id = ?", userID).Order("date DESC").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat absen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// GetAllAttendance handler for admin
func GetAllAttendance(c *gin.Context) {
	var records []models.Attendance
	if err := config.DB.Preload("User").Order("date DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data absen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}
