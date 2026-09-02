package models

import "time"

type CardReq struct {
	Pan    string
	Cvc    string
	Date   time.Time
	Holder string
	Title  string
}
