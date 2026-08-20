package models

// todo скорее всего эти модели запросов и ответов уйдут в общие модели для клиента и сервера

type LoginReq struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}
