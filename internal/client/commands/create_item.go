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
	title := fs.String("title", "", "Title")
	userID := fs.Int64("user-id", 0, "User ID")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *login == "" || *pass == "" || *title == "" || *userID == 0 {
		return "", errors.New("login, password, title and user-id are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	req := models.LoginPassReq{
		Login: *login,
		Pass:  *pass,
		Title: *title,
		JWT:   auth.Token,
	}

	if err := c.Service.CreateUserPassItem(ctx, req, c.Cfg.EncryptKey, auth.UserID); err != nil {
		return "", fmt.Errorf("cannot create userpass: %w", err)
	}

	return "Item successfully created", nil
}
