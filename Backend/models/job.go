package models

import (
	"time"
	"gorm.io/gorm"
)

type Job struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`       // Sesuai dengan "Posisi"
	Company     string         `gorm:"not null" json:"company"`     // Sesuai dengan "Perusahaan"
	Type        string         `json:"type"`                        // Sesuai dengan "Tipe" (misal: Full-time, Internship)
	Description string         `json:"description"`                 // Sesuai dengan "Catatan"
	Status      string         `gorm:"default:'applied'" json:"status"`
	UserID      uint           `json:"user_id"`                     // Sesuai dengan "UserID"
	CreatedAt   time.Time      `json:"created_at"`                  // Sesuai dengan "Tanggal"
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}