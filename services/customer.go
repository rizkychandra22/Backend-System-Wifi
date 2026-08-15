package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)

func CreateCustomer(name, phone string, address *string, registeredByID uint) (*models.User, *utils.AppError) {
	// Cek apakah nomor telepon sudah digunakan
	var existing models.User
	if err := config.DB.Where("phone = ?", phone).First(&existing).Error; err == nil {
		return nil, utils.NewAppError(http.StatusConflict, "Nomor telepon sudah digunakan")
	}

	user := models.User{
		Name:           name,
		Phone:          phone,
		Role:           models.RoleCustomer,
		Address:        address,
		RegisteredByID: &registeredByID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal membuat pelanggan")
	}

	return &user, nil
}

func GetCustomers(employeeID *uint) ([]models.User, *utils.AppError) {
	var users []models.User
	query := config.DB.Preload("RegisteredBy").Where("role = ?", models.RoleCustomer)

	if employeeID != nil {
		query = query.Where("registered_by_id = ?", *employeeID)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data pelanggan")
	}
	return users, nil
}

func GetCustomerSubscription(customerID string) (*models.WifiPackage, *utils.AppError) {
	var payment models.Payment
	err := config.DB.Preload("WifiPackage").Where("customer_id = ?", customerID).Order("created_at desc").First(&payment).Error
	if err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Tidak ditemukan langganan aktif")
	}
	return payment.WifiPackage, nil
}
