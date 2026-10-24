package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title"`   // Sesuai dengan "Posisi"
	Company     string `gorm:"not null" json:"company"` // Sesuai dengan "Perusahaan"
	Link        string `json:"link"`                    // Sesuai dengan "Link"
	Type        string `json:"type"`                    // Sesuai dengan "Tipe" (misal: Full-time, Internship)
	Description string `json:"description"`             // Sesuai dengan "Catatan"
	// Default 'Terkirim' agar konsisten dengan daftar ValidStatuses di handlers
	Status    string         `gorm:"default:'Terkirim'" json:"status"`
	UserID    uint           `gorm:"index" json:"user_id"` // Index untuk mempercepat filter per user
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// StatusHistory: nyimpen riwayat tiap kali status job berubah
type StatusHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JobID     uint      `gorm:"not null;index" json:"job_id"` // Index karena sering di-query per job
	Status    string    `gorm:"not null" json:"status"`
	ChangedAt time.Time `json:"changed_at"`
}
