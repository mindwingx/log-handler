package models

import (
	"time"
)

type Clients struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at,omitempty" gorm:"default:null"`
	ID         uint64    `json:"id" gorm:"primaryKey"`
	NationalId string    `json:"national_id"`
	ClientUuid string    `json:"client_uuid"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Address    string    `json:"address,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Email      string    `json:"email"`
}
