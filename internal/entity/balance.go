package entity

import "time"

type UserBalance struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Current   float64   `json:"current" db:"current"`
	Withdrawn float64   `json:"withdrawn" db:"withdrawn"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
