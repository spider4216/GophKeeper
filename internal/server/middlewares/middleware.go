package middlewares

import (
	"github.com/spider4216/GophKeeper/internal/server/config"
	"github.com/spider4216/GophKeeper/internal/server/services"
	"go.uber.org/zap"
)

type Middleware struct {
	logger  *zap.SugaredLogger
	cfg     *config.Config
	service *services.Service
}

func New(logger *zap.SugaredLogger, cfg *config.Config, service *services.Service) Middleware {
	return Middleware{
		logger:  logger,
		cfg:     cfg,
		service: service,
	}
}
