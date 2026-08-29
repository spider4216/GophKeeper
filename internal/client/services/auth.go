package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

const (
	authURL string = "/auth/login"
	regURL  string = "/auth/register"
)

func (s *Service) Register(ctx context.Context, req shrModel.RegisterReq) (*shrModel.RegisterResp, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(s.host, regURL)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(r)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Cannot close body in register method")
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("response status is %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res shrModel.RegisterResp

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *Service) Login(ctx context.Context, req shrModel.LoginReq) (*shrModel.LoginResp, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(s.host, authURL)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(r)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Cannot close body in login method")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status is %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res shrModel.LoginResp

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
