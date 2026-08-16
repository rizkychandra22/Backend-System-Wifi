package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
	"time"
)

func CreateOvertime(userID uint, title, description string, date, startTime, endTime time.Time) (*models.Overtime, *utils.AppError) {
	// Calculate Price: Rp 5.000 / hour, proportional
	durationMinutes := endTime.Sub(startTime).Minutes()
	if durationMinutes < 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "Waktu selesai harus lebih besar dari waktu mulai")
	}
	price := (durationMinutes / 60.0) * 5000.0

	overtime := models.Overtime{
		UserID:      userID,
		Title:       title,
		Description: description,
		Date:        date,
		StartTime:   startTime,
		EndTime:     endTime,
		Price:       price,
	}

	if err := config.DB.Create(&overtime).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mencatat lembur")
	}

	return &overtime, nil
}

func GetAllOvertimes() ([]models.Overtime, *utils.AppError) {
	var overtimes []models.Overtime
	if err := config.DB.Preload("User").Order("date desc, start_time desc").Find(&overtimes).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data lembur")
	}
	return overtimes, nil
}

func GetUserOvertimes(userID uint) ([]models.Overtime, *utils.AppError) {
	var overtimes []models.Overtime
	if err := config.DB.Preload("User").Where("user_id = ?", userID).Order("date desc, start_time desc").Find(&overtimes).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal mengambil data lembur")
	}
	return overtimes, nil
}

func GetOvertimeByID(id string) (*models.Overtime, *utils.AppError) {
	var overtime models.Overtime
	if err := config.DB.Preload("User").First(&overtime, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Data lembur tidak ditemukan")
	}
	return &overtime, nil
}

func UpdateOvertime(id string, title, description string, date, startTime, endTime time.Time, userID uint) (*models.Overtime, *utils.AppError) {
	var overtime models.Overtime
	if err := config.DB.First(&overtime, id).Error; err != nil {
		return nil, utils.NewAppError(http.StatusNotFound, "Data lembur tidak ditemukan")
	}

	durationMinutes := endTime.Sub(startTime).Minutes()
	if durationMinutes < 0 {
		return nil, utils.NewAppError(http.StatusBadRequest, "Waktu selesai harus lebih besar dari waktu mulai")
	}
	price := (durationMinutes / 60.0) * 5000.0

	overtime.Title = title
	overtime.Description = description
	overtime.Date = date
	overtime.StartTime = startTime
	overtime.EndTime = endTime
	overtime.Price = price
	overtime.UserID = userID

	if err := config.DB.Save(&overtime).Error; err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal memperbarui data lembur")
	}

	return &overtime, nil
}

func DeleteOvertime(id string) *utils.AppError {
	var overtime models.Overtime
	if err := config.DB.First(&overtime, id).Error; err != nil {
		return utils.NewAppError(http.StatusNotFound, "Data lembur tidak ditemukan")
	}

	if err := config.DB.Delete(&overtime).Error; err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "Gagal menghapus data lembur")
	}

	return nil
}
