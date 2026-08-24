package models

import (
	"time"

	"github.com/spider4216/GophKeeper/internal/enum"
)

type ItemRepo struct {
	ID         int64
	Type       enum.SecretType
	Ciphertext string
	UserID     int64
	CreatedAt  time.Time
}
