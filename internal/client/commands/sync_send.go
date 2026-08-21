package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
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
	claims := &claims{}
	_, err := jwt.ParseWithClaims(*token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(c.Cfg.EncryptKey), nil
		})

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
