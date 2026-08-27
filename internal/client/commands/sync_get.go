package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) SyncGet(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(SyncGet.String(), flag.ExitOnError)

	token := fs.String("token", "", "Token")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *token == "" {
		return "", errors.New("token is required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	if err := c.Service.SyncGet(ctx, claims.UserID, *token); err != nil {
		return "", fmt.Errorf("cannot sync: %s", err)
	}

	return "Sync was successfully done", nil
}
