package services

import "net/http"

type Service struct {
	client *http.Client
	host   string
}

func New(client *http.Client, host string) *Service {
	return &Service{
		client: client,
		host:   host,
	}
}
