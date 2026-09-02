package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

func (c *Command) CreateCard(ctx context.Context, args []string) (string, error) {
	fs := flag.NewFlagSet(InsertLoginPass.String(), flag.ExitOnError)

	pan := fs.String("pan", "", "PAN")
	dateStr := fs.String("date", "", "Exp date")
	cvcStr := fs.String("cvc", "", "cvc")
	holder := fs.String("holder", "", "Holder")
	userID := fs.Int64("user-id", 0, "User ID")
	title := fs.String("title", "", "Title")

	// todo validate pan and cvc

	if err := fs.Parse(args); err != nil {
		return "", err
	}

	if *pan == "" || *dateStr == "" || *cvcStr == "" || *userID == 0 || *holder == "" || *title == "" {
		return "", errors.New("pan, date, cvc, holder, title and user-id are required")
	}

	auth, err := c.Service.GetToken(ctx, *userID)
	if err != nil {
		return "", err
	}

	parsedDate, err := time.Parse("02.01.2006", *dateStr)
	if err != nil {
		return "", fmt.Errorf("error parse date: %w", err)
	}

	cvc, err := strconv.Atoi(*cvcStr)
	if err != nil {
		return "", fmt.Errorf("cannot convert cvc to int: %w", err)
	}

	req := models.CardReq{
		Pan:    *pan,
		Cvc:    cvc,
		Date:   parsedDate,
		Holder: *holder,
		Title:  *title,
	}

	if err := c.Service.CreateCardItem(ctx, req, c.Cfg.EncryptKey, auth.UserID); err != nil {
		return "", fmt.Errorf("cannot create card item: %w", err)
	}

	return "Card successfully created", nil
}
