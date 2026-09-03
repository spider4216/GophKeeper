package reptest

import (
	"context"
	"database/sql"
	"encoding/json"
	"iter"
	"log/slog"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/repository"
	commonRep "github.com/spider4216/GophKeeper/internal/repository/reptest"
)

type SliceRepository struct {
	data      map[string][]string
	logger    *slog.Logger
	commonRep *commonRep.SliceRepository
}

func NewRepository(logger *slog.Logger, commonRep *commonRep.SliceRepository, store map[string][]string) *SliceRepository {
	return &SliceRepository{
		data:      store,
		logger:    logger,
		commonRep: commonRep,
	}
}

func (r *SliceRepository) Source() any {
	return nil
}
func (r *SliceRepository) CreateItem(ctx context.Context, item models.ItemRepo) (string, error) {
	_, err := json.Marshal(item)

	if err != nil {
		return "", err
	}

	return "", nil
}

func (r *SliceRepository) GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error) {
	return nil, nil
}
func (r *SliceRepository) GetItemsByIDs(ctx context.Context, itemIDs []string) ([]models.ItemRepo, error) {
	return nil, nil
}
func (r *SliceRepository) GetUserItemByID(ctx context.Context, itemID string, userID int64) (*models.ItemRepo, error) {
	return nil, nil
}
func (r *SliceRepository) DeletePendingByItemIDs(ctx context.Context, itemIDs []string) error {
	return nil
}
func (r *SliceRepository) DeletePendingByItemIDsTx(ctx context.Context, tx *sql.Tx, itemIDs []string) error {
	return nil
}
func (r *SliceRepository) UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return nil
}
func (r *SliceRepository) UpdateLastUserRevTx(ctx context.Context, tx *sql.Tx, userID int64, rev int64) error {
	return nil
}
func (r *SliceRepository) CreateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return nil
}
func (r *SliceRepository) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	return 0, nil
}
func (r *SliceRepository) GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error) {
	return nil, nil
}
func (r *SliceRepository) UpdateUserItemTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64, val string) error {
	return nil
}
func (r *SliceRepository) GetMetadataByItemID(ctx context.Context, itemID string) ([]shrModel.MetadataRepo, error) {
	return nil, nil
}
func (r *SliceRepository) GetCommonRepo() repository.CommonRepositoryInterface {
	return r.commonRep
}
func (r *SliceRepository) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (string, error) {
	data := map[string]any{
		"item_id":    item.ID,
		"user_id":    item.UserID,
		"type":       item.Type,
		"ciphertext": item.Ciphertext,
		"created_at": item.CreatedAt,
	}

	b, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	r.data["items"] = append(r.data["items"], string(b))

	return item.ID, nil
}
func (r *SliceRepository) CreatePendingChangeTx(ctx context.Context, tx *sql.Tx, itemID string, op enum.OpType, userID int64) error {
	data := map[string]any{
		"item_id":   itemID,
		"operation": op,
		"user_id":   userID,
	}

	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	r.data["pending_changes"] = append(r.data["pending_changes"], string(b))

	return nil
}
func (r *SliceRepository) CreateUserPassItem(ctx context.Context, item models.ItemRepo, userID int64, title string) error {
	itemID, err := r.CreateItemTx(ctx, nil, item)

	if err != nil {
		return err
	}

	_, err = r.commonRep.CreateMetaTx(ctx, nil, itemID, "Title", title)

	if err != nil {
		return err
	}

	return r.CreatePendingChangeTx(ctx, nil, itemID, enum.OpCreate, userID)
}
func (r *SliceRepository) DeleteUserItem(ctx context.Context, itemID string, userID int64) error {
	return nil
}
func (r *SliceRepository) ApplySync(ctx context.Context, userID int64, res shrModel.SyncGet) error {
	return nil
}
func (r *SliceRepository) UpdateLoginPass(ctx context.Context, itemID string, userID int64, encrypted string, metaID int64, title string) error {
	return nil
}
func (r *SliceRepository) CommitSyncChunkTx(ctx context.Context, ids []string, userID int64, lastRev int64) error {
	return nil
}
func (r *SliceRepository) SaveUserToken(ctx context.Context, userID int64, token string) error {
	return nil
}
func (r *SliceRepository) GetToken(ctx context.Context, userID int64) (*models.Auth, error) {
	return nil, nil
}
func (r *SliceRepository) GetBinaryData(ctx context.Context, itemID string) iter.Seq[shrModel.ChunkRepo] {
	return nil
}
func (r *SliceRepository) CreateItemForBinary(ctx context.Context, userID int64, meta []shrModel.MetadataRepo) (string, error) {
	return "", nil
}
