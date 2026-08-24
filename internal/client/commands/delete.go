package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strconv"
)

func (c *Command) DeleteItem(args []string) (string, error) {
	// todo  command name to constant
	fs := flag.NewFlagSet("delete-item", flag.ExitOnError)

	itemID := fs.String("id", "", "Item id")
	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *token == "" || *itemID == "" {
		return "", errors.New("login, password, title and jwt are required")
	}

	claims, err := c.getClaims(*token)

	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	itemIDInt, err := strconv.ParseInt(*itemID, 10, 64)

	if err != nil {
		return "", fmt.Errorf("cannot parse itemID: %s", err)
	}

	if err := c.Service.DeleteUserItem(ctx, itemIDInt, claims.UserID); err != nil {
		return "", fmt.Errorf("cannot delete item: %s", err)
	}

	return "Item successfully deleted", nil
}
