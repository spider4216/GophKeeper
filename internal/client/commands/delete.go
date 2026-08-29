package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

func (c *Command) DeleteItem(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(Delete.String(), flag.ExitOnError)

	itemID := fs.String("id", "", "Item id")
	token := fs.String("token", "", "JWT from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *token == "" || *itemID == "" {
		return "", errors.New("login, password, title and jwt are required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %w", err)
	}

	if err != nil {
		return "", fmt.Errorf("cannot parse itemID: %w", err)
	}

	if err := c.Service.DeleteUserItem(ctx, *itemID, claims.UserID); err != nil {
		return "", fmt.Errorf("cannot delete item: %w", err)
	}

	return "Item successfully deleted", nil
}
