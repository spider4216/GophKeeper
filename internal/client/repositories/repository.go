package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"log/slog"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
)

// PgxStorage хранилище где данные складываются в БД PostgreSQL.
type ClientRepository struct {
	con       *sql.DB
	logger    *slog.Logger
	commonRep commonRep.CommonRepositoryInterface
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *slog.Logger, common commonRep.CommonRepositoryInterface) (*ClientRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &ClientRepository{con: db, logger: logger, commonRep: common}, nil
}

func (repo *ClientRepository) GetCommonRepo() commonRep.CommonRepositoryInterface {
	return repo.commonRep
}

func (repo *ClientRepository) Source() any {
	return repo.con
}

func (repo *ClientRepository) CreateItem(ctx context.Context, item models.ItemRepo) (string, error) {
	return repo.createItem(ctx, repo.con, item)
}

func (repo *ClientRepository) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (string, error) {
	return repo.createItem(ctx, tx, item)
}

func (repo *ClientRepository) createItem(ctx context.Context, db commonRep.Querier, item models.ItemRepo) (string, error) {
	sql := "INSERT INTO items (id, type,ciphertext,user_id) VALUES ($1,$2,$3,$4) RETURNING id"

	var id string

	if err := db.QueryRowContext(ctx, sql, item.ID, item.Type, item.Ciphertext, item.UserID).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (repo *ClientRepository) CreatePendingChangeTx(ctx context.Context, tx *sql.Tx, itemID string, op enum.OpType, userID int64) error {
	sql := "INSERT INTO pending_changes (item_id,operation,user_id) VALUES ($1,$2,$3)"

	_, err := tx.ExecContext(ctx, sql, itemID, op, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error) {
	sql := "SELECT item_id,operation,user_id FROM pending_changes WHERE user_id=$1"

	rows, err := repo.con.QueryContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []models.PendChangesRepo

	for rows.Next() {
		var item models.PendChangesRepo

		if err := rows.Scan(
			&item.ItemID,
			&item.Operation,
			&item.UserID,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ClientRepository) GetItemsByIDs(ctx context.Context, itemIDs []string) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, ciphertext, created_at FROM items WHERE id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []models.ItemRepo

	for rows.Next() {
		var item models.ItemRepo

		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Ciphertext,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ClientRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []string) ([]shrModel.MetadataRepo, error) {
	sql := "SELECT id, item_id, key, value FROM metadata WHERE item_id = ANY($1);"
	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []shrModel.MetadataRepo

	for rows.Next() {
		var item shrModel.MetadataRepo

		if err := rows.Scan(
			&item.ID,
			&item.ItemID,
			&item.Key,
			&item.Value,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ClientRepository) DeletePendingByItemIDs(ctx context.Context, itemIDs []string) error {
	return repo.deletePendingByItemIDs(ctx, repo.con, itemIDs)
}

func (repo *ClientRepository) DeletePendingByItemIDsTx(ctx context.Context, tx *sql.Tx, itemIDs []string) error {
	return repo.deletePendingByItemIDs(ctx, tx, itemIDs)
}

func (repo *ClientRepository) deletePendingByItemIDs(ctx context.Context, db commonRep.Querier, itemIDs []string) error {
	sql := "DELETE FROM pending_changes WHERE item_id = ANY($1)"

	_, err := db.ExecContext(ctx, sql, itemIDs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	return repo.updateLastUserRev(ctx, repo.con, userID, rev)
}

func (repo *ClientRepository) UpdateLastUserRevTx(ctx context.Context, tx *sql.Tx, userID int64, rev int64) error {
	return repo.updateLastUserRev(ctx, tx, userID, rev)
}

func (repo *ClientRepository) updateLastUserRev(ctx context.Context, db commonRep.Querier, userID int64, rev int64) error {
	sql := "UPDATE sync_state SET last_server_revision=$1 WHERE user_id=$2"

	_, err := db.ExecContext(ctx, sql, rev, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) CreateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	sql := "INSERT INTO sync_state (user_id, last_server_revision) VALUES ($1,$2)"

	_, err := repo.con.ExecContext(ctx, sql, userID, rev)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	sql := "SELECT last_server_revision FROM sync_state WHERE user_id = $1;"

	var rev int64

	err := repo.con.QueryRowContext(ctx, sql, userID).Scan(&rev)
	if err != nil {
		return 0, err
	}

	return rev, nil
}

func (repo *ClientRepository) GetMetadataByItemID(ctx context.Context, itemID string) ([]shrModel.MetadataRepo, error) {
	sql := "SELECT id, item_id, key, value FROM metadata WHERE item_id = $1;"
	rows, err := repo.con.QueryContext(ctx, sql, itemID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []shrModel.MetadataRepo

	for rows.Next() {
		var item shrModel.MetadataRepo

		if err := rows.Scan(
			&item.ID,
			&item.ItemID,
			&item.Key,
			&item.Value,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ClientRepository) GetUserItems(ctx context.Context, userID int64) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, ciphertext, user_id, created_at FROM items WHERE user_id=$1;"

	rows, err := repo.con.QueryContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []models.ItemRepo

	for rows.Next() {
		var item models.ItemRepo

		if err := rows.Scan(
			&item.ID,
			&item.Type,
			&item.Ciphertext,
			&item.UserID,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ClientRepository) GetUserItemByID(ctx context.Context, itemID string, userID int64) (*models.ItemRepo, error) {
	sql := "SELECT id, type, ciphertext, user_id, created_at FROM items WHERE id=$1 and user_id=$2;"

	var item models.ItemRepo

	err := repo.con.QueryRowContext(ctx, sql, itemID, userID).Scan(&item.ID, &item.Type, &item.Ciphertext, &item.UserID, &item.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *ClientRepository) UpdateUserItemTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64, val string) error {
	sql := "UPDATE items SET ciphertext=$1 WHERE id=$2 AND user_id=$3"

	_, err := tx.ExecContext(ctx, sql, val, itemID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) CreateUserPassItem(ctx context.Context, item models.ItemRepo, userID int64, title string) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in create userpass", "error", err)
		}
	}()

	itemID, err := repo.CreateItemTx(ctx, tx, item)
	if err != nil {
		return fmt.Errorf("cannot create item: %w", err)
	}

	if _, err := repo.commonRep.CreateMetaTx(ctx, tx, itemID, "Title", title); err != nil {
		return fmt.Errorf("cannot create meta: %w", err)
	}

	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, enum.OpCreate, userID); err != nil {
		return fmt.Errorf("cannot create penfing: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) DeleteUserItem(ctx context.Context, itemID string, userID int64) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in selete user item", "error", err)
		}
	}()

	if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete meta: %w", err)
	}

	if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete item: %w", err)
	}

	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, enum.OpDelete, userID); err != nil {
		return fmt.Errorf("cannot create pending: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) ApplySync(ctx context.Context, userID int64, res shrModel.SyncGet) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in apply sync", "error", err)
		}
	}()

	for _, change := range res.Changes {
		switch change.Operation {
		case enum.OpCreate:
			if err := repo.syncCreate(ctx, tx, change, userID); err != nil {
				return fmt.Errorf("sync error. Cannot create: %w", err)
			}

		case enum.OpDelete:
			if err := repo.syncDelete(ctx, tx, change.Item.ID, userID); err != nil {
				return fmt.Errorf("sync error. Cannot delete: %w", err)
			}

		case enum.OpUpdate:
			if err := repo.syncUpdate(ctx, tx, change, userID); err != nil {
				return fmt.Errorf("sync error. Cannot update: %w", err)
			}

		default:
			return fmt.Errorf("unknown operation: %s", change.Operation)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) syncUpdate(ctx context.Context, tx *sql.Tx, change shrModel.SyncGetChange, userID int64) error {
	repo.logger.Debug("Update strategy on client")

	if err := repo.UpdateUserItemTx(ctx, tx, change.Item.ID, userID, change.Item.Ciphertext); err != nil {
		return fmt.Errorf("cannot update item: %w", err)
	}

	for k, v := range change.Metadata {
		if err := repo.GetCommonRepo().UpdateMetaByItemIDAndKeyTx(ctx, tx, change.Item.ID, k, v); err != nil {
			return fmt.Errorf("cannot update metadata: %w", err)
		}
	}

	return nil
}

func (repo *ClientRepository) syncDelete(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete user metadata: %w", err)
	}

	if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete user item: %w", err)
	}

	return nil
}

func (repo *ClientRepository) syncCreate(ctx context.Context, tx *sql.Tx, change shrModel.SyncGetChange, userID int64) error {
	item := models.ItemRepo{
		ID:         change.Item.ID,
		Type:       enum.SecretType(change.Item.Type),
		Ciphertext: change.Item.Ciphertext,
		UserID:     userID,
	}

	itemID, err := repo.CreateItemTx(ctx, tx, item)
	if err != nil {
		return fmt.Errorf("cannot create item: %w", err)
	}

	for k, v := range change.Metadata {
		if _, err := repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, k, v); err != nil {
			return fmt.Errorf("cannot create metadata: %w", err)
		}
	}

	return nil
}

func (repo *ClientRepository) UpdateLoginPass(ctx context.Context, itemID string, userID int64, encrypted string, metaID int64, title string) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in update userpass", "error", err)
		}
	}()

	if err := repo.UpdateUserItemTx(ctx, tx, itemID, userID, encrypted); err != nil {
		return fmt.Errorf("cannot update user item: %w", err)
	}

	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, enum.OpUpdate, userID); err != nil {
		return fmt.Errorf("cannot create pending change: %w", err)
	}

	if err := repo.GetCommonRepo().UpdateMetaByIDTx(ctx, tx, metaID, userID, title); err != nil {
		return fmt.Errorf("cannot update metadata: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) CommitSyncChunkTx(ctx context.Context, ids []string, userID int64, lastRev int64) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in commit sync", "error", err)
		}
	}()

	// Удаляем Pending для чанка
	if err := repo.DeletePendingByItemIDsTx(ctx, tx, ids); err != nil {
		return fmt.Errorf("cannot delete pending changes for chunk: %w", err)
	}

	// Обновляем последнюю ревизию для чанка
	if err := repo.UpdateLastUserRevTx(ctx, tx, userID, lastRev); err != nil {
		return fmt.Errorf("cannot update latest revision after chunk: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) SaveUserToken(ctx context.Context, userID int64, token string) error {
	sql := "INSERT INTO auth (user_id, token) VALUES ($1,$2) ON CONFLICT (user_id) DO UPDATE SET token=EXCLUDED.token"

	_, err := repo.con.ExecContext(ctx, sql, userID, token)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) GetToken(ctx context.Context, userID int64) (*models.Auth, error) {
	sql := "SELECT id,user_id,token FROM auth WHERE user_id=$1;"

	var auth models.Auth

	err := repo.con.QueryRowContext(ctx, sql, userID).Scan(&auth.ID, &auth.UserID, &auth.Token)
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (repo *ClientRepository) CreateItemForHugeText(ctx context.Context, userID int64, title string) (string, error) {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in save huge text", "error", err)
		}
	}()

	item := models.ItemRepo{
		ID:     uuid.NewString(),
		Type:   enum.Text,
		UserID: userID,
	}

	// Создаем Item
	itemID, err := repo.CreateItemTx(ctx, tx, item)
	if err != nil {
		return "", err
	}

	// Создаем Metadata для свободного текста
	// todo errgroup
	_, err = repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, "title", title)
	if err != nil {
		return "", err
	}

	_, err = repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, "encoding", "UTF-8")
	if err != nil {
		return "", err
	}

	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, enum.OpCreate, userID); err != nil {
		return "", fmt.Errorf("cannot create pending changes: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return itemID, nil
}

func (repo *ClientRepository) GetTextHugeData(ctx context.Context, itemID string) iter.Seq[shrModel.ChunkRepo] {
	return func(yield func(shrModel.ChunkRepo) bool) {
		sql := "SELECT item_id, chunk_number, ciphertext FROM item_chunks WHERE item_id=$1 ORDER BY chunk_number ASC;"

		rows, err := repo.con.QueryContext(ctx, sql, itemID)
		if err != nil {
			repo.logger.Error("Cannot get text data", "error", err)
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				repo.logger.Warn("Cannot close rows", "error", err)
			}
		}()

		for rows.Next() {
			var item shrModel.ChunkRepo

			if err := rows.Scan(
				&item.ItemID,
				&item.ChunkNumber,
				&item.Ciphertext,
			); err != nil {
				repo.logger.Error("cannot scan text data", "error", err)
				return
			}

			if !yield(item) {
				return
			}
		}

		if err := rows.Err(); err != nil {
			repo.logger.Error("row error in get text data", "error", err)
			return
		}
	}
}
