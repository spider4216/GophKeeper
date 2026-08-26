package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/spider4216/GophKeeper/internal/logger"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	migSrv "github.com/spider4216/GophKeeper/migrations/server"
	"go.uber.org/zap"
)

type app struct {
	logger *zap.SugaredLogger
	repo   repositories.Repository
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

	if err := app.initRepo(); err != nil {
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

func (app *app) initRepo() error {
	common, err := commonRep.NewRepository(app.cfg.DbDsn, app.logger)

	if err != nil {
		return err
	}

	repo, err := repositories.NewRepository(app.cfg.DbDsn, app.logger, common)
	if err != nil {
		return err
	}

	app.repo = repo

	return nil
}

func (app *app) initMigrations() error {
	app.logger.Debug("Up migrations")

	repo, ok := app.repo.(*repositories.SrvRepository)

	if !ok {
		return fmt.Errorf("cannot cast to pgx store type in init migration")
	}

	src := repo.Source().(*sql.DB)

	if err := migSrv.MigrateSrv(src); err != nil {
		return err
	}

	app.logger.Debug("Migration done")

	return nil
}

func (app *app) initLogger() error {
	logger, err := logger.InitZap(app.cfg.LogLvl)
	if err != nil {
		return err
	}

	app.logger = logger

	return nil
}
