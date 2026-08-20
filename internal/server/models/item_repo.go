package models

import "time"

type ItemRepo struct {
	ID        int64
	UserID    int64
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
