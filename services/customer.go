package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)


func GetCustomerSubscription(customerID string) (*models.Subscription, *utils.AppError) {
	var sub models.Subscription
	err := config.DB.Preload("WifiPackage").Preload("Customer").Where("customer_id = ?", customerID).First(&sub).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Tidak ditemukan langganan aktif")
	}
	return &sub, nil
}

