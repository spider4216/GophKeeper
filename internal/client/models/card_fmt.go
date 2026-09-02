package models

import "time"

type CardFmt struct {
	Pan    string    `json:"pan"`
	Cvc    string    `json:"cvc"`
	Date   time.Time `json:"date"`
	Holder string    `json:"holder"`
}
