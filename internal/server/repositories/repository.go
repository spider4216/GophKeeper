package repositories

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

// хранилище где данные складываются в БД PostgreSQL.
type SrvRepository struct {
	con    *sql.DB
	logger *zap.SugaredLogger
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewRepository(dsn string, logger *zap.SugaredLogger) (*SrvRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &SrvRepository{con: db, logger: logger}, nil
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

func (repo *SrvRepository) CreateItem(ctx context.Context, item models.ItemRepo) (int64, error) {
	sql := "INSERT INTO items (id, type ,user_id) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, item.ID, item.Type, item.UserID).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) CreateItemPayload(ctx context.Context, item models.ItemPayloadRepo) error {
	sql := "INSERT INTO item_payloads (item_id,ciphertext) VALUES ($1,$2)"

	_, err := repo.con.ExecContext(ctx, sql, item.ItemID, item.Ciphertext)

	if err != nil {
		return err
	}

	return nil
}

func (repo *SrvRepository) CreateMeta(ctx context.Context, itemID int64, k string, v string) (int64, error) {
	sql := "INSERT INTO metadata (item_id,key,value) VALUES ($1,$2,$3) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, itemID, k, v).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) CreateSyncChanges(ctx context.Context, itemID int64, op string) (int64, error) {
	sql := "INSERT INTO sync_changes (item_id,operation) VALUES ($1,$2) RETURNING id"

	var id int64

	if err := repo.con.QueryRowContext(ctx, sql, itemID, op).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *SrvRepository) GetSyncChangesByID(ctx context.Context, ID int64) (*models.SyncChangesRepo, error) {
	sql := "SELECT id,item_id,revision,operation,created_at FROM sync_changes WHERE id=$1"

	var change models.SyncChangesRepo

	err := repo.con.QueryRowContext(ctx, sql, ID).Scan(&change.ID, &change.ID, &change.ItemID, &change.Revision, &change.Operation, &change.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &change, nil
}

func (repo *SrvRepository) GetLatestUserRev(ctx context.Context, userID int64) (int64, error) {
	sql := "SELECT MAX(sc.revision) FROM sync_changes sc INNER JOIN items i ON sc.item_id=i.id WHERE i.user_id=$1 GROUP BY i.user_id"

	var rev int64

	err := repo.con.QueryRow(sql, userID).Scan(&rev)

	if err != nil {
		return 0, err
	}

	return rev, nil
}

func (repo *SrvRepository) GetUserSyncChanges(ctx context.Context, userID int64, since int64) ([]models.SyncChangesRepo, error) {
	sql := "SELECT sc.id, sc.item_id, sc.revision, sc.operation, sc.created_at FROM sync_changes sc INNER JOIN items i ON i.id=sc.item_id WHERE i.user_id=$1 and sc.revision > $2;"

	rows, err := repo.con.QueryContext(ctx, sql, userID, since)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
		}
	}()

	var items []models.SyncChangesRepo

	for rows.Next() {
		var item models.SyncChangesRepo

		if err := rows.Scan(
			&item.ID,
			&item.ItemID,
			&item.Revision,
			&item.Operation,
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

func (repo *SrvRepository) GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]models.MetadataRepo, error) {
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

func (repo *SrvRepository) GetItemsByIDs(ctx context.Context, itemIDs []int64) ([]models.ItemRepo, error) {
	sql := "SELECT id, type, created_at FROM items WHERE id = ANY($1);"

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

func (repo *SrvRepository) GetPayloadByItemIDs(ctx context.Context, itemIDs []int64) ([]models.ItemPayloadRepo, error) {
	sql := "SELECT item_id, ciphertext FROM item_payloads WHERE item_id = ANY($1);"

	rows, err := repo.con.QueryContext(ctx, sql, itemIDs)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			repo.logger.Warnf("Cannot close rows: %s", err)
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
