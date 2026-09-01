package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/models"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
)

type syncChange struct {
	itemID    string
	operation enum.OpType
}

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

	changes := make([]syncChange, 0, limit)
	itemIDs := make([]string, 0, limit)

	var lastRevision int64

	for change := range s.repo.GetUserSyncChanges(ctx, userID, since, limit) {
		changes = append(changes, syncChange{
			itemID:    change.ItemID,
			operation: change.Operation,
		})

		itemIDs = append(itemIDs, change.ItemID)
		lastRevision = change.Revision
	}

	count := len(changes)
	nextRev := since

	if count > 0 {
		nextRev = lastRevision
	}

	s.logger.Debug("Get Items...")

	var items []models.ItemRepo
	var payloads []models.ItemPayloadRepo
	var meta []shrModel.MetadataRepo

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3)

	g.Go(func() error {
		s.logger.Debug("Get items...")
		var err error
		items, err = s.repo.GetItemsByIDs(ctx, itemIDs)

		return err
	})

	g.Go(func() error {
		s.logger.Debug("Get payloads...")
		var err error
		payloads, err = s.repo.GetPayloadByItemIDs(ctx, itemIDs)

		return err
	})

	g.Go(func() error {
		s.logger.Debug("Get metadata...")
		var err error
		meta, err = s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)

		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Есть ли еще записи
	hasMore := count == limit
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

func (s *Service) InsertItemChunk(ctx context.Context, itemID string, chunkNum int, ciphertext string) error {
	return s.repo.GetCommonRepo().InsertItemChunk(ctx, itemID, chunkNum, ciphertext)
}
