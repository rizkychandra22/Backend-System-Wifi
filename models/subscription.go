package models

import (
	"time"
)

type Subscription struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	CustomerID    uint         `json:"customer_id" gorm:"uniqueIndex;not null"`
	Customer      *User        `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
	WifiPackageID uint         `json:"wifi_package_id" gorm:"not null"`
	WifiPackage   *WifiPackage `json:"wifi_package" gorm:"foreignKey:WifiPackageID"`
	BillingDay    int          `json:"billing_day" gorm:"not null"` // Hari jatuh tempo (1 - 31)
	NextDueDate   time.Time    `json:"next_due_date" gorm:"not null"` // Tanggal jatuh tempo berikutnya
	Status        string       `json:"status" gorm:"type:varchar(20);default:'active'"` // 'active', 'suspended', 'cancelled'
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
