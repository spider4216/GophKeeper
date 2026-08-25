package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
)

// PgxStorage хранилище где данные складываются в БД PostgreSQL.
type ClientRepository struct {
	con       *sql.DB
	logger    *zap.SugaredLogger
	commonRep commonRep.CommonRepositoryInterface
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *zap.SugaredLogger, common commonRep.CommonRepositoryInterface) (*ClientRepository, error) {
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

func (repo *ClientRepository) CreateItem(ctx context.Context, item models.ItemRepo) (int64, error) {
	sql := "INSERT INTO items (type,ciphertext,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, item.Type, item.Ciphertext, item.UserID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ClientRepository) CreateItemTx(ctx context.Context, tx *sql.Tx, item models.ItemRepo) (int64, error) {
	sql := "INSERT INTO items (type,ciphertext,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := tx.QueryRowContext(ctx, sql, item.Type, item.Ciphertext, item.UserID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ClientRepository) CreatePendingChange(ctx context.Context, itemID int64, op string, userID int64) error {
	sql := "INSERT INTO pending_changes (item_id,operation,user_id) VALUES ($1,$2,$3)"

	_, err := repo.con.ExecContext(ctx, sql, itemID, op, userID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) CreatePendingChangeTx(ctx context.Context, tx *sql.Tx, itemID int64, op string, userID int64) error {
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
			repo.logger.Warnf("Cannot close rows: %s", err)
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

func (repo *ClientRepository) GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, ciphertext, created_at FROM items WHERE id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
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

func (repo *ClientRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]shrModel.MetadataRepo, error) {
	sql := "SELECT id, item_id, key, value FROM metadata WHERE item_id = ANY($1);"
	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
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

func (repo *ClientRepository) DeletePendingByItemIDs(ctx context.Context, itemIDs []int64) error {
	sql := "DELETE FROM pending_changes WHERE item_id = ANY($1)"

	_, err := repo.con.ExecContext(ctx, sql, itemIDs)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) UpdateLastUserRev(ctx context.Context, userID int64, rev int64) error {
	sql := "UPDATE sync_state SET last_server_revision=$1 WHERE user_id=$2"

	_, err := repo.con.ExecContext(ctx, sql, rev, userID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) UpdateLastUserRevTx(ctx context.Context, tx *sql.Tx, userID int64, rev int64) error {
	sql := "UPDATE sync_state SET last_server_revision=$1 WHERE user_id=$2"

	_, err := tx.ExecContext(ctx, sql, rev, userID)

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

func (repo *ClientRepository) GetMetadataByItemID(ctx context.Context, itemID int64) ([]shrModel.MetadataRepo, error) {
	sql := "SELECT id, item_id, key, value FROM metadata WHERE item_id = $1;"
	rows, err := repo.con.QueryContext(ctx, sql, itemID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
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
			repo.logger.Warnf("Cannot close rows: %s", err)
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

func (repo *ClientRepository) GetUserItemByID(ctx context.Context, itemID int64, userID int64) (*models.ItemRepo, error) {
	sql := "SELECT id, type, ciphertext, user_id, created_at FROM items WHERE id=$1 and user_id=$2;"

	var item models.ItemRepo

	err := repo.con.QueryRowContext(ctx, sql, itemID, userID).Scan(&item.ID, &item.Type, &item.Ciphertext, &item.UserID, &item.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (repo *ClientRepository) UpdateUserItem(ctx context.Context, itemID int64, userID int64, val string) error {
	sql := "UPDATE items SET ciphertext=$1 WHERE id=$2 AND user_id=$3"

	_, err := repo.con.ExecContext(ctx, sql, val, itemID, userID)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) UpdateUserItemTx(ctx context.Context, tx *sql.Tx, itemID int64, userID int64, val string) error {
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
		// todo in errorf everywhere as %w
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	itemID, err := repo.CreateItemTx(ctx, tx, item)

	if err != nil {
		return fmt.Errorf("cannot create item: %w", err)
	}

	if _, err := repo.commonRep.CreateMetaTx(ctx, tx, itemID, "Title", title); err != nil {
		return fmt.Errorf("cannot create meta: %w", err)
	}

	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, "CREATE", userID); err != nil {
		return fmt.Errorf("cannot create penfing: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (repo *ClientRepository) DeleteUserItem(ctx context.Context, itemID int64, userID int64) error {
	tx, err := repo.con.BeginTx(ctx, nil)

	if err != nil {
		// todo in errorf everywhere as %w
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete meta: %w", err)
	}

	if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, itemID, userID); err != nil {
		return fmt.Errorf("cannot delete item: %w", err)
	}

	// todo op to const
	if err := repo.CreatePendingChangeTx(ctx, tx, itemID, "DELETE", userID); err != nil {
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
		// todo error
		_ = tx.Rollback()
	}()

	for _, change := range res.Changes {
		switch change.Operation {
		case "CREATE":
			// todo to method
			item := models.ItemRepo{
				Type:       enum.SecretType(change.Item.Type),
				Ciphertext: change.Item.Ciphertext,
				UserID:     userID,
			}

			itemID, err := repo.CreateItemTx(ctx, tx, item)

			if err != nil {
				return err
			}

			for k, v := range change.Metadata {
				if _, err := repo.GetCommonRepo().CreateMetaTx(ctx, tx, itemID, k, v); err != nil {
					return err
				}
			}

		case "DELETE":
			// todo to method
			if err := repo.GetCommonRepo().DeleteUserMetaByItemIDTx(ctx, tx, change.Item.ID, userID); err != nil {
				return err
			}

			if err := repo.GetCommonRepo().DeleteUserItemByIDTx(ctx, tx, change.Item.ID, userID); err != nil {
				return err
			}

		case "UPDATE":
			if err := repo.UpdateUserItemTx(ctx, tx, change.Item.ID, userID, change.Item.Ciphertext); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown operation: %s", change.Operation)
		}
	}

	if err := repo.UpdateLastUserRevTx(ctx, tx, userID, res.NextRev); err != nil {
		return fmt.Errorf("update last revision: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
