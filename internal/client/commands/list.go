package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
)

func (c *Command) UserList(args []string) (string, error) {
	// todo  command name to constant
	fs := flag.NewFlagSet(List.String(), flag.ExitOnError)

	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *token == "" {
		return "", errors.New("jwt is required")
	}

	claims, err := c.getClaims(*token)

	if err != nil {
		return "", fmt.Errorf("cannot parse jwt: %s", err)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	items, err := c.Service.GetUserItemsWithMeta(ctx, claims.UserID)

	if err != nil {
		return "", fmt.Errorf("cannot get user items: %s", err)
	}

	var builder strings.Builder

	for _, item := range items {
		builder.WriteString(fmt.Sprintf("ID: %d\n", item.ID))
		builder.WriteString(fmt.Sprintf("Type: %s\n", item.Type))
		builder.WriteString("Meta:")

		for _, meta := range item.Metadata {
			builder.WriteString(fmt.Sprintf("%s: %s\n", meta.Key, meta.Value))
		}

		builder.WriteString("---------------\n")
	}

	return builder.String(), nil
}
