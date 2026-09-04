package models

import (
	"time"

	"github.com/spider4216/GophKeeper/internal/enum"
)

type ItemRepo struct {
	ID         string
	Type       enum.SecretType
	Ciphertext string
	UserID     int64
	CreatedAt  time.Time
}
