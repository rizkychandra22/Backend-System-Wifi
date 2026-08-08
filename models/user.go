package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
	RoleCustomer Role = "customer"
)

type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone          string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Role           Role           `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`
	IPAddress      *string        `gorm:"type:varchar(45)" json:"ip_address"` // Pointer to allow null (belum pernah login)
	ProfilePicture *string        `gorm:"type:text" json:"profile_picture"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
