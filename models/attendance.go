package models

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceStatus string

const (
	StatusProses  AttendanceStatus = "Proses"
	StatusHadir   AttendanceStatus = "Hadir"
	StatusLibur   AttendanceStatus = "Libur"
	StatusIzin    AttendanceStatus = "Izin"
)

type Attendance struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	UserID      *uint            `json:"user_id"`
	Date        string           `gorm:"type:varchar(10);not null" json:"date"`
	ClockIn     *time.Time       `json:"clock_in"`
	ClockOut    *time.Time       `json:"clock_out"`
	Grade       string           `gorm:"type:varchar(50)" json:"grade"`
	ClockInLat  *float64         `json:"clock_in_lat"`
	ClockInLng  *float64         `json:"clock_in_lng"`
	ClockOutLat *float64         `json:"clock_out_lat"`
	ClockOutLng *float64         `json:"clock_out_lng"`
	Status      AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	Notes       *string          `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
	User        User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
