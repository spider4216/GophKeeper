package repositories

import (
	"context"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
	CreateItem(ctx context.Context, item models.ItemRepo) (int64, error)
	CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error)
	// todo op as custom type
	CreatePendingChange(ctx context.Context, itemID int64, op string) error

	GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error)
	GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error)
	GetUserItemByID(ctx context.Context, itemID int64, userID int64) (*models.ItemRepo, error)
	GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]models.MetadataRepo, error)
	DeletePendingByItemIDs(ctx context.Context, itemIDs []int64) error
	UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error
	CreateLastUserRev(ctx context.Context, userID int64, rev int64) error
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error)
}
