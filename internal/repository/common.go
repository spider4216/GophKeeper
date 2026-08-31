package repository

import (
	"context"
	"database/sql"
	"log/slog"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

// хранилище где данные складываются в БД PostgreSQL.
type CommonRepository struct {
	con    *sql.DB
	logger *slog.Logger
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *slog.Logger) (*CommonRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &CommonRepository{con: db, logger: logger}, nil
}

func (repo *CommonRepository) Source() any {
	return repo.con
}

func (repo *CommonRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []string) ([]shrModel.MetadataRepo, error) {
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

func (repo *CommonRepository) DeleteUserItemByID(ctx context.Context, itemID string, userID int64) error {
	return repo.deleteUserItemByID(ctx, repo.con, itemID, userID)
}

func (repo *CommonRepository) DeleteUserItemByIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	return repo.deleteUserItemByID(ctx, tx, itemID, userID)
}

func (repo *CommonRepository) deleteUserItemByID(ctx context.Context, db Querier, itemID string, userID int64) error {
	sql := "DELETE FROM items WHERE user_id=$1 and id=$2"

	_, err := db.ExecContext(ctx, sql, userID, itemID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CommonRepository) DeleteUserMetaByItemID(ctx context.Context, itemID string, userID int64) error {
	return repo.deleteUserMetaByItemID(ctx, repo.con, itemID, userID)
}

func (repo *CommonRepository) DeleteUserMetaByItemIDTx(ctx context.Context, tx *sql.Tx, itemID string, userID int64) error {
	return repo.deleteUserMetaByItemID(ctx, tx, itemID, userID)
}

func (repo *CommonRepository) deleteUserMetaByItemID(ctx context.Context, db Querier, itemID string, userID int64) error {
	sql := "DELETE FROM metadata md USING items i WHERE i.id = md.item_id AND md.item_id=$1 AND i.user_id=$2"

	_, err := db.ExecContext(ctx, sql, itemID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *CommonRepository) CreateMeta(ctx context.Context, itemID string, k string, v string) (int64, error) {
	return repo.createMeta(ctx, repo.con, itemID, k, v)
}

func (repo *CommonRepository) CreateMetaTx(ctx context.Context, tx *sql.Tx, itemID string, k string, v string) (int64, error) {
	return repo.createMeta(ctx, tx, itemID, k, v)
}

func (repo *CommonRepository) createMeta(ctx context.Context, db Querier, itemID string, k string, v string) (int64, error) {
	sql := "INSERT INTO metadata (item_id,key,value) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := db.QueryRowContext(ctx, sql, itemID, k, v).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *CommonRepository) UpdateMetaByIDTx(ctx context.Context, tx *sql.Tx, id int64, userID int64, v string) error {
	sql := "UPDATE metadata md SET value=$1 FROM items i WHERE md.id=$2 AND i.user_id=$3"

	r, err := tx.ExecContext(ctx, sql, v, id, userID)
	if err != nil {
		return err
	}

	aff, err := r.RowsAffected()
	if err != nil {
		return err
	}

	repo.logger.Debug("Rows affected", "count", aff)

	return nil
}

func (repo *CommonRepository) UpdateMetaByItemIDAndKeyTx(ctx context.Context, tx *sql.Tx, itemID string, key string, v string) error {
	sql := "UPDATE metadata SET value=$1 WHERE item_id=$2 AND key=$3"

	_, err := tx.ExecContext(ctx, sql, v, itemID, key)
	if err != nil {
		return err
	}

	return nil
}
