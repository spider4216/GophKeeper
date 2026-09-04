package server

import (
	"database/sql"
	"embed"

	"github.com/spider4216/GophKeeper/migrations"
)

//go:embed *.sql
var FS embed.FS

func MigrateSrv(db *sql.DB) error {
	return migrations.Run(FS, db)
}
