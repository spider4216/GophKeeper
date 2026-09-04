package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/spider4216/GophKeeper/internal/client/models"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

const (
	authURL string = "/auth/login"
	regURL  string = "/auth/register"
)

func (s *Service) Register(ctx context.Context, req shrModel.RegisterReq) (*shrModel.RegisterResp, error) {
	url, err := url.JoinPath(s.host, regURL)
	if err != nil {
		return nil, err
	}

	return postRequest[shrModel.RegisterReq, shrModel.RegisterResp](ctx, url, req, s.client, *s.logger, http.StatusCreated)
}

func (s *Service) Login(ctx context.Context, req shrModel.LoginReq) (*shrModel.LoginResp, error) {
	url, err := url.JoinPath(s.host, authURL)
	if err != nil {
		return nil, err
	}

	return postRequest[shrModel.LoginReq, shrModel.LoginResp](ctx, url, req, s.client, *s.logger, http.StatusOK)
}

func postRequest[TReq, TResp any](ctx context.Context, url string, req TReq, cli *http.Client, logger slog.Logger, expCode int) (*TResp, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := cli.Do(r)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Warn("Cannot close body in login method")
		}
	}()

	if resp.StatusCode != expCode {
		return nil, fmt.Errorf("response status is %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res TResp

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *Service) SaveUserToken(ctx context.Context, userID int64, token string) error {
	return s.repo.SaveUserToken(ctx, userID, token)
}

func (s *Service) GetToken(ctx context.Context, userID int64) (*models.Auth, error) {
	return s.repo.GetToken(ctx, userID)
}
