package services

import (
	"backend-wifi/config"
	"backend-wifi/models"
	"backend-wifi/utils"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(phone string, password *string, deviceID string) (map[string]interface{}, *utils.AppError) {

	// Cek apakah device_id sedang di-lockout
	var lockout models.DeviceLockout
	if err := config.DB.Where("device_id = ?", deviceID).First(&lockout).Error; err == nil {
		if lockout.LockedUntil.After(time.Now()) {
			return nil, utils.NewAppError(http.StatusForbidden, "Akses login dari perangkat Anda diblokir selama 24 jam.")
		}
	}

	var user models.User
	if err := config.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "Nomor telepon tidak terdaftar")
	}

	if user.Role == "admin" {
		// Cek status terkunci akun admin
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			return nil, utils.NewAppError(http.StatusForbidden, "Akun terkunci karena terlalu banyak percobaan salah. Coba lagi dalam 30 menit.")
		}

		// Cek apakah deviceID ini pernah dipakai login oleh non-admin
		var nonAdminCount int64
		config.DB.Model(&models.User{}).Where("role IN ('employee', 'customer') AND device_id = ?", deviceID).Count(&nonAdminCount)
		if nonAdminCount > 0 {
			// Blokir device ini selama 24 jam
			lockedUntil := time.Now().Add(24 * time.Hour)
			if lockout.ID != 0 {
				lockout.LockedUntil = lockedUntil
				config.DB.Save(&lockout)
			} else {
				config.DB.Create(&models.DeviceLockout{
					DeviceID:    deviceID,
					LockedUntil: lockedUntil,
				})
			}
			return nil, utils.NewAppError(http.StatusForbidden, "Perangkat ini telah diblokir selama 24 jam karena pelanggaran keamanan percobaan akses admin.")
		}

		// Cek Password kosong
		if password == nil || *password == "" {
			return nil, utils.NewAppErrorWithData(http.StatusPreconditionRequired, "Password diperlukan untuk akun Admin", map[string]interface{}{"requires_password": true})
		}

		// Validasi Password
		if user.Password == nil || bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*password)) != nil {
			user.FailedLoginAttempts++
			if user.FailedLoginAttempts >= 5 {
				lockedTime := time.Now().Add(30 * time.Minute)
				user.LockedUntil = &lockedTime
				user.FailedLoginAttempts = 0
				config.DB.Save(&user)
				return nil, utils.NewAppError(http.StatusForbidden, "Terlalu banyak percobaan salah. Akun terkunci selama 30 menit.")
			}
			config.DB.Save(&user)
			return nil, utils.NewAppError(http.StatusUnauthorized, "Password salah")
		}

		// Benar password
		user.FailedLoginAttempts = 0
		user.LockedUntil = nil
		config.DB.Save(&user)

	} else {
		// Device Lock Logic (Admin dibebaskan dari lock device)
		if user.DeviceID == nil || *user.DeviceID == "" {
			// First time login, save this DeviceID
			user.DeviceID = &deviceID
			config.DB.Save(&user)
		} else if *user.DeviceID != deviceID {
			// Device mismatch, block login
			return nil, utils.NewAppError(http.StatusForbidden, "Akun ini sudah login di device lain. Silakan hubungi Admin.")
		}
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"role":  user.Role,
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 3 days expiration
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, utils.NewAppError(http.StatusInternalServerError, "JWT_SECRET belum dikonfigurasi")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Gagal membuat token")
	}

	return map[string]interface{}{
		"message": "Login berhasil",
		"token":   tokenString,
		"user": map[string]interface{}{
			"id":      user.ID,
			"name":    user.Name,
			"phone":   user.Phone,
			"role":    user.Role,
			"address": user.Address,
		},
	}, nil
}
