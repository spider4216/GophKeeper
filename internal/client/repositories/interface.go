package repositories

import (
	"context"
	"database/sql"

	"github.com/spider4216/GophKeeper/internal/client/models"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/repository"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
	CreateItem(ctx context.Context, item models.ItemRepo) (int64, error)
	// todo op as custom type
	CreatePendingChange(ctx context.Context, itemID int64, op string, userID int64) error

	GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error)
	GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error)
	GetUserItemByID(ctx context.Context, itemID int64, userID int64) (*models.ItemRepo, error)
	DeletePendingByItemIDs(ctx context.Context, itemIDs []int64) error
	UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error
	UpdateLastUserRevTx(ctx context.Context, tx *sql.Tx, userID int64, rev int64) error
	CreateLastUserRev(ctx context.Context, userID int64, rev int64) error
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error)
	UpdateUserItem(ctx context.Context, itemID int64, userID int64, val string) error
	UpdateUserItemTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64, val string) error
	GetMetadataByItemID(ctx context.Context, itemID int64) ([]shrModel.MetadataRepo, error)
	GetCommonRepo() repository.CommonRepositoryInterface
	CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (int64, error)
	CreatePendingChangeTx(ctx context.Context, tx *sql.Tx, itemID int64, op string, userID int64) error
	CreateUserPassItem(ctx context.Context, item models.ItemRepo, userID int64, title string) error
	DeleteUserItem(ctx context.Context, itemID int64, userID int64) error
	ApplySync(ctx context.Context, userID int64, res shrModel.SyncGet) error
	UpdateLoginPass(ctx context.Context, itemID int64, userID int64, encrypted string) error
}
