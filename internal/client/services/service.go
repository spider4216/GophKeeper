package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/client/repositories"
	"go.uber.org/zap"
)

type Service struct {
	client *http.Client
	host   string
	repo   repositories.Repository
	logger *zap.SugaredLogger
}

func New(client *http.Client, host string, repo repositories.Repository, logger *zap.SugaredLogger) *Service {
	return &Service{
		client: client,
		host:   host,
		repo:   repo,
		logger: logger,
	}
}

// todo переименовать метода на CreateLoginPassItem
func (s *Service) CreateItem(ctx context.Context, t string, data models.LoginPassReq, key string, userID int64) (int64, error) {
	// todo t - as custom type

	// Формат хранения для типа
	d := models.LoginPassFmt{
		Login: data.Login,
		Pass:  data.Pass,
	}

	b, err := json.Marshal(d)

	if err != nil {
		return 0, err
	}

	encrypted, err := s.SignData(b, key)

	if err != nil {
		return 0, err
	}

	item := models.ItemRepo{
		Type:       t,
		Ciphertext: encrypted,
		UserID:     userID,
	}

	return s.repo.CreateItem(ctx, item)
}

func (s *Service) SignData(data []byte, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	if _, err := h.Write(data); err != nil {
		return "", err
	}

	sign := h.Sum(nil)

	return hex.EncodeToString(sign), nil
}

func (s *Service) CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error) {
	return s.repo.CreateMeta(ctx, itemID, k, v)
}

func (s *Service) CreatePendingChange(ctx context.Context, itemID int64, op string) error {
	return s.repo.CreatePendingChange(ctx, itemID, op)
}

func (s *Service) UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return s.repo.UpdateLastUserRev(ctx, userID, rev)
}

func (s *Service) CreateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return s.repo.CreateLastUserRev(ctx, userID, rev)
}

func (s *Service) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetLatestUserRev(ctx, userID)
}
