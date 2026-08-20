package models

type UserReq struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}
