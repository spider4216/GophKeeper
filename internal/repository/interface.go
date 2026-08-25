package repository

import (
	"context"
	"database/sql"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

type CommonRepositoryInterface interface {
	GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]shrModel.MetadataRepo, error)
	DeleteUserItemByID(ctx context.Context, itemID int64, userID int64) error
	DeleteUserItemByIDTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64) error
	DeleteUserMetaByItemID(ctx context.Context, itemID int64, userID int64) error
	DeleteUserMetaByItemIDTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64) error
	CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error)
	CreateMetaTx(ctx context.Context, tx *sql.Tx, itemID int64, k string, v string) (int64, error)
}
