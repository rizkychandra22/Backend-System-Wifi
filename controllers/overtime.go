package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"backend-wifi/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OvertimeInput struct {
	UserID      uint   `json:"user_id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Date        string `json:"date" binding:"required"` // format: YYYY-MM-DD
	StartTime   string `json:"start_time" binding:"required"` // format: HH:mm
	EndTime     string `json:"end_time" binding:"required"` // format: HH:mm
}

func CreateOvertime(c *gin.Context) {
	var input OvertimeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))

	// Jika employee, paksa userID = diri sendiri. Jika admin, gunakan dari input (jika ada)
	targetUserID := userID
	if userRole == string(models.RoleAdmin) && input.UserID != 0 {
		targetUserID = input.UserID
	}

	dateParsed, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah, gunakan YYYY-MM-DD"})
		return
	}

	startTimeParsed, err := time.Parse("2006-01-02 15:04", input.Date+" "+input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format jam mulai salah, gunakan HH:mm"})
		return
	}

	endTimeParsed, err := time.Parse("2006-01-02 15:04", input.Date+" "+input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format jam selesai salah, gunakan HH:mm"})
		return
	}

	overtime, appErr := services.CreateOvertime(targetUserID, input.Title, input.Description, dateParsed, startTimeParsed, endTimeParsed)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil mencatat lembur",
		"data":    overtime,
	})
}

func GetOvertimes(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))

	var overtimes []models.Overtime
	var appErr *utils.AppError

	if userRole == string(models.RoleAdmin) {
		overtimes, appErr = services.GetAllOvertimes()
	} else {
		overtimes, appErr = services.GetUserOvertimes(userID)
	}

	if appErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data lembur"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": overtimes,
	})
}

func GetOvertime(c *gin.Context) {
	id := c.Param("id")
	overtime, appErr := services.GetOvertimeByID(id)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	// Cek akses
	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))

	if userRole != string(models.RoleAdmin) && overtime.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak melihat data ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": overtime,
	})
}

func UpdateOvertime(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	if userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengedit data lembur secara bebas"})
		return
	}

	id := c.Param("id")
	var input OvertimeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid"})
		return
	}

	dateParsed, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal salah, gunakan YYYY-MM-DD"})
		return
	}

	startTimeParsed, err := time.Parse("2006-01-02 15:04", input.Date+" "+input.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format jam mulai salah, gunakan HH:mm"})
		return
	}

	endTimeParsed, err := time.Parse("2006-01-02 15:04", input.Date+" "+input.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format jam selesai salah, gunakan HH:mm"})
		return
	}

	overtime, appErr := services.UpdateOvertime(id, input.Title, input.Description, dateParsed, startTimeParsed, endTimeParsed, input.UserID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui data lembur",
		"data":    overtime,
	})
}

func DeleteOvertime(c *gin.Context) {
	id := c.Param("id")

	// Cek kepemilikan jika karyawan
	userRole, _ := c.Get("userRole")
	if userRole != string(models.RoleAdmin) {
		userIDVal, _ := c.Get("userID")
		userID := uint(userIDVal.(float64))
		
		overtime, err := services.GetOvertimeByID(id)
		if err != nil {
			c.JSON(err.StatusCode, gin.H{"error": err.Message})
			return
		}
		
		if overtime.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak menghapus data ini"})
			return
		}
	}

	if appErr := services.DeleteOvertime(id); appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus data lembur",
	})
}
