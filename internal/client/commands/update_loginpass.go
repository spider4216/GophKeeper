package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) UpdateLoginPass(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(UpdateLoginPass.String(), flag.ExitOnError)

	itemID := fs.Int64("id", 0, "Item id")
	login := fs.String("login", "", "Login")
	metaID := fs.Int64("meta-id", 0, "Metadata ID")
	title := fs.String("title", "", "Title")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *login == "" || *pass == "" || *token == "" || *title == "" || *metaID == 0 {
		return "", errors.New("login, password, title, metadata id and jwt are required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %w", err)
	}

	if err := c.Service.UpdateLoginPass(ctx, *itemID, claims.UserID, *login, *pass, c.Cfg.EncryptKey, *title, *metaID); err != nil {
		return "", fmt.Errorf("cannot update: %w", err)
	}

	return "Item successfully updates", nil
}
