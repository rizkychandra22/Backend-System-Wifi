package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetCustomerSubscription(c *gin.Context) {
	customerID := c.Param("id")

	userRoleClaim, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userRole := userRoleClaim.(string)

	if userRole == string(models.RoleCustomer) {
		userIDClaim, _ := c.Get("userID")
		userID := uint(userIDClaim.(float64))

		if fmt.Sprintf("%d", userID) != customerID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own subscription"})
			return
		}
	}

	subscription, appErr := services.GetCustomerSubscription(customerID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

