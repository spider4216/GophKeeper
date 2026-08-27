package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (c *Command) Register(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(Register.String(), flag.ExitOnError)

	email := fs.String("email", "", "user email")
	pass := fs.String("password", "", "user password")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *email == "" || *pass == "" {
		return "", errors.New("email and password are required")
	}

	if ok := c.Service.ValidateStrongPassword(*pass); !ok {
		return "", errors.New("password not strong enough (need one upper, one lower, one digit and >= 8 characters)")
	}

	if ok := c.Service.ValidateEmailFormat(*email); !ok {
		return "", errors.New("incorrect email format")
	}

	req := shrModel.RegisterReq{
		Email: *email,
		Pass:  *pass,
	}

	resp, err := c.Service.Register(ctx, req)
	if err != nil {
		return "", fmt.Errorf("cannot create user: %s", err)
	}

	if err := c.Service.CreateLastUserRev(ctx, resp.UserID, 0); err != nil {
		return "", fmt.Errorf("cannot create latest revision for user: %s", err)
	}

	return "Register is OK", nil
}
