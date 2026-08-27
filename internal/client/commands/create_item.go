package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

func (c *Command) CreateLoginpass(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(InsertLoginPass.String(), flag.ExitOnError)

	login := fs.String("login", "", "Login")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")
	title := fs.String("title", "", "Title")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *login == "" || *pass == "" || *token == "" || *title == "" {
		return "", errors.New("login, password, title and jwt are required")
	}

	req := models.LoginPassReq{
		Login: *login,
		Pass:  *pass,
		Title: *title,
		JWT:   *token,
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %w", err)
	}

	if err := c.Service.CreateUserPassItem(ctx, req, c.Cfg.EncryptKey, claims.UserID); err != nil {
		return "", fmt.Errorf("cannot create userpass: %w", err)
	}

	return "Item successfully created", nil
}
