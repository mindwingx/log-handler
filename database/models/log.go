package models

import (
	"gorm.io/gorm"
	"time"
)

type PrimaryLogs struct {
	gorm.Model
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Uuid         string    `json:"uuid"`
	File         string    `json:"file"` // file name and file path
	ProcessState string    `json:"process_state"`
	//relation
	Details []LogDetails `json:"details" gorm:"foreignKey:PrimaryLogId"`
}
