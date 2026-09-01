package repository

import (
	"context"
	"database/sql"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type CommonRepositoryInterface interface {
	GetMetadataByItemIDs(ctx context.Context, itemIDs []string) ([]shrModel.MetadataRepo, error)
	DeleteUserItemByID(ctx context.Context, itemID string, userID int64) error
	DeleteUserItemByIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error
	DeleteUserMetaByItemID(ctx context.Context, itemID string, userID int64) error
	DeleteUserMetaByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error
	CreateMeta(ctx context.Context, itemID string, k string, v string) (int64, error)
	CreateMetaTx(ctx context.Context, tx *sql.Tx, itemID string, k string, v string) (int64, error)
	UpdateMetaByIDTx(ctx context.Context, tx *sql.Tx, id int64, userID int64, v string) error
	UpdateMetaByItemIDAndKeyTx(ctx context.Context, tx *sql.Tx, itemID string, key string, v string) error
	InsertItemChunk(ctx context.Context, itemID string, chunkNum int, ciphertext string) error
	GetItemChunks(ctx context.Context, itemID string) ([]shrModel.ChunkRepo, error)
}
