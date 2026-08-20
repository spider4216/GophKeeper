package client

import (
	"database/sql"
	"embed"

	"github.com/spider4216/GophKeeper/migrations"
)

//go:embed *.sql
var FS embed.FS

func MigrateClient(db *sql.DB) error {
	return migrations.Run(FS, db)
}
