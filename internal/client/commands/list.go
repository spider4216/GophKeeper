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

	token := fs.String("token", "", "JWT from server")

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *token == "" {
		return "", errors.New("jwt is required")
	}

	claims, err := c.getClaims(*token)
	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %w", err)
	}

	items, err := c.Service.GetUserItemsWithMeta(ctx, claims.UserID)
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
