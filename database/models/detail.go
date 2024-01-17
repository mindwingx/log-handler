package models

import (
	"gorm.io/gorm"
	"time"
)

type LogDetails struct {
	gorm.Model
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"default:null"`
	DeletedAt    time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
	ID           uint64    `json:"id" gorm:"primaryKey"`
	PrimaryLogId uint64    `json:"primary_log_id"`
	LogLevel     string    `json:"log_level"`
	Log          string    `json:"log"`
	// relations
	LogModel PrimaryLogs `json:"log_model" gorm:"foreignKey:PrimaryLogId"`
}
