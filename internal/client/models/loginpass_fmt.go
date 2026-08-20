package models

// Формат хранения login_pass
type LoginPassFmt struct {
	Login string `json:"login"`
	Pass  string `json:"password"`
}
