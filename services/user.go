package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
)

func CreateUser(name, phone string, role models.Role, address *string, registeredByID *uint) (*models.User, *utils.AppError) {
	user := models.User{
		Name:           name,
		Phone:          phone,
		Role:           role,
		Address:        address,
		RegisteredByID: registeredByID,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal membuat pengguna")
	}

	return &user, nil
}

func GetUsers() ([]models.User, *utils.AppError) {
	var users []models.User
	if err := config.DB.Preload("RegisteredBy").Find(&users).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data pengguna")
	}
	return users, nil
}

func GetUserByID(id string) (*models.User, *utils.AppError) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}
	return &user, nil
}

func UpdateUser(id string, name, phone *string, role *models.Role, address *string) (*models.User, *utils.AppError) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
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
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui pengguna")
	}

	return &user, nil
}

func DeleteUser(id string) *utils.AppError {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menghapus pengguna")
	}

	return nil
}

func ResetUserIP(id string) *utils.AppError {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}

	user.IPAddress = nil
	if err := config.DB.Save(&user).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal me-reset IP pengguna")
	}

	return nil
}
