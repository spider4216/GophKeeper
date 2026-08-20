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
	CreateSyncChanges(ctx context.Context, itemID int64, op string) (int64, error)
}
