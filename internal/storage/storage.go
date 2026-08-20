package storage

import (
	"context"

	"go.uber.org/zap"
)

type Storage interface {
	// Ping проверяет доступность хранилища.
	Ping(ctx context.Context) error

	// Source возвращает инкапсулированное хранилище (источник).
	Source() any
}

func NewPgx(dsn string, logger *zap.SugaredLogger) (Storage, error) {
	pgxStore, err := NewPgxStorage(dsn, logger)
	if err != nil {
		return nil, err
	}

	return pgxStore, nil
}
