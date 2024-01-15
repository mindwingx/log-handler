package models

import (
	"time"
)

type PrimaryLogs struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Uuid      string    `json:"uuid"`
	File      string    `json:"file"` // file name and file path
	State     string    `json:"state"`
}
