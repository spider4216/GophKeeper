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
}
