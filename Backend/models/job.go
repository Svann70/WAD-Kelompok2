package models

import (
	"time"
)

type Job struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Perusahaan string    `gorm:"type:varchar(255);not null" json:"perusahaan"`
	Posisi     string    `gorm:"type:varchar(255);not null" json:"posisi"`
	Tipe       string    `gorm:"type:varchar(100)" json:"tipe"`
	Tanggal    string    `gorm:"type:varchar(50)" json:"tanggal"`
	Catatan    string    `gorm:"type:text" json:"catatan"`
	UserID     uint      `gorm:"index" json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
