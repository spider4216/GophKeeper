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
}
