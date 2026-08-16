package controllers

import (
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login function for all users (using phone number)
func Login(c *gin.Context) {
	var input struct {
		Phone    string  `json:"phone" binding:"required"`
		Password *string `json:"password"`
		DeviceID string  `json:"device_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, appErr := services.Login(input.Phone, input.Password, input.DeviceID)
	if appErr != nil {
		if appErr.Data != nil {
			// Merge data from appErr.Data (e.g., requires_password: true)
			dataMap, ok := appErr.Data.(map[string]interface{})
			if ok {
				response := gin.H{"error": appErr.Message}
				for k, v := range dataMap {
					response[k] = v
				}
				c.JSON(appErr.StatusCode, response)
				return
			}
		}
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, result)
}
