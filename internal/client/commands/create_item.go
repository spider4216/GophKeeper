package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

func (c *Command) CreateLoginpass(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(InsertLoginPass.String(), flag.ExitOnError)

	login := fs.String("login", "", "Login")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")
	title := fs.String("title", "", "Title")

	fs.Parse(args)

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
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	if err := c.Service.CreateUserPassItem(ctx, req, c.Cfg.EncryptKey, claims.UserID); err != nil {
		return "", err
	}

	return "Item successfully created", nil
}
