package models

import (
	"time"

	"gorm.io/gorm"
)

type Allowance struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Month       string         `json:"month" gorm:"type:varchar(7);not null"`
	TargetType  string         `json:"target_type" gorm:"type:varchar(20);not null;default:'personal'"`
	UserID      *uint          `json:"user_id"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedByID uint           `json:"created_by_id" gorm:"not null"`
	CreatedBy   *User          `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
