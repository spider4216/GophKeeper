package models

// todo можно сделать кастомный тип токена и статуса
// todo вынести в общие модели клиента и сервера, структура повторяется с сервером
type LoginResp struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
	ExpiredAt int64  `json:"expired_at,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
	UserID    int64  `json:"user_id"`
}
