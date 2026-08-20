package models

import "time"

type UserRepo struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
