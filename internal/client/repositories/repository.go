package repositories

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/client/models"
)

// PgxStorage хранилище где данные складываются в БД PostgreSQL.
type ClientRepository struct {
	con    *sql.DB
	logger *zap.SugaredLogger
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *zap.SugaredLogger) (*ClientRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &ClientRepository{con: db, logger: logger}, nil
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

func (repo *ClientRepository) CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error) {
	sql := "INSERT INTO metadata (item_id,key,value) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, itemID, k, v).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *ClientRepository) CreatePendingChange(ctx context.Context, itemID int64, op string) error {
	sql := "INSERT INTO pending_changes (item_id,operation) VALUES ($1,$2)"

	_, err := repo.con.ExecContext(ctx, sql, itemID, op)

	if err != nil {
		return err
	}

	return nil
}

func (repo *ClientRepository) GetPendingUserChanges(ctx context.Context, userID int) ([]models.PendChangesRepo, error) {
	sql := "SELECT item_id,operation FROM pending_changes pc INNER JOIN items i ON i.id = pc.item_id WHERE i.user_id=$1"

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

func (repo *ClientRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]models.MetadataRepo, error) {
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

	var items []models.MetadataRepo

	for rows.Next() {
		var item models.MetadataRepo

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
