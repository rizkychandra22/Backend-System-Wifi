package models

import (
	"time"

	"gorm.io/gorm"
)

type WifiPackage struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Price     float64        `gorm:"not null" json:"price"` // e.g. 150000
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
