package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"net/url"

	"github.com/spider4216/GophKeeper/internal/client/models"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

const (
	syncURL string = "/sync"
)

func (s *Service) SyncSend(ctx context.Context, userID int64, token string, syncChunkSize int) error {
	s.logger.Debug("Get pendings...")

	pends, err := s.repo.GetPendingUserChanges(ctx, int(userID))
	if err != nil {
		return fmt.Errorf("cannot get pending changes: %w", err)
	}

	s.logger.Debug("Pendings count", "count", len(pends))

	for start := 0; start < len(pends); start += syncChunkSize {
		end := start + syncChunkSize

		// тут защита от выхода за пределы
		if end > len(pends) {
			end = len(pends)
		}

		// Извлекаем чанк
		chunk := pends[start:end]

		// Отправляем чанк на сервер
		s.logger.Debug("Sync chunk", "start", start+1, "end", end, "len", len(pends))

		// Получаем ревизию изменений
		lastRev, err := s.syncChunk(ctx, userID, token, chunk)
		if err != nil {
			return fmt.Errorf("cannot sync chunk %d-%d: %w", start+1, end, err)
		}

		// Коммитим чанк (удляеем pending, обновляем ревизию на клиенте)
		if err := s.repo.CommitSyncChunkTx(ctx, s.pendingItemIDs(chunk), userID, lastRev); err != nil {
			return fmt.Errorf("cannot commit sync: %w", err)
		}
	}

	s.logger.Debug("Sync done")

	return nil
}

func (s *Service) pendingItemIDs(pends []models.PendChangesRepo) []string {
	ids := make([]string, 0, len(pends))

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

	s.logger.Debug("User items sync", "len items", len(items))

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
			s.logger.Warn("error closing response body", "error", err)
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

	s.logger.Debug("Latest revision for user", "rev", res.LastRev, "user", userID)

	return res.LastRev, nil
}

func (s *Service) syncGetPage(ctx context.Context, token string, syncLimit int, rev int64) iter.Seq[shrModel.SyncGet] {
	return func(yield func(shrModel.SyncGet) bool) {
		lastRev := rev

		for {
			s.logger.Debug("Get sync changes since revision", "rev", lastRev)

			url, err := url.JoinPath(s.host, syncURL)
			if err != nil {
				s.logger.Error("cannot join sync url", "error", err)
				return
			}

			url = fmt.Sprintf("%s?since=%d&limit=%d", url, lastRev, syncLimit)

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				s.logger.Error("cannot create request for sync operation", "error", err)
				return
			}

			r.Header.Set("Authorization", token)

			s.logger.Debug("Get sync...")

			resp, err := s.client.Do(r)
			if err != nil {
				s.logger.Error("cannot do request", "error", err)
				return
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					s.logger.Warn("error closing body", "error", err)
				}
			}()

			if resp.StatusCode != http.StatusOK {
				s.logger.Error("status sync is not OK", "status", resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Error("cannot read sync response", "error", err)
				return
			}

			s.logger.Debug("Response from server", "content", string(body))

			var res shrModel.SyncGet

			if err := json.Unmarshal(body, &res); err != nil {
				s.logger.Error("cannot unmarshal sync response", "error", err)
				return
			}

			s.logger.Debug("Changes", "count", len(res.Changes))

			lastRev = res.NextRev

			if !yield(res) {
				return
			}

			if !res.HasMore {
				return
			}
		}
	}
}

func (s *Service) SyncGet(ctx context.Context, userID int64, token string, syncLimit int) error {
	// Получаем последнюю версию синхронизации
	rev, err := s.GetLatestUserRev(ctx, userID)
	if err != nil {
		return err
	}

	for page := range s.syncGetPage(ctx, token, syncLimit, rev) {
		if err := s.repo.ApplySync(ctx, userID, page); err != nil {
			s.logger.Error("cannot apply sync", "error", err)
			return err
		}
	}

	return nil
}
