package models

import "time"

type CardReq struct {
	Pan    string
	Cvc    int
	Date   time.Time
	Holder string
	Title  string
}
