package controllers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdateProfile allows a user to update their own profile
func UpdateProfile(c *gin.Context) {
	// Get user ID from JWT context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input struct {
		Name           *string `json:"name"`
		Phone          *string `json:"phone"`
		Address        *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Address != nil {
		user.Address = input.Address
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    user,
	})
}

func UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Fitur ganti password hanya untuk admin"})
		return
	}

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == nil || bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.OldPassword)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password lama salah"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	passwordStr := string(hashedPassword)
	user.Password = &passwordStr

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
