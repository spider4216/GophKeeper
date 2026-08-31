package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	shrConfig "github.com/spider4216/GophKeeper/config"
	"github.com/spider4216/GophKeeper/internal/logger"
	commonRep "github.com/spider4216/GophKeeper/internal/repository"
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/repositories"
	migSrv "github.com/spider4216/GophKeeper/migrations/server"
)

type app struct {
	logger *slog.Logger
	repo   repositories.Repository
	cfg    *config.Config
}

func newApp() *app {
	return &app{}
}

func (a *app) Run() error {
	_, err := shrConfig.NewBuilder(a).
		Step((*app).initConfig).
		Step((*app).initLogger).
		Step((*app).initRepo).
		Step((*app).initMigrations).
		Build()

	return err
}

func (a *app) initConfig() error {
	var cfg *config.Config

	cfg, err := config.New()
	if err != nil {
		return err
	}

	if cfg.DbDsn == "" {
		return errors.New("dsn didtn't passed")
	}

	a.cfg = cfg

	return nil
}

func (a *app) initRepo() error {
	common, err := commonRep.NewRepository(a.cfg.DbDsn, a.logger)
	if err != nil {
		return err
	}

	repo, err := repositories.NewRepository(a.cfg.DbDsn, a.logger, common)
	if err != nil {
		return err
	}

	a.repo = repo

	return nil
}

func (a *app) initMigrations() error {
	a.logger.Debug("Up migrations")

	repo, ok := a.repo.(*repositories.SrvRepository)

	if !ok {
		return fmt.Errorf("cannot cast to pgx store type in init migration")
	}

	src := repo.Source().(*sql.DB)

	if err := migSrv.MigrateSrv(src); err != nil {
		return err
	}

	a.logger.Debug("Migration done")

	return nil
}

func (a *app) initLogger() error {
	logger := logger.Init(a.cfg.LogLvl)

	a.logger = logger

	return nil
}
