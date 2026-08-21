package commands

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func (c *Command) UserList(args []string) {
	// todo  command name to constant
	fs := flag.NewFlagSet("register", flag.ExitOnError)

	token := fs.String("token", "", "JWT from server")

	fs.Parse(args)

	if *token == "" {
		// todo instead fmt use something with out source
		fmt.Println("jwt is required")
		os.Exit(1)
	}

	claims := &claims{}
	_, err := jwt.ParseWithClaims(*token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(c.Cfg.JWTKey), nil
		})

	if err != nil {
		fmt.Printf("cannot parse jwt: %s", err)
		os.Exit(1)
	}

	// todo подумать как сдесь релизовать middleware
	ctx := context.WithValue(context.Background(), "userID", claims.UserID)

	items, err := c.Service.GetUserItemsWithMeta(ctx, claims.UserID)

	if err != nil {
		// todo everywhere \n
		fmt.Printf("cannot get user items: %s\n", err)
		os.Exit(1)
	}

	for _, item := range items {
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Type: %s\n", item.Type)
		fmt.Println("Meta:")
		for _, meta := range item.Metadata {
			fmt.Printf("%s: %s\n", meta.Key, meta.Value)
		}
		fmt.Println("---------------")
	}
}
