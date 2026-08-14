package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	CustomerID    uint           `json:"customer_id" gorm:"not null"`
	Customer      *User          `json:"customer" gorm:"foreignKey:CustomerID"`
	WifiPackageID uint           `json:"wifi_package_id" gorm:"not null"`
	WifiPackage   *WifiPackage   `json:"wifi_package" gorm:"foreignKey:WifiPackageID"`
	PackagePrice  float64        `json:"package_price" gorm:"not null"`
	PPN           float64        `json:"ppn" gorm:"not null"`
	TotalAmount   float64        `json:"total_amount" gorm:"not null"`
	Status        string         `json:"status" gorm:"default:'paid'"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
