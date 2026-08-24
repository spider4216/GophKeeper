package commands

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/services"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
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
	Register(args []string) (string, error)
	Login(args []string) (string, error)
	CreateLoginpass(args []string) (string, error)
	SyncSend(args []string) (string, error)
	UserList(args []string) (string, error)
	View(args []string) (string, error)
	SyncGet(args []string) (string, error)
	DeleteItem(args []string) (string, error)
	UpdateLoginPass(args []string) (string, error)
}

type Command struct {
	Service *services.Service
	Cfg     *config.Config
}

// todo подумать над middleware

func New(service *services.Service, cfg *config.Config) *Command {
	return &Command{
		Service: service,
		Cfg:     cfg,
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
