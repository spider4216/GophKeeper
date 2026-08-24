package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func (c *Command) SyncSend(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("sync-send", flag.ExitOnError)

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

	if err := c.Service.SyncSend(ctx, claims.UserID, *token); err != nil {
		fmt.Printf("cannot sync: %s", err)
		os.Exit(1)
	}

	fmt.Println("Sync was successfully done")
}
