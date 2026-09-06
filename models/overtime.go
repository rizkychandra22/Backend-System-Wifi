package models

import (
	"time"

	"gorm.io/gorm"
)

type Overtime struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	Date        string         `json:"date" gorm:"type:varchar(10);not null"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Price       float64        `json:"price" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
