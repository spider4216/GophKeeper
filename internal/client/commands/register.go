package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (c *Command) Register(args []string) (string, error) {
	// todo transaction
	fs := flag.NewFlagSet(Register.String(), flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	fs.Parse(args)

	if *email == "" || *pass == "" {
		return "", errors.New("email and password are required")
	}

	// todo подумать что то с контекстом
	ctx := context.Background()

	req := shrModel.RegisterReq{
		Email: *email,
		Pass:  *pass,
	}

	// todo condext for register user
	resp, err := c.Service.Register(req)

	if err != nil {
		return "", fmt.Errorf("cannot create user: %s", err)
	}

	// todo command logger

	if err := c.Service.CreateLastUserRev(ctx, resp.UserID, 0); err != nil {
		return "", fmt.Errorf("cannot create latest revision for user: %s", err)
	}

	return "Register is OK", nil
}
