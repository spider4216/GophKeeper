package repositories

import (
	"context"
	"database/sql"

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
	CreateItem(ctx context.Context, item models.ItemRepo) (int64, error)
	CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (int64, error)
	CreateItemPayloadTx(ctx context.Context, tx *sql.Tx, item models.ItemPayloadRepo) error
	// todo op custom type
	CreateSyncChangesTx(ctx context.Context, tx *sql.Tx, itemID int64, op enum.OpType, userID int64) (int64, error)
	GetSyncChangesByID(ctx context.Context, ID int64) (*models.SyncChangesRepo, error)
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserSyncChanges(ctx context.Context, userID int64, since int64) ([]models.SyncChangesRepo, error)
	GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error)
	GetPayloadByItemIDs(ctx context.Context, itemIDs []int64) ([]models.ItemPayloadRepo, error)
	DeletePayloadByItemIDTx(ctx context.Context, tx *sql.Tx, itemID int64) error
	UpdateUserItemPayloadTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64, val string) error
	GetCommonRepo() repository.CommonRepositoryInterface
	ApplySync(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error
}
