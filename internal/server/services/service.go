package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
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

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) ApplySync(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error {
	return s.repo.ApplySync(ctx, in, userID)
}

func (s *Service) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	return s.repo.GetLatestUserRev(ctx, userID)
}

func (s *Service) SyncGet(ctx context.Context, userID int64, since int64, limit int) (*shrModel.SyncGet, error) {
	s.logger.Debug("Get changes...")

	changes, err := s.repo.GetUserSyncChanges(ctx, userID, since, limit)
	if err != nil {
		return nil, err
	}

	var itemIDs []int64

	for _, change := range changes {
		itemIDs = append(itemIDs, change.ItemID)
	}

	s.logger.Debug("Get Items...")

	items, err := s.repo.GetItemsByIDs(ctx, itemIDs)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Get payloads...")

	payloads, err := s.repo.GetPayloadByItemIDs(ctx, itemIDs)
	if err != nil {
		return nil, err
	}

	s.logger.Debug("Get metadata...")

	meta, err := s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)
	if err != nil {
		return nil, err
	}

	var nextRev int64

	// Если изменения есть, то след. ревизия - краяняя
	if len(changes) > 0 {
		nextRev = changes[len(changes)-1].Revision
	} else {
		// Если изменений для синхронизации нет, то следующая
		// ревизия - это переданная
		nextRev = since
	}

	// Есть ли еще записи
	hasMore := len(changes) == limit

	return s.mapSyncResponse(
		changes,
		items,
		payloads,
		meta,
		since,
		nextRev,
		hasMore,
	), nil
}
