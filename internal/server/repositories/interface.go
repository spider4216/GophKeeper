package repositories

import (
	"context"
	"database/sql"
	"iter"

	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/repository"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any

	CreateUser(ctx context.Context, user models.UserRepo) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error)
	CreateItem(ctx context.Context, item models.ItemRepo) (string, error)
	CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (string, error)
	CreateItemPayloadTx(ctx context.Context, tx *sql.Tx, item models.ItemPayloadRepo) error
	CreateSyncChangesTx(ctx context.Context, tx *sql.Tx, itemID string, op enum.OpType, userID int64) (int64, error)
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserSyncChanges(ctx context.Context, userID int64, since int64, limit int) iter.Seq[models.SyncChangesRepo]
	GetItemsByIDs(ctx context.Context, itemIDs []string) ([]models.ItemRepo, error)
	GetPayloadByItemIDs(ctx context.Context, itemIDs []string) ([]models.ItemPayloadRepo, error)
	DeletePayloadByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string) error
	UpdateUserItemPayloadTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64, val string) error
	GetCommonRepo() repository.CommonRepositoryInterface
	ApplySync(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error
}
