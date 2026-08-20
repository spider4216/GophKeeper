package services

import (
	"net/http"

	"github.com/spider4216/GophKeeper/internal/client/repositories"
)

type Service struct {
	client *http.Client
	host   string
	repo   repositories.Repository
}

func New(client *http.Client, host string, repo repositories.Repository) *Service {
	return &Service{
		client: client,
		host:   host,
		repo:   repo,
	}
}
