package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)

func CreateUser(name, phone string, role models.Role, address *string) (*models.User, *utils.AppError) {
	user := models.User{
		Name:    name,
		Phone:   phone,
		Role:    role,
		Address: address,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to create user")
	}

	return &user, nil
}

func GetUsers() ([]models.User, *utils.AppError) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to retrieve users")
	}
	return users, nil
}

func GetUserByID(id string) (*models.User, *utils.AppError) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "User not found")
	}
	return &user, nil
}

func UpdateUser(id string, name, phone *string, role *models.Role, address *string) (*models.User, *utils.AppError) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "User not found")
	}

	if name != nil {
		user.Name = *name
	}
	if phone != nil {
		user.Phone = *phone
	}
	if role != nil {
		user.Role = *role
	}
	if address != nil {
		user.Address = address
	}
	if err := config.DB.Save(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to update user")
	}

	return &user, nil
}

func DeleteUser(id string) *utils.AppError {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "User not found")
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Failed to delete user")
	}

	return nil
}

func ResetUserIP(id string) *utils.AppError {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "User not found")
	}

	user.IPAddress = nil
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Failed to reset user IP")
	}

	return nil
}
