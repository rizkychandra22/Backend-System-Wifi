package models

import (
	"time"
)

type DeviceLockout struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	DeviceID    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"device_id"`
	LockedUntil time.Time `json:"locked_until"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
