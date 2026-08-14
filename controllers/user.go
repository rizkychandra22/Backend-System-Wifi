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
	userID         := uint(userIDClaim.(float64))

	user, appErr := services.CreateUser(input.Name, input.Phone, input.Role, input.Address, &userID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}

// Get all users
func GetUsers(c *gin.Context) {
	users, appErr := services.GetUsers()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
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

	user, appErr := services.UpdateUser(c.Param("id"), input.Name, input.Phone, input.Role, input.Address)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": user})
}

// Delete User
func DeleteUser(c *gin.Context) {
	appErr := services.DeleteUser(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Reset User IP (Admin only) - Digunakan jika device pengguna hilang/rusak
func ResetUserIP(c *gin.Context) {
	appErr := services.ResetUserIP(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User IP reset successfully. User can now login from a new device."})
}
