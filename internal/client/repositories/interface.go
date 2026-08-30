package repositories

import (
	"context"
	"database/sql"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/repository"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
	CreateItem(ctx context.Context, item models.ItemRepo) (string, error)

	GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error)
	GetItemsByIDs(ctx context.Context, itemIDs []string) ([]models.ItemRepo, error)
	GetUserItemByID(ctx context.Context, itemID string, userID int64) (*models.ItemRepo, error)
	DeletePendingByItemIDs(ctx context.Context, itemIDs []string) error
	DeletePendingByItemIDsTx(ctx context.Context, tx *sql.Tx, itemIDs []string) error
	UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error
	UpdateLastUserRevTx(ctx context.Context, tx *sql.Tx, userID int64, rev int64) error
	CreateLastUserRev(ctx context.Context, userID int64, rev int64) error
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error)
	UpdateUserItemTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64, val string) error
	GetMetadataByItemID(ctx context.Context, itemID string) ([]shrModel.MetadataRepo, error)
	GetCommonRepo() repository.CommonRepositoryInterface
	CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (string, error)
	CreatePendingChangeTx(ctx context.Context, tx *sql.Tx, itemID string, op enum.OpType, userID int64) error
	CreateUserPassItem(ctx context.Context, item models.ItemRepo, userID int64, title string) error
	DeleteUserItem(ctx context.Context, itemID string, userID int64) error
	ApplySync(ctx context.Context, userID int64, res shrModel.SyncGet) error
	UpdateLoginPass(ctx context.Context, itemID string, userID int64, encrypted string, metaID int64, title string) error
	CommitSyncChunkTx(ctx context.Context, ids []string, userID int64, lastRev int64) error
	SaveUserToken(ctx context.Context, userID int64, token string) error
	GetToken(ctx context.Context, userID int64) (*models.Auth, error)
}
