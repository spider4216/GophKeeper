package commands

import "github.com/spider4216/GophKeeper/internal/client/services"

type Command struct {
	Service *services.Service
}

func New(service *services.Service) *Command {
	return &Command{
		Service: service,
	}
}
