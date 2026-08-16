package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create User (Admin only)
func CreateUser(c *gin.Context) {
	var input struct {
		Name    string      `json:"name" binding:"required"`
		Phone   string      `json:"phone" binding:"required"`
		Role    models.Role `json:"role" binding:"required"`
		Address *string     `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDClaim, _ := c.Get("userID")
	userID := uint(userIDClaim.(float64))

	userRoleClaim, _ := c.Get("userRole")
	userRole := userRoleClaim.(string)

	if userRole == string(models.RoleEmployee) {
		input.Role = models.RoleCustomer // Employee hanya bisa membuat customer
	}

	user, appErr := services.CreateUser(input.Name, input.Phone, input.Role, input.Address, &userID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pengguna berhasil dibuat", "data": user})
}

// Get all users
func GetUsers(c *gin.Context) {
	userRoleClaim, _ := c.Get("userRole")
	userRole := userRoleClaim.(string)

	var employeeID *uint
	if userRole == string(models.RoleEmployee) {
		userIDClaim, _ := c.Get("userID")
		id := uint(userIDClaim.(float64))
		employeeID = &id
	}

	users, appErr := services.GetUsers(employeeID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// Get admin contact
func GetAdminContact(c *gin.Context) {
	admin, appErr := services.GetAdminContact()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	// Hanya return nama dan phone number untuk privasi
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"name": admin.Name, "phone": admin.Phone}})
}

// Get user by ID
func GetUserByID(c *gin.Context) {
	user, appErr := services.GetUserByID(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Update User
func UpdateUser(c *gin.Context) {
	var input struct {
		Name    *string      `json:"name"`
		Phone   *string      `json:"phone"`
		Role    *models.Role `json:"role"`
		Address *string      `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRoleClaim, _ := c.Get("userRole")
	userRole := userRoleClaim.(string)

	var employeeID *uint
	if userRole == string(models.RoleEmployee) {
		userIDClaim, _ := c.Get("userID")
		id := uint(userIDClaim.(float64))
		employeeID = &id

		// Karyawan tidak bisa mengubah role
		custRole := models.RoleCustomer
		input.Role = &custRole
	}

	user, appErr := services.UpdateUser(c.Param("id"), input.Name, input.Phone, input.Role, input.Address, employeeID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil diperbarui", "data": user})
}

// Delete User
func DeleteUser(c *gin.Context) {
	appErr := services.DeleteUser(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}

// Reset User IP (Admin only) - Digunakan jika device pengguna hilang/rusak
func ResetUserIP(c *gin.Context) {
	appErr := services.ResetUserIP(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "IP Pengguna berhasil di-reset. Pengguna sekarang dapat login dari perangkat baru."})
}
