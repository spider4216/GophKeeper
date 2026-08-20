package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

// PgxStorage хранилище где данные складываются в БД PostgreSQL.
type PgxStorage struct {
	Con    *sql.DB
	logger *zap.SugaredLogger
}

// NewPgxStorage создание хранилища с БД PostgreSQL.
func NewPgxStorage(con string, logger *zap.SugaredLogger) (*PgxStorage, error) {
	db, err := sql.Open("pgx", con)
	if err != nil {
		return nil, err
	}

	return &PgxStorage{Con: db, logger: logger}, nil
}

func (db *PgxStorage) Ping(ctx context.Context) error {
	return db.Con.PingContext(ctx)
}

func (db *PgxStorage) Source() any {
	return db.Con
}
