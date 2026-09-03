package reptest

import (
	"context"
	"database/sql"
	"log/slog"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

type SliceRepository struct {
	data   []any
	logger *slog.Logger
}

func NewRepository(logger *slog.Logger) *SliceRepository {
	return &SliceRepository{
		data:   []any{},
		logger: logger,
	}
}

func (r *SliceRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []string) ([]shrModel.MetadataRepo, error) {
	return nil, nil
}
func (r *SliceRepository) DeleteUserItemByID(ctx context.Context, itemID string, userID int64) error {
	return nil
}
func (r *SliceRepository) DeleteUserItemByIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	return nil
}
func (r *SliceRepository) DeleteUserMetaByItemID(ctx context.Context, itemID string, userID int64) error {
	return nil
}
func (r *SliceRepository) DeleteUserMetaByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	return nil
}
func (r *SliceRepository) CreateMeta(ctx context.Context, itemID string, k string, v string) (int64, error) {
	return 0, nil
}
func (r *SliceRepository) CreateMetaTx(ctx context.Context, tx *sql.Tx, itemID string, k string, v string) (int64, error) {
	return 0, nil
}
func (r *SliceRepository) UpdateMetaByIDTx(ctx context.Context, tx *sql.Tx, id int64, userID int64, v string) error {
	return nil
}
func (r *SliceRepository) UpdateMetaByItemIDAndKeyTx(ctx context.Context, tx *sql.Tx, itemID string, key string, v string) error {
	return nil
}
func (r *SliceRepository) InsertItemChunk(ctx context.Context, itemID string, chunkNum int, ciphertext string) error {
	return nil
}
func (r *SliceRepository) GetItemChunks(ctx context.Context, itemID string) ([]shrModel.ChunkRepo, error) {
	return nil, nil
}
func (r *SliceRepository) GetItemChunk(ctx context.Context, itemID string, chunkNum int) (*shrModel.ChunkRepo, error) {
	return nil, nil
}
func (r *SliceRepository) DeleteChunksByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string) error {
	return nil
}
