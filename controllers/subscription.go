package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrUpdateSubscription(c *gin.Context) {
	var input struct {
		CustomerID    uint      `json:"customer_id" binding:"required"`
		WifiPackageID uint      `json:"wifi_package_id" binding:"required"`
		BillingDay    int       `json:"billing_day" binding:"required,min=1,max=31"`
		NextDueDate   time.Time `json:"next_due_date" binding:"required"`
		Status        string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDClaim, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDClaim.(float64))

	subscription, appErr := services.CreateOrUpdateSubscription(
		input.CustomerID,
		input.WifiPackageID,
		input.BillingDay,
		input.NextDueDate,
		input.Status,
		userID,
	)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data langganan berhasil disimpan",
		"data":    subscription,
	})
}

func GetSubscriptionByCustomerID(c *gin.Context) {
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
			c.JSON(http.StatusForbidden, gin.H{"error": "Anda hanya bisa melihat data langganan Anda sendiri"})
			return
		}
	}

	subscription, appErr := services.GetSubscriptionByCustomerID(customerID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscription})
}

func GetAllSubscriptions(c *gin.Context) {
	// Hanya Admin dan Karyawan yang boleh melihat semua subscription
	userRoleClaim, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userRole := userRoleClaim.(string)
	if userRole == string(models.RoleCustomer) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	subscriptions, appErr := services.GetAllSubscriptions()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subscriptions})
}

func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	userRoleClaim, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userRole := userRoleClaim.(string)
	if userRole != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang diperbolehkan menghapus data langganan"})
		return
	}

	appErr := services.DeleteSubscription(id)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data langganan berhasil dihapus"})
}
