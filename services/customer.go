package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)


func GetCustomerSubscription(customerID string) (*models.WifiPackage, *utils.AppError) {
	var payment models.Payment
	err := config.DB.Preload("WifiPackage").Where("customer_id = ?", customerID).Order("created_at desc").First(&payment).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Tidak ditemukan langganan aktif")
	}
	return payment.WifiPackage, nil
}

