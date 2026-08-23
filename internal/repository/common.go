package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

// хранилище где данные складываются в БД PostgreSQL.
type CommonRepository struct {
	con    *sql.DB
	logger *zap.SugaredLogger
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *zap.SugaredLogger) (*CommonRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &CommonRepository{con: db, logger: logger}, nil
}

func (repo *CommonRepository) Source() any {
	return repo.con
}

func (repo *CommonRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]shrModel.MetadataRepo, error) {
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

func (repo *CommonRepository) DeleteUserItemByID(ctx context.Context, itemID int64, userID int64) error {
	sql := "DELETE FROM items WHERE user_id=$1 and id=$2"

	_, err := repo.con.ExecContext(ctx, sql, userID, itemID)

	if err != nil {
		return err
	}

	return nil
}
