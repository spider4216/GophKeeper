package models

import "time"

type CardFmt struct {
	Pan    string    `json:"pan"`
	Cvc    int       `json:"cvc"`
	Date   time.Time `json:"date"`
	Holder string    `json:"holder"`
}
