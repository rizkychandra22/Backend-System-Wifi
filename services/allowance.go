package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"errors"
	"net/http"
	"regexp"

	"gorm.io/gorm"
)

var monthRegex = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

func CreateAllowance(title, description string, amount float64, month, targetType string, userID *uint, createdByID uint) (*models.Allowance, *utils.AppError) {
	if title == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "Nama tunjangan/bonus tidak boleh kosong")
	}

	if amount <= 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "Nominal tunjangan/bonus harus lebih besar dari 0")
	}

	if !monthRegex.MatchString(month) {
		return nil, utils.NewAppError(http.StatusBadRequest, "Format bulan tidak valid, gunakan format YYYY-MM (cth: 2026-09)")
	}

	if targetType != "personal" && targetType != "global" {
		return nil, utils.NewAppError(http.StatusBadRequest, "Target penerima harus berupa 'personal' atau 'global'")
	}

	var targetUserID *uint
	if targetType == "personal" {
		if userID == nil || *userID == 0 {
			return nil, utils.NewAppError(http.StatusBadRequest, "Karyawan harus dipilih jika target adalah personal")
		}
		var user models.User
		if err := config.DB.First(&user, *userID).Error; err != nil {
			return nil, utils.NewAppError(http.StatusNotFound, "Karyawan yang dipilih tidak ditemukan")
		}
		if user.Role != models.RoleEmployee {
			return nil, utils.NewAppError(http.StatusBadRequest, "Tunjangan hanya dapat diberikan kepada pengguna dengan role karyawan")
		}
		targetUserID = userID
	} else {
		targetUserID = nil
	}

	allowance := models.Allowance{
		Title:       title,
		Description: description,
		Amount:      amount,
		Month:       month,
		TargetType:  targetType,
		UserID:      targetUserID,
		CreatedByID: createdByID,
	}

	if err := config.DB.Create(&allowance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal menyimpan data tunjangan/bonus")
	}

	config.DB.Preload("User").Preload("CreatedBy").First(&allowance, allowance.ID)

	return &allowance, nil
}

func GetAllAllowances(monthFilter string) ([]models.Allowance, *utils.AppError) {
	var allowances []models.Allowance
	query := config.DB.Preload("User").Preload("CreatedBy")

	if monthFilter != "" && monthFilter != "all" {
		query = query.Where("month = ?", monthFilter)
	}

	if err := query.Order("month DESC, id DESC").Find(&allowances).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data tunjangan/bonus")
	}

	return allowances, nil
}

func GetEmployeeAllowances(employeeID uint, monthFilter string) ([]models.Allowance, *utils.AppError) {
	var allowances []models.Allowance
	query := config.DB.Preload("CreatedBy").Where("user_id = ? OR target_type = 'global'", employeeID)

	if monthFilter != "" && monthFilter != "all" {
		query = query.Where("month = ?", monthFilter)
	}

	if err := query.Order("month DESC, id DESC").Find(&allowances).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data tunjangan/bonus karyawan")
	}

	return allowances, nil
}

func GetAllowanceByID(id string) (*models.Allowance, *utils.AppError) {
	var allowance models.Allowance
	if err := config.DB.Preload("User").Preload("CreatedBy").First(&allowance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusNotFound, "Data tunjangan/bonus tidak ditemukan")
		}
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data tunjangan/bonus")
	}

	return &allowance, nil
}

func UpdateAllowance(id string, title, description string, amount float64, month, targetType string, userID *uint) (*models.Allowance, *utils.AppError) {
	var allowance models.Allowance
	if err := config.DB.First(&allowance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(http.StatusNotFound, "Data tunjangan/bonus tidak ditemukan")
		}
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mencari data tunjangan/bonus")
	}

	if title == "" {
		return nil, utils.NewAppError(http.StatusBadRequest, "Nama tunjangan/bonus tidak boleh kosong")
	}

	if amount <= 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "Nominal tunjangan/bonus harus lebih besar dari 0")
	}

	if !monthRegex.MatchString(month) {
		return nil, utils.NewAppError(http.StatusBadRequest, "Format bulan tidak valid, gunakan format YYYY-MM (cth: 2026-09)")
	}

	if targetType != "personal" && targetType != "global" {
		return nil, utils.NewAppError(http.StatusBadRequest, "Target penerima harus berupa 'personal' atau 'global'")
	}

	var targetUserID *uint
	if targetType == "personal" {
		if userID == nil || *userID == 0 {
			return nil, utils.NewAppError(http.StatusBadRequest, "Karyawan harus dipilih jika target adalah personal")
		}
		var user models.User
		if err := config.DB.First(&user, *userID).Error; err != nil {
			return nil, utils.NewAppError(http.StatusNotFound, "Karyawan yang dipilih tidak ditemukan")
		}
		if user.Role != models.RoleEmployee {
			return nil, utils.NewAppError(http.StatusBadRequest, "Tunjangan hanya dapat diberikan kepada pengguna dengan role karyawan")
		}
		targetUserID = userID
	} else {
		targetUserID = nil
	}

	allowance.Title = title
	allowance.Description = description
	allowance.Amount = amount
	allowance.Month = month
	allowance.TargetType = targetType
	allowance.UserID = targetUserID

	if err := config.DB.Save(&allowance).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui data tunjangan/bonus")
	}

	config.DB.Preload("User").Preload("CreatedBy").First(&allowance, allowance.ID)
	return &allowance, nil
}

func DeleteAllowance(id string) *utils.AppError {
	var allowance models.Allowance
	if err := config.DB.First(&allowance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(http.StatusNotFound, "Data tunjangan/bonus tidak ditemukan")
		}
		return utils.NewAppError(http.StatusInternalServerError, "Gagal mencari data tunjangan/bonus")
	}

	if err := config.DB.Delete(&allowance).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menghapus data tunjangan/bonus")
	}

	return nil
}
