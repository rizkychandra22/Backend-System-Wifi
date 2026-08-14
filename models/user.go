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
	ID                  uint           `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"type:varchar(100);not null" json:"name"`
	Phone               string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone"`
	Role                Role           `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`
	Password            *string        `gorm:"type:varchar(255)" json:"-"`
	FailedLoginAttempts int            `gorm:"default:0" json:"-"`
	LockedUntil         *time.Time     `json:"locked_until,omitempty"`
	Address             *string        `gorm:"type:text" json:"address"`
	IPAddress           *string        `gorm:"type:varchar(45)" json:"ip_address"`
	RegisteredByID      *uint          `json:"registered_by_id,omitempty"`
	RegisteredBy        *User          `gorm:"foreignKey:RegisteredByID" json:"registered_by,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}
