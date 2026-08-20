package repositories

import (
	"context"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

type Repository interface {
	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
	CreateItem(ctx context.Context, item models.ItemRepo) (int64, error)
}
