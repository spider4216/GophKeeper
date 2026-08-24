package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) UpdateLoginPass(args []string) (string, error) {
	fs := flag.NewFlagSet(UpdateLoginPass.String(), flag.ExitOnError)

	itemID := fs.Int64("id", 0, "Item id")
	login := fs.String("login", "", "Login")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")
	// todo подумать об обновлении meta

	fs.Parse(args)

	if *login == "" || *pass == "" || *token == "" {
		return "", errors.New("login, password, title and jwt are required")
	}

	claims, err := c.getClaims(*token)

	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	if err := c.Service.UpdateLoginPass(ctx, *itemID, claims.UserID, *login, *pass, c.Cfg.EncryptKey); err != nil {
		return "", fmt.Errorf("cannot update: %s", err)
	}

	return "Item successfully updates", nil
}
