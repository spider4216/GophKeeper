package model

type LoginReq struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}
