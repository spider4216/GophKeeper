package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) SyncGet(args []string) (string, error) {
	// todo  command name to constant
	fs := flag.NewFlagSet(SyncGet.String(), flag.ExitOnError)

	token := fs.String("token", "", "Token")

	fs.Parse(args)

	if *token == "" {
		return "", errors.New("Token is required")
	}

	// move to middleware
	claims, err := c.getClaims(*token)

	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	if err := c.Service.SyncGet(ctx, claims.UserID, *token); err != nil {
		return "", fmt.Errorf("cannot sync: %s", err)
	}

	return "Sync was successfully done", nil
}
