package models

import "time"

type LogDetails struct {
	CreatedAt    time.Time `json:"created_at"`
	ID           uint64    `json:"id" gorm:"primaryKey"`
	PrimaryLogId uint64    `json:"primary_log_id"`
	LogLevel     string    `json:"log_level"`
	Log          string    `json:"log"`
}
