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
