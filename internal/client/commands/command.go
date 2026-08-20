package commands

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/services"
)

type Command struct {
	Service *services.Service
	Cfg     *config.Config
}

// todo подумать над middleware
type claims struct {
	jwt.RegisteredClaims
	UserID int64
}

func New(service *services.Service, cfg *config.Config) *Command {
	return &Command{
		Service: service,
		Cfg:     cfg,
	}
}
