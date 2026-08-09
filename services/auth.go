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

func Login(phone string, password *string, clientIP string) (map[string]interface{}, *utils.AppError) {
	// 1. Cek tabel IPLockout
	var lockout models.IPLockout
	if err := config.DB.Where("ip_address = ?", clientIP).First(&lockout).Error; err == nil {
		if lockout.LockedUntil.After(time.Now()) {
			return nil, utils.NewAppError(http.StatusForbidden, "Anda diblokir karena mencoba mengakses halaman admin tanpa izin.")
		}
	}

	var user models.User
	if err := config.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "Nomor telepon tidak terdaftar")
	}

	if user.Role == "admin" {
		// Cek IP: apakah terikat dengan customer/employee
		var existingUser models.User
		if err := config.DB.Where("ip_address = ? AND role != ?", clientIP, "admin").First(&existingUser).Error; err == nil {
			// IP ini dipakai oleh customer/employee! Blokir IP 1 hari
			lockedUntil := time.Now().Add(24 * time.Hour)
			if lockout.ID != 0 {
				lockout.LockedUntil = lockedUntil
				config.DB.Save(&lockout)
			} else {
				newLockout := models.IPLockout{
					IPAddress:   clientIP,
					LockedUntil: lockedUntil,
				}
				config.DB.Create(&newLockout)
			}
			return nil, utils.NewAppError(http.StatusForbidden, "Anda sebagai "+string(existingUser.Role)+" tidak memiliki akses untuk ke halaman admin")
		}

		// Cek status terkunci
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			return nil, utils.NewAppError(http.StatusForbidden, "Akun terkunci karena terlalu banyak percobaan salah. Coba lagi dalam 30 menit.")
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
		// Device IP Lock Logic (Admin dibebaskan dari lock device)
		if user.IPAddress == nil || *user.IPAddress == "" {
			// First time login, save this IP
			user.IPAddress = &clientIP
			config.DB.Save(&user)
		} else if *user.IPAddress != clientIP {
			// IP mismatch, block login
			return nil, utils.NewAppError(http.StatusForbidden, "Akun ini sudah login di device lain. Silakan hubungi Admin.")
		}
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"role":  user.Role,
		"phone": user.Phone,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // 3 days expiration
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, utils.NewAppError(http.StatusInternalServerError, "JWT_SECRET not configured")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "Failed to generate token")
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
