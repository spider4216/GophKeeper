package commands

import (
	"context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/services"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"go.uber.org/zap"
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
}

type Command struct {
	Service *services.Service
	Cfg     *config.Config
	logger  *zap.SugaredLogger
}

func New(service *services.Service, cfg *config.Config, logger *zap.SugaredLogger) CommandInterface {
	return &Command{
		Service: service,
		Cfg:     cfg,
		logger:  logger,
	}
}

func (c *Command) getClaims(token string) (*shrModel.Claims, error) {
	claims := &shrModel.Claims{}
	_, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(c.Cfg.JWTKey), nil
		})

	if err != nil {
		return nil, err
	}

	return claims, nil
}
