package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	"go.uber.org/zap"
)

type Service struct {
	repo   repositories.Repository
	logger *zap.SugaredLogger
}

func New(repo repositories.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) CreateUser(ctx context.Context, email string, plainPass string) (int64, error) {
	// todo посолить пароль
	b := sha256.Sum256([]byte(plainPass))
	hash := hex.EncodeToString(b[:])

	u := models.UserRepo{
		Email:        email,
		PasswordHash: hash,
	}

	return s.repo.CreateUser(ctx, u)
}
