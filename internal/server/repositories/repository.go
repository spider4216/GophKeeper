package repositories

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

// PgxStorage хранилище где данные складываются в БД PostgreSQL.
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
