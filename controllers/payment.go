package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePayment(c *gin.Context) {
	var input struct {
		CustomerID    uint   `json:"customer_id" binding:"required"`
		WifiPackageID uint   `json:"wifi_package_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDClaim, _ := c.Get("userID")
	userID := uint(userIDClaim.(float64))

	payment, appErr := services.CreatePayment(input.CustomerID, input.WifiPackageID, userID, input.PaymentMethod)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pembayaran berhasil dicatat", "data": payment})
}

func GetCustomerPayments(c *gin.Context) {
	userRoleClaim, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userRole := userRoleClaim.(string)
	
	if userRole == string(models.RoleCustomer) {
		userIDClaim, _ := c.Get("userID")
		userID := uint(userIDClaim.(float64))
		
		if fmt.Sprintf("%d", userID) != c.Param("customer_id") {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own payments"})
			return
		}
	}

	payments, appErr := services.GetCustomerPayments(c.Param("customer_id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}


func GetAllPayments(c *gin.Context) {
	payments, appErr := services.GetAllPayments()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

func UpdatePayment(c *gin.Context) {
	var input struct {
		CustomerID    uint   `json:"customer_id" binding:"required"`
		WifiPackageID uint   `json:"wifi_package_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID := c.Param("id")
	payment, appErr := services.UpdatePayment(paymentID, input.CustomerID, input.WifiPackageID, input.PaymentMethod)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil diperbarui", "data": payment})
}

func DeletePayment(c *gin.Context) {
	paymentID := c.Param("id")
	appErr := services.DeletePayment(paymentID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pembayaran berhasil dihapus"})
}
