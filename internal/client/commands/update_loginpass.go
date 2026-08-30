package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) UpdateLoginPass(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(UpdateLoginPass.String(), flag.ExitOnError)

	itemID := fs.String("id", "", "Item id")
	login := fs.String("login", "", "Login")
	metaID := fs.Int64("meta-id", 0, "Metadata ID")
	title := fs.String("title", "", "Title")
	pass := fs.String("password", "", "Password")
	userID := fs.Int64("user-id", 0, "User ID from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *login == "" || *pass == "" || *userID == 0 || *title == "" || *metaID == 0 {
		return "", errors.New("login, password, title, metadata id and jwt are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	if err := c.Service.UpdateLoginPass(ctx, *itemID, auth.UserID, *login, *pass, c.Cfg.EncryptKey, *title, *metaID); err != nil {
		return "", fmt.Errorf("cannot update: %w", err)
	}

	return "Item successfully updates", nil
}
