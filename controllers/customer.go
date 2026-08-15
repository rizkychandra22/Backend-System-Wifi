package controllers

import (
	"backend-wifi/models"
	"backend-wifi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCustomer(c *gin.Context) {
	var input struct {
		Name    string  `json:"name" binding:"required"`
		Phone   string  `json:"phone" binding:"required"`
		Address *string `json:"address"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDClaim, _ := c.Get("userID")
	userID := uint(userIDClaim.(float64))

	userRoleClaim, _ := c.Get("userRole")
	userRole := userRoleClaim.(string)
	if userRole != string(models.RoleAdmin) && userRole != string(models.RoleEmployee) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to create customer"})
		return
	}

	customer, appErr := services.CreateCustomer(input.Name, input.Phone, input.Address, userID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pelanggan berhasil dibuat", "data": customer})
}

func GetCustomers(c *gin.Context) {
	userRoleClaim, _ := c.Get("userRole")
	userRole := userRoleClaim.(string)

	var employeeID *uint
	if userRole == string(models.RoleEmployee) {
		userIDClaim, _ := c.Get("userID")
		id := uint(userIDClaim.(float64))
		employeeID = &id
	}

	customers, appErr := services.GetCustomers(employeeID)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

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
