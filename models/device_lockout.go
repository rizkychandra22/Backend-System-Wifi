package models

import (
	"time"
)

type IPLockout struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IPAddress   string    `gorm:"type:varchar(45);uniqueIndex;not null" json:"ip_address"`
	LockedUntil time.Time `json:"locked_until"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
