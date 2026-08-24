package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func (c *Command) SyncGet(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("sync-get", flag.ExitOnError)

	token := fs.String("token", "", "Token")

	fs.Parse(args)

	if *token == "" {
		// todo instead fmt use something with out source
		fmt.Println("Token is required")
		os.Exit(1)
	}

	// move to middleware
	claims, err := c.GetClaims(*token)

	if err != nil {
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	if err := c.Service.SyncGet(ctx, claims.UserID, *token); err != nil {
		fmt.Printf("cannot sync: %s", err)
		os.Exit(1)
	}

	fmt.Println("Sync was successfully done")
}
