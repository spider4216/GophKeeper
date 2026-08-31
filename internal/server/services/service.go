package services

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   repositories.Repository
	logger *slog.Logger
}

func New(repo repositories.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) CreateUser(ctx context.Context, email string, plainPass string) (int64, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	u := models.UserRepo{
		Email:        email,
		PasswordHash: string(b),
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

	changes := slices.Collect(
		s.repo.GetUserSyncChanges(ctx, userID, since, limit),
	)

	var itemIDs []string

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
		nextRev,
		hasMore,
	), nil
}

// IsErrAsDuplicate проверка, является ли переданная ошибка типум дубликата.
func (s *Service) IsErrAsDuplicate(err error) bool {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == pgerrcode.UniqueViolation
}
