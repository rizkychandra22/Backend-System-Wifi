package controllers

import (
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateProfile allows a user to update their own profile
func UpdateProfile(c *gin.Context) {
	// Get user ID from JWT context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Name    *string `json:"name"`
		Phone   *string `json:"phone"`
		Address *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, appErr := services.UpdateProfile(userID.(float64), input.Name, input.Phone, input.Address)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
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

	var input struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appErr := services.UpdatePassword(userID.(float64), input.OldPassword, input.NewPassword)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
