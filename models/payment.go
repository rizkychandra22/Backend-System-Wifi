package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CustomerID    uint           `gorm:"not null" json:"customer_id"`
	Customer      *User          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	WifiServiceID uint           `gorm:"not null" json:"wifi_service_id"`
	WifiService   *WifiService   `gorm:"foreignKey:WifiServiceID" json:"wifi_service,omitempty"`
	PackagePrice  float64        `gorm:"not null" json:"package_price"`
	PPN           float64        `gorm:"not null" json:"ppn"`
	TotalAmount   float64        `gorm:"not null" json:"total_amount"`
	Status        string         `gorm:"type:varchar(20);default:'paid'" json:"status"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
