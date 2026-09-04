package model

type RegisterReq struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}
