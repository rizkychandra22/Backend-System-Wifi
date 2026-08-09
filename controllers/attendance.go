package controllers

import (
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	attendance, appErr := services.ClockIn(userID.(float64), input.Lat, input.Lng)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
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

	attendance, appErr := services.ClockOut(userID.(float64), input.Lat, input.Lng)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
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

	attendance, appErr := services.RequestIzin(userID.(float64), input.Notes)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
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

	attendance, appErr := services.GetTodayAttendance(userID.(float64))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
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

	history, appErr := services.GetAttendanceHistory(userID.(float64))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": history})
}

// GetAllAttendance handler for admin
func GetAllAttendance(c *gin.Context) {
	records, appErr := services.GetAllAttendance()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}
