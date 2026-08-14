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
		return nil, utils.NewAppError(http.StatusConflict, "Phone number already in use")
	}

	user := models.User{
		Name:           name,
		Phone:          phone,
		Role:           models.RoleCustomer,
		Address:        address,
		RegisteredByID: &registeredByID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to create customer")
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
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve customers")
	}
	return users, nil
}
