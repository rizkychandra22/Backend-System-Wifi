package controllers

import (
	"backend-wifi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateWifiPackage(c *gin.Context) {
	var input struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ws, appErr := services.CreateWifiPackage(input.Name, input.Price)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Wifi service created", "data": ws})
}

func GetWifiPackages(c *gin.Context) {
	ws, appErr := services.GetWifiPackages()
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ws})
}

func UpdateWifiPackage(c *gin.Context) {
	var input struct {
		Name  *string  `json:"name"`
		Price *float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ws, appErr := services.UpdateWifiPackage(c.Param("id"), input.Name, input.Price)
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wifi service updated", "data": ws})
}

func DeleteWifiPackage(c *gin.Context) {
	appErr := services.DeleteWifiPackage(c.Param("id"))
	if appErr != nil {
		c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wifi service deleted"})
}
