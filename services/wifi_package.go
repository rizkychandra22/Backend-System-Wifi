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
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal membuat layanan WiFi")
	}
	return &ws, nil
}

func GetWifiPackages() ([]models.WifiPackage, *utils.AppError) {
	var ws []models.WifiPackage
	if err := config.DB.Find(&ws).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data layanan WiFi")
	}
	return ws, nil
}

func UpdateWifiPackage(id string, name *string, price *float64) (*models.WifiPackage, *utils.AppError) {
	var ws models.WifiPackage
	if err := config.DB.First(&ws, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Layanan WiFi tidak ditemukan")
	}

	if name != nil {
		ws.Name = *name
	}
	if price != nil {
		ws.Price = *price
	}
	if err := config.DB.Save(&ws).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui layanan WiFi")
	}
	return &ws, nil
}

func DeleteWifiPackage(id string) *utils.AppError {
	var ws models.WifiPackage
	if err := config.DB.First(&ws, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Layanan WiFi tidak ditemukan")
	}
	if err := config.DB.Delete(&ws).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menghapus layanan WiFi")
	}
	return nil
}
