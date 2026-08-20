package models

// Формат запроса типа login_pass
type LoginPassReq struct {
	Login string
	Pass  string
	Title string
	JWT   string
}
