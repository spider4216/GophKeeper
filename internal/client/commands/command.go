package commands

import (
	"context"
	"log/slog"

	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

type CmdName string

func (cn CmdName) String() string {
	return string(cn)
}

const (
	Register        CmdName = "register"
	Login           CmdName = "login"
	InsertLoginPass CmdName = "insert-loginpass"
	UpdateLoginPass CmdName = "update-loginpass"
	List            CmdName = "list"
	View            CmdName = "view"
	Delete          CmdName = "delete-item"
	SyncSend        CmdName = "sync-send"
	SyncGet         CmdName = "sync-get"
	Version         CmdName = "version"
	CreateText      CmdName = "create-text"
	GetText         CmdName = "get-text"
)

type CommandInterface interface {
	Register(ctx context.Context, args []string) (string, error)
	Login(ctx context.Context, args []string) (string, error)
	CreateLoginpass(ctx context.Context, args []string) (string, error)
	SyncSend(ctx context.Context, args []string) (string, error)
	UserList(ctx context.Context, args []string) (string, error)
	View(ctx context.Context, args []string) (string, error)
	SyncGet(ctx context.Context, args []string) (string, error)
	DeleteItem(ctx context.Context, args []string) (string, error)
	UpdateLoginPass(ctx context.Context, args []string) (string, error)
	PrintVersion() (string, error)
	CreateText(ctx context.Context, args []string) (string, error)
	GetText(ctx context.Context, args []string) (string, error)
}

type Command struct {
	Service *services.Service
	Cfg     *config.Config
	logger  *slog.Logger
}

func New(service *services.Service, cfg *config.Config, logger *slog.Logger) CommandInterface {
	return &Command{
		Service: service,
		Cfg:     cfg,
		logger:  logger,
	}
}
