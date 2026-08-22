package repositories

import (
	"context"

	"github.com/spider4216/GophKeeper/internal/server/models"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any

	CreateUser(ctx context.Context, user models.UserRepo) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error)
	CreateItem(ctx context.Context, item models.ItemRepo) (int64, error)
	CreateItemPayload(ctx context.Context, item models.ItemPayloadRepo) error
	CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error)
	// todo op custom type
	CreateSyncChanges(ctx context.Context, itemID int64, op string, userID int64) (int64, error)
	GetSyncChangesByID(ctx context.Context, ID int64) (*models.SyncChangesRepo, error)
	GetLatestUserRev(ctx context.Context, userID int64) (int64, error)
	GetUserSyncChanges(ctx context.Context, userID int64, since int64) ([]models.SyncChangesRepo, error)
	GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]models.MetadataRepo, error)
	GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error)
	GetPayloadByItemIDs(ctx context.Context, itemIDs []int64) ([]models.ItemPayloadRepo, error)
	DeleteUserItemByID(ctx context.Context, itemID int64, userID int64) error
	DeleteUserMetaByItemID(ctx context.Context, itemID int64, userID int64) error
	DeletePayloadByItemID(ctx context.Context, itemID int64) error
}
