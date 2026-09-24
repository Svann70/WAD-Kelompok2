package models

import (
	"time"
	"gorm.io/gorm"
)

type Job struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Company     string         `gorm:"not null" json:"company"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'applied'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}