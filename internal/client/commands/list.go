package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
)

func (c *Command) UserList(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(List.String(), flag.ExitOnError)

	userID := fs.Int64("user-id", 0, "User ID")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *userID == 0 {
		return "", errors.New("jwt is required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)

	if err != nil {
		return "", err
	}

	items, err := c.Service.GetUserItemsWithMeta(ctx, auth.UserID)
	if err != nil {
		return "", fmt.Errorf("cannot get user items: %w", err)
	}

	var builder strings.Builder

	for _, item := range items {
		fmt.Fprintf(&builder, "ID: %s\n", item.ID)
		fmt.Fprintf(&builder, "Type: %s\n", item.Type)
		fmt.Fprint(&builder, "Meta:\n")

		for _, meta := range item.Metadata {
			fmt.Fprintf(&builder, "%d) %s: %s\n", meta.ID, meta.Key, meta.Value)
		}

		builder.WriteString("---------------\n")
	}

	return builder.String(), nil
}
