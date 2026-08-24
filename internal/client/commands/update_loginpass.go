package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func (c *Command) UpdateLoginPass(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("update-loginpass", flag.ExitOnError)

	itemID := fs.Int64("id", 0, "Item id")
	login := fs.String("login", "", "Login")
	pass := fs.String("password", "", "Password")
	token := fs.String("token", "", "JWT from server")
	// todo подумать об обновлении meta

	fs.Parse(args)

	if *login == "" || *pass == "" || *token == "" {
		// todo instead fmt use something with out source
		fmt.Println("login, password, title and jwt are required")
		os.Exit(1)
	}

	claims, err := c.GetClaims(*token)

	if err != nil {
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	if err := c.Service.UpdateLoginPass(ctx, *itemID, claims.UserID, *login, *pass, c.Cfg.EncryptKey); err != nil {
		fmt.Printf("cannot update: %s", err)
		os.Exit(1)
	}

	fmt.Println("Item successfully updates")
}
