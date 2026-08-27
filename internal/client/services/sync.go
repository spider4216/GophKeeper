package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/spider4216/GophKeeper/internal/client/models"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

const (
	syncURL       string = "/sync"
	syncChunkSize int    = 1
)

func (s *Service) SyncSend(ctx context.Context, userID int64, token string) error {
	s.logger.Debug("Get pendings...")

	pends, err := s.repo.GetPendingUserChanges(ctx, int(userID))
	if err != nil {
		return fmt.Errorf("cannot get pending changes: %w", err)
	}

	s.logger.Debugf("Pendings count: %d", len(pends))

	for start := 0; start < len(pends); start += syncChunkSize {
		end := start + syncChunkSize

		// тут защита от выхода за пределы
		if end > len(pends) {
			end = len(pends)
		}

		// Извлекаем чанк
		chunk := pends[start:end]

		// Отправляем чанк на сервер
		s.logger.Debugf("Sync chunk: %d-%d of %d", start+1, end, len(pends))

		// Получаем ревизию изменений
		lastRev, err := s.syncChunk(ctx, userID, token, chunk)

		if err != nil {
			return fmt.Errorf("cannot sync chunk %d-%d: %w", start+1, end, err)
		}

		// Удаляем Pending для чанка
		if err := s.repo.DeletePendingByItemIDs(ctx, s.pendingItemIDs(chunk)); err != nil {
			return fmt.Errorf("cannot delete pending changes for chunk %d-%d: %w", start+1, end, err)
		}

		// Обновляем последнюю ревизию
		if err := s.repo.UpdateLastUserRev(ctx, userID, lastRev); err != nil {
			return fmt.Errorf("cannot update latest revision after chunk %d-%d: %w", start+1, end, err)
		}
	}

	s.logger.Debug("Sync done")

	return nil
}

func (s *Service) pendingItemIDs(pends []models.PendChangesRepo) []int64 {
	ids := make([]int64, 0, len(pends))

	for _, pend := range pends {
		ids = append(ids, pend.ItemID)
	}

	return ids
}

func (s *Service) syncChunk(ctx context.Context, userID int64, token string, pends []models.PendChangesRepo) (int64, error) {
	itemIDs := s.pendingItemIDs(pends)

	items, err := s.repo.GetItemsByIDs(ctx, itemIDs)
	if err != nil {
		return 0, fmt.Errorf("cannot get items: %w", err)
	}

	s.logger.Debugf("User items sync: %d", len(items))

	meta, err := s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)
	if err != nil {
		return 0, fmt.Errorf("cannot get metadata: %w", err)
	}

	req := s.buildSyncRequest(pends, items, meta)

	data, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("cannot build sync request: %w", err)
	}

	url, err := url.JoinPath(s.host, syncURL)
	if err != nil {
		return 0, fmt.Errorf("cannot join sync url to host: %w", err)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))

	if err != nil {
		return 0, fmt.Errorf("cannot create sync request: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", token)

	s.logger.Debug("Send sync chunk...")
	s.logger.Debug(string(data))

	resp, err := s.client.Do(r)

	if err != nil {
		return 0, fmt.Errorf("cannot send sync request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warnf("error closing response body: %s", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("sync status is not OK: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, fmt.Errorf("cannot read sync response: %w", err)
	}

	var res shrModel.SyncPutResp

	if err := json.Unmarshal(body, &res); err != nil {
		return 0, fmt.Errorf("cannot decode sync response: %w", err)
	}

	s.logger.Debugf("Latest revision: %d, for user: %d", res.LastRev, userID)

	return res.LastRev, nil
}

func (s *Service) SyncGet(ctx context.Context, userID int64, token string) error {
	// Получаем последнюю версию синхронизации
	rev, err := s.GetLatestUserRev(ctx, userID)

	if err != nil {
		return err
	}

	url, err := url.JoinPath(s.host, syncURL)

	if err != nil {
		return err
	}

	revStr := strconv.FormatInt(rev, 10)

	url = url + "?since=" + revStr

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("cannot create request for sync operation: %w", err)
	}

	r.Header.Add("Authorization", token)

	s.logger.Debug("Get sync...")
	resp, err := s.client.Do(r)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status sync is not OK: %d", resp.StatusCode)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warnf("Error closing body: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var res shrModel.SyncGet

	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("cannot unmarshal: %w", err)
	}

	s.logger.Debugf("Changes: %d", len(res.Changes))

	if err := s.repo.ApplySync(ctx, userID, res); err != nil {
		return fmt.Errorf("cannot apply sync: %w", err)
	}

	return nil
}
