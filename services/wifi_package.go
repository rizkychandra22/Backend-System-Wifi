package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)

func CreateWifiPackage(name string, price float64) (*models.WifiPackage, *utils.AppError) {
	ws := models.WifiPackage{Name: name, Price: price}
	if err := config.DB.Create(&ws).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to create wifi service")
	}
	return &ws, nil
}

func GetWifiPackages() ([]models.WifiPackage, *utils.AppError) {
	var ws []models.WifiPackage
	if err := config.DB.Find(&ws).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve wifi services")
	}
	return ws, nil
}

func UpdateWifiPackage(id string, name *string, price *float64) (*models.WifiPackage, *utils.AppError) {
	var ws models.WifiPackage
	if err := config.DB.First(&ws, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Wifi service not found")
	}

	if name != nil {
		ws.Name = *name
	}
	if price != nil {
		ws.Price = *price
	}
	if err := config.DB.Save(&ws).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to update wifi service")
	}
	return &ws, nil
}

func DeleteWifiPackage(id string) *utils.AppError {
	var ws models.WifiPackage
	if err := config.DB.First(&ws, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Wifi service not found")
	}
	if err := config.DB.Delete(&ws).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Failed to delete wifi service")
	}
	return nil
}
