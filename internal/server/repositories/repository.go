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
