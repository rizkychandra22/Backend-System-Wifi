package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func UpdateProfile(userID float64, name, phone, address *string) (*models.User, *utils.AppError) {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}

	if name != nil {
		user.Name = *name
	}
	if phone != nil {
		user.Phone = *phone
	}
	if address != nil {
		user.Address = address
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui profil")
	}

	return &user, nil
}

func UpdatePassword(userID float64, oldPassword, newPassword string) *utils.AppError {
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Pengguna tidak ditemukan")
	}

	if user.Role != "admin" {
		return utils.NewAppError(http.StatusForbidden, "Fitur ganti password hanya untuk admin")
	}

	if user.Password == nil || bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(oldPassword)) != nil {
		return utils.NewAppError(http.StatusUnauthorized, "Password lama salah")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal mengenkripsi password")
	}

	passwordStr := string(hashedPassword)
	user.Password = &passwordStr

	if err := config.DB.Save(&user).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menyimpan password baru")
	}

	return nil
}
