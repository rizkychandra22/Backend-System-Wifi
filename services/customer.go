package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"errors"
	"net/http"

	"gorm.io/gorm"
)


func GetCustomerSubscription(customerID string) (*models.Subscription, *utils.AppError) {
	var sub models.Subscription
	err := config.DB.Preload("WifiPackage").Preload("Customer").Where("customer_id = ?", customerID).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data langganan")
	}
	return &sub, nil
}

