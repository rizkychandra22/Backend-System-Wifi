package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
	"time"
)

func CreateOrUpdateSubscription(customerID uint, wifiPackageID uint, billingDay int, nextDueDate time.Time, status string) (*models.Subscription, *utils.AppError) {
	// 1. Validasi customer ada dan role-nya customer
	var customer models.User
	if err := config.DB.First(&customer, customerID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Pelanggan tidak ditemukan")
	}
	if customer.Role != models.RoleCustomer {
		return nil, utils.NewAppError(http.StatusBadRequest, "Pengguna terpilih bukan merupakan pelanggan")
	}

	// 2. Validasi paket wifi ada
	var wifiPkg models.WifiPackage
	if err := config.DB.First(&wifiPkg, wifiPackageID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Paket WiFi tidak ditemukan")
	}

	// 3. Validasi status
	validStatus := map[string]bool{"active": true, "suspended": true, "cancelled": true}
	if !validStatus[status] {
		status = "active" // fallback
	}

	// 4. Cari apakah sudah ada langganan aktif
	var sub models.Subscription
	err := config.DB.Where("customer_id = ?", customerID).First(&sub).Error

	if err == nil {
		// Update existing subscription
		sub.WifiPackageID = wifiPackageID
		sub.BillingDay = billingDay
		sub.NextDueDate = nextDueDate
		sub.Status = status
		if err := config.DB.Save(&sub).Error; err != nil {
			return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui langganan")
		}
	} else {
		// Create new subscription
		sub = models.Subscription{
			CustomerID:    customerID,
			WifiPackageID: wifiPackageID,
			BillingDay:    billingDay,
			NextDueDate:   nextDueDate,
			Status:        status,
		}
		if err := config.DB.Create(&sub).Error; err != nil {
			return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal membuat langganan baru")
		}
	}

	// Preload data lengkap untuk response
	config.DB.Preload("WifiPackage").Preload("Customer").First(&sub, sub.ID)
	return &sub, nil
}

func GetSubscriptionByCustomerID(customerID string) (*models.Subscription, *utils.AppError) {
	var sub models.Subscription
	err := config.DB.Preload("WifiPackage").Preload("Customer").Where("customer_id = ?", customerID).First(&sub).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Data langganan aktif tidak ditemukan")
	}
	return &sub, nil
}

func GetAllSubscriptions() ([]models.Subscription, *utils.AppError) {
	var subs []models.Subscription
	err := config.DB.Preload("WifiPackage").Preload("Customer").Order("created_at DESC").Find(&subs).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil semua data langganan")
	}
	return subs, nil
}

func DeleteSubscription(id string) *utils.AppError {
	var sub models.Subscription
	if err := config.DB.First(&sub, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Data langganan tidak ditemukan")
	}

	if err := config.DB.Delete(&sub).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menghapus data langganan")
	}

	return nil
}
