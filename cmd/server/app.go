package main

import (
	"errors"
	"fmt"

	"github.com/spider4216/GophKeeper/internal/logger"
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/storage"
	migSrv "github.com/spider4216/GophKeeper/migrations/server"
	"go.uber.org/zap"
)

type app struct {
	logger *zap.SugaredLogger
	store  storage.Storage
	cfg    *config.Config
}

func newApp() *app {
	return &app{}
}

func (app *app) Run() error {

	if err := app.initConfig(); err != nil {
		return err
	}

	if err := app.initLogger(); err != nil {
		return err
	}

	if err := app.initStore(); err != nil {
		return err
	}

	if err := app.initMigrations(); err != nil {
		return err
	}

	return nil
}

func (app *app) initConfig() error {
	var cfg *config.Config

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.DbDsn == "" {
		return errors.New("dsn didtn't passed")
	}

	app.cfg = cfg

	return nil
}

func (app *app) initStore() error {
	store, err := storage.NewPgx(app.cfg.DbDsn, app.logger)
	if err != nil {
		return err
	}

	app.store = store

	return nil
}

func (app *app) initMigrations() error {
	app.logger.Debug("Up migrations")

	st, ok := app.store.(*storage.PgxStorage)

	if !ok {
		return fmt.Errorf("cannot cast to pgx store type in init migration")
	}

	if err := migSrv.MigrateSrv(st.Con); err != nil {
		return err
	}

	app.logger.Debug("Migration done")

	return nil
}

func (app *app) initLogger() error {
	// todo move lvl to cfg
	logger, err := logger.InitZap("debug")
	if err != nil {
		return err
	}

	app.logger = logger

	return nil
}
