package models

import "time"

type ItemRepo struct {
	ID        string
	UserID    int64
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
