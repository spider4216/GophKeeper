package repositories

import (
	"database/sql"

	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
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
