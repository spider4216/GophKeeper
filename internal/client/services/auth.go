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

func (s *Service) Register(ctx context.Context, req shrModel.RegisterReq) (*shrModel.RegisterResp, error) {
	data, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	// todo endpoint to const
	url, err := url.JoinPath(s.host, "/auth/register")

	if err != nil {
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}

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

func (s *Service) Login(req shrModel.LoginReq) (*shrModel.LoginResp, error) {
	data, err := json.Marshal(req)

	if err != nil {
		return nil, err
	}

	// todo endpoint to const
	url, err := url.JoinPath(s.host, "/auth/login")

	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")

	// todo Может здесь нужен контекст в клиенте ???
	resp, err := s.client.Do(r)
	if err != nil {
		return nil, err
	}

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
