package commands

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/spider4216/GophKeeper/internal/client/config"
	"github.com/spider4216/GophKeeper/internal/client/services"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

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

func (c *Command) GetClaims(token string) (*shrModel.Claims, error) {
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
