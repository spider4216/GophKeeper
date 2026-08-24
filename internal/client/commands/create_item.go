package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/spider4216/GophKeeper/internal/client/models"
	"github.com/spider4216/GophKeeper/internal/enum"
)

func (c *Command) CreateLoginpass(args []string) (string, error) {
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

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	// todo transaction
	// todo подумать что сделать с контекстом
	// todo тип в коснтанту
	itemID, err := c.Service.CreateLoginPassItem(ctx, enum.LoginPass, req, c.Cfg.EncryptKey, claims.UserID)

	if err != nil {
		return "", fmt.Errorf("cannot create item: %s", err)
	}

	_, err = c.Service.CreateMeta(ctx, itemID, "Title", req.Title)

	if err != nil {
		return "", fmt.Errorf("cannot create meta for item: %s", err)
	}

	// op to const and custom type
	err = c.Service.CreatePendingChange(ctx, itemID, "CREATE", claims.UserID)

	if err != nil {
		return "", fmt.Errorf("cannot create pending for item: %s", err)
	}

	return "Item successfully created", nil
}
