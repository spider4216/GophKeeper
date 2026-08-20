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
