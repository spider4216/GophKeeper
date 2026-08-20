package models

type ItemRepo struct {
	ID int64
	// todo type as go type instead string
	Type       string
	Ciphertext string
	UserID     int64
}
