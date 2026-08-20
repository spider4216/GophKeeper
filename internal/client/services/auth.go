package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spider4216/GophKeeper/internal/client/models"
)

func (s *Service) Register(req models.RegisterReq) error {
	data, err := json.Marshal(req)

	if err != nil {
		return err
	}

	// todo endpoint to const
	url, err := url.JoinPath(s.host, "/auth/register")

	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(r)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("response status is %d", resp.StatusCode)
	}

	return nil
}
