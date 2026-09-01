package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

// хранилище где данные складываются в БД PostgreSQL.
type SrvRepository struct {
	con        *sql.DB
	logger     *slog.Logger
	commonRepo commonRep.CommonRepositoryInterface
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *slog.Logger, common commonRep.CommonRepositoryInterface) (*SrvRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &SrvRepository{con: db, logger: logger, commonRepo: common}, nil
}

func (repo *SrvRepository) GetCommonRepo() commonRep.CommonRepositoryInterface {
	return repo.commonRepo
}

func (repo *SrvRepository) Source() any {
	return repo.con
}

func (repo *SrvRepository) CreateUser(ctx context.Context, user models.UserRepo) (int64, error) {
	sql := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"
	var lastInsertId int64

	err := repo.con.QueryRowContext(ctx, sql, user.Email, user.PasswordHash).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func (repo *SrvRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserRepo, error) {
	sql := "SELECT id,email,password, created_at FROM users WHERE email=$1"

	user := models.UserRepo{}

	err := repo.con.QueryRow(sql, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *SrvRepository) CreateItem(ctx context.Context, item models.ItemRepo) (string, error) {
	return repo.createItem(ctx, repo.con, item)
}

func (repo *SrvRepository) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (string, error) {
	return repo.createItem(ctx, tx, item)
}

func (repo *SrvRepository) createItem(ctx context.Context, db commonRep.Querier, item models.ItemRepo) (string, error) {
	sql := "INSERT INTO items (id, type ,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id string

	if err := db.QueryRowContext(ctx, sql, item.ID, item.Type, item.UserID).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (repo *SrvRepository) CreateItemPayloadTx(ctx context.Context, tx *sql.Tx, item models.ItemPayloadRepo) error {
	sql := "INSERT INTO item_payloads (item_id,ciphertext) VALUES ($1,$2)"

	_, err := tx.ExecContext(ctx, sql, item.ItemID, item.Ciphertext)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) CreateSyncChangesTx(ctx context.Context, tx *sql.Tx, itemID string, op enum.OpType, userID int64) (int64, error) {
	sql := "INSERT INTO sync_changes (item_id,operation,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := tx.QueryRowContext(ctx, sql, itemID, op, userID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	sql := "SELECT MAX(revision) FROM sync_changes WHERE user_id=$1 GROUP BY user_id"

	var rev int64

	err := repo.con.QueryRow(sql, userID).Scan(&rev)
	if err != nil {
		return 0, err
	}

	return rev, nil
}

func (repo *SrvRepository) GetUserSyncChanges(ctx context.Context, userID int64, since int64, limit int) iter.Seq[models.SyncChangesRepo] {
	return func(yield func(models.SyncChangesRepo) bool) {
		sql := "SELECT id, item_id, revision, operation, created_at FROM sync_changes WHERE user_id=$1 and revision > $2 ORDER BY revision LIMIT $3;"

		rows, err := repo.con.QueryContext(ctx, sql, userID, since, limit)
		if err != nil {
			repo.logger.Error("Cannot get sync rows", "error", err)
			return
		}

		defer func() {
			if err := rows.Close(); err != nil {
				repo.logger.Warn("Cannot close rows", "error", err)
			}
		}()

		for rows.Next() {
			var item models.SyncChangesRepo

			if err := rows.Scan(
				&item.ID,
				&item.ItemID,
				&item.Revision,
				&item.Operation,
				&item.CreatedAt,
			); err != nil {
				repo.logger.Error("cannot scan sync changes", "error", err)
				return
			}

			if !yield(item) {
				return
			}
		}

		if err := rows.Err(); err != nil {
			repo.logger.Error("row error in get sync changes", "error", err)
			return
		}
	}
}

func (repo *SrvRepository) GetItemsByIDs(ctx context.Context, itemIDs []string) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, created_at FROM items WHERE id = ANY($1);"

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

func (repo *SrvRepository) GetPayloadByItemIDs(ctx context.Context, itemIDs []string) ([]models.ItemPayloadRepo, error) {
	sql := "SELECT item_id, ciphertext FROM item_payloads WHERE item_id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warn("Cannot close rows", "error", err)
		}
	}()

	var items []models.ItemPayloadRepo

	for rows.Next() {
		var item models.ItemPayloadRepo

		if err := rows.Scan(
			&item.ItemID,
			&item.Ciphertext,
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

func (repo *SrvRepository) DeletePayloadByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string) error {
	sql := "DELETE FROM item_payloads WHERE item_id=$1"

	_, err := tx.ExecContext(ctx, sql, itemID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) UpdateUserItemPayloadTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64, val string) error {
	sql := "UPDATE item_payloads p SET ciphertext=$1 FROM items i WHERE p.item_id=$2 AND i.user_id=$3"

	_, err := tx.ExecContext(ctx, sql, val, itemID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) ApplySync(ctx context.Context, in []shrModel.SyncPutChange, userID int64) error {
	tx, err := repo.con.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			repo.logger.Warn("cannot rollback in apply sync", "error", err)
		}
	}()

	for _, item := range in {
		repo.logger.Debug("Sync Operation", "operation", item.Operation)

		line := models.ItemRepo{
			ID:     item.Item.ID,
			UserID: userID,
			Type:   item.Item.Type,
		}

		switch item.Operation {
		case enum.OpCreate:
			if err := repo.syncCreate(ctx, tx, line, item, userID); err != nil {
				return fmt.Errorf("error sync. Cannot create: %w", err)
			}

		case enum.OpDelete:
			if err := repo.syncDelete(ctx, tx, item.Item.ID, userID); err != nil {
				return fmt.Errorf("error sync. Cannot delete: %w", err)
			}

		case enum.OpUpdate:
			if err := repo.syncUpdate(ctx, tx, userID, item); err != nil {
				return fmt.Errorf("error sync. Cannot update: %w", err)
			}

		default:
			return fmt.Errorf("unknown operation: %s", item.Operation)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *SrvRepository) syncUpdate(ctx context.Context, tx *sql.Tx, userID int64, item shrModel.SyncPutChange) error {
	repo.logger.Debug("Update strategy")
	if err := repo.UpdateUserItemPayloadTx(ctx, tx, item.Item.ID, userID, item.Item.Ciphertext); err != nil {
		return fmt.Errorf("cannot update item payload: %w", err)
	}

	for k, v := range item.Metadata {
		if err := repo.GetCommonRepo().UpdateMetaByItemIDAndKeyTx(ctx, tx, item.Item.ID, k, v); err != nil {
			return fmt.Errorf("cannot update metadata: %w", err)
		}
	}

	if _, err := repo.CreateSyncChangesTx(ctx, tx, item.Item.ID, enum.OpUpdate, userID); err != nil {
		return fmt.Errorf("cannot create sync changes: %w", err)
	}

	return nil
}

func (repo *SrvRepository) syncDelete(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	repo.logger.Debug("Delete strategy")
	if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete metadata: %w", err)
	}

	if err := repo.DeletePayloadByItemIDTx(ctx, tx, itemID); err != nil {
		return fmt.Errorf("cannot delete payload: %w", err)
	}

	if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete user item: %w", err)
	}

	if _, err := repo.CreateSyncChangesTx(ctx, tx, itemID, enum.OpDelete, userID); err != nil {
		return fmt.Errorf("cannot delete sync changes: %w", err)
	}

	// todo delete chanks only when it need
	if err := repo.GetCommonRepo().DeleteChunksByItemIDTx(ctx, tx, itemID); err != nil {
		return fmt.Errorf("cannot delete chunks: %w", err)
	}

	return nil
}

func (repo *SrvRepository) syncCreate(ctx context.Context, tx *sql.Tx, line models.ItemRepo, item shrModel.SyncPutChange, userID int64) error {
	repo.logger.Debug("create strategy")
	itemID, err := repo.CreateItemTx(ctx, tx, line)
	if err != nil {
		return fmt.Errorf("cannot create item: %w", err)
	}

	pl := models.ItemPayloadRepo{
		ItemID:     itemID,
		Ciphertext: item.Item.Ciphertext,
	}

	if err := repo.CreateItemPayloadTx(ctx, tx, pl); err != nil {
		return fmt.Errorf("cannot create item payload: %w", err)
	}

	for k, v := range item.Metadata {
		if _, err := repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, k, v); err != nil {
			return fmt.Errorf("cannot create meta: %w", err)
		}
	}

	if _, err := repo.CreateSyncChangesTx(ctx, tx, itemID, item.Operation, userID); err != nil {
		return fmt.Errorf("cannot create sync changes: %w", err)
	}

	return nil
}
