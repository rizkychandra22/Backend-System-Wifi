package controllers

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login function for all users (using phone number)
func Login(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("phone = ?", input.Phone).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nomor telepon tidak terdaftar"})
		return
	}

	clientIP := c.ClientIP()

	// Device IP Lock Logic (Admin dibebaskan dari lock device)
	if user.Role != "admin" {
		if user.IPAddress == nil || *user.IPAddress == "" {
			// First time login, save this IP
			user.IPAddress = &clientIP
			config.DB.Save(&user)
		} else if *user.IPAddress != clientIP {
			// IP mismatch, block login
			c.JSON(http.StatusForbidden, gin.H{"error": "Akun ini sudah login di device lain. Silakan hubungi Admin."})
			return
		}
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"role":  user.Role,
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 3 days expiration
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT_SECRET not configured"})
		return
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"phone":   user.Phone,
			"role":    user.Role,
			"address": user.Address,
		},
	})
}
