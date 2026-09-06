package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AllowanceInput struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required"`
	Month       string  `json:"month" binding:"required"`
	TargetType  string  `json:"target_type" binding:"required"`
	UserID      *uint   `json:"user_id"`
}

func CreateAllowance(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	if userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat menambahkan tunjangan/bonus"})
		return
	}

	var input AllowanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	userIDVal, _ := c.Get("userID")
	adminID := uint(userIDVal.(float64))

	allowance, appErr := services.CreateAllowance(
		input.Title,
		input.Description,
		input.Amount,
		input.Month,
		input.TargetType,
		input.UserID,
		adminID,
	)

	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menambahkan tunjangan/bonus karyawan",
		"data":    allowance,
	})
}

func GetAllowances(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))
	monthFilter := c.Query("month")

	var allowances []models.Allowance
	var appErr error

	if userRole == string(models.RoleAdmin) {
		var err *models.Allowance // dummy for type
		_ = err
		allList, serviceErr := services.GetAllAllowances(monthFilter)
		allowances = allList
		if serviceErr != nil {
			appErr = serviceErr
		}
	} else if userRole == string(models.RoleEmployee) {
		empList, serviceErr := services.GetEmployeeAllowances(userID, monthFilter)
		allowances = empList
		if serviceErr != nil {
			appErr = serviceErr
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
		return
	}

	if appErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data tunjangan/bonus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": allowances,
	})
}

func GetAllowance(c *gin.Context) {
	id := c.Param("id")
	allowance, appErr := services.GetAllowanceByID(id)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	userRole, _ := c.Get("userRole")
	userIDVal, _ := c.Get("userID")
	userID := uint(userIDVal.(float64))

	// Jika employee, hanya boleh lihat jika global atau milik sendiri
	if userRole == string(models.RoleEmployee) {
		if allowance.TargetType == "personal" && (allowance.UserID == nil || *allowance.UserID != userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": allowance,
	})
}

func UpdateAllowance(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	if userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengubah tunjangan/bonus"})
		return
	}

	var input AllowanceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	id := c.Param("id")
	allowance, appErr := services.UpdateAllowance(
		id,
		input.Title,
		input.Description,
		input.Amount,
		input.Month,
		input.TargetType,
		input.UserID,
	)

	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui data tunjangan/bonus",
		"data":    allowance,
	})
}

func DeleteAllowance(c *gin.Context) {
	userRole, _ := c.Get("userRole")
	if userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat menghapus tunjangan/bonus"})
		return
	}

	id := c.Param("id")
	if appErr := services.DeleteAllowance(id); appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus data tunjangan/bonus",
	})
}
