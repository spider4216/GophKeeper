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
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

const (
	syncURL  string = "/sync"
	chunkURL string = "/items/%s/chunks/%d"
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
		lastRev, err := s.syncChunkSend(ctx, userID, token, chunk)
		if err != nil {
			return fmt.Errorf("cannot sync chunk %d-%d: %w", start+1, end, err)
		}

		// Отправляем загрузку данных, если такая имеется
		if err := s.uploadChunk(ctx, chunk, token); err != nil {
			return fmt.Errorf("cannot upload chunk: %w", err)
		}

		// Коммитим чанк (удляеем pending, обновляем ревизию на клиенте)
		if err := s.repo.CommitSyncChunk(ctx, s.pendingItemIDs(chunk), userID, lastRev); err != nil {
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

func (s *Service) syncChunkSend(ctx context.Context, userID int64, token string, pends []models.PendChangesRepo) (int64, error) {
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

func (s *Service) uploadChunk(ctx context.Context, pends []models.PendChangesRepo, token string) error {
	itemIDs := s.pendingItemIDs(pends)

	items, err := s.repo.GetItemsByIDs(ctx, itemIDs)
	if err != nil {
		return err
	}

	// Оставляем только items с типом binary
	var filtered []models.ItemRepo

	for _, item := range items {
		if item.Type == enum.Binary {
			filtered = append(filtered, item)
		}
	}

	if len(filtered) <= 0 {
		s.logger.Debug("nothing upload for chunk")
		return nil
	}

	s.logger.Debug("uploads count", "count", len(items))

	for _, i := range filtered {
		// получаем чанки item
		chunks, err := s.repo.GetCommonRepo().GetItemChunks(ctx, i.ID)
		if err != nil {
			return err
		}

		// Отправляем чанки
		for num, chunk := range chunks {
			num = num + 1
			path := fmt.Sprintf(chunkURL, i.ID, num)

			url, err := url.JoinPath(s.host, path)
			if err != nil {
				return fmt.Errorf("cannot join sync url to host: %w", err)
			}

			req := models.UploadReq{
				Ciphertext: chunk.Ciphertext,
			}

			data, err := json.Marshal(req)
			if err != nil {
				return err
			}

			r, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
			if err != nil {
				return fmt.Errorf("cannot create upload request: %w", err)
			}

			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("Authorization", token)

			s.logger.Debug("Send upload chunk...", "num", num, "item", chunk.ItemID)

			resp, err := s.client.Do(r)
			if err != nil {
				return fmt.Errorf("cannot send upload request: %w", err)
			}

			if resp.StatusCode != http.StatusCreated {
				if err := resp.Body.Close(); err != nil {
					s.logger.Warn("error closing response body", "error", err)
				}

				return fmt.Errorf("upload status is not OK: %d", resp.StatusCode)
			}

			if err := resp.Body.Close(); err != nil {
				s.logger.Warn("error closing response body", "error", err)
			}
		}
	}

	return nil
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
		s.logger.Debug("Apply page")
		if err := s.repo.ApplySync(ctx, userID, page); err != nil {
			s.logger.Error("cannot apply sync", "error", err)
			return fmt.Errorf("cannot apply sync: %w", err)
		}

		s.logger.Debug("Download data for page")
		if err := s.downloadsData(ctx, token, page); err != nil {
			s.logger.Error("cannot downloads data", "error", err)
			return fmt.Errorf("cannot downloads data: %w", err)
		}

		s.logger.Debug("Commit revision for page")
		if err := s.repo.UpdateLastUserRev(ctx, userID, page.NextRev); err != nil {
			return fmt.Errorf("update last revision: %w", err)
		}
	}

	return nil
}

func (s *Service) downloadsData(ctx context.Context, token string, page shrModel.SyncGet) error {
	// Извлекаем идентификаторы страницы с Items которые имеют тип binary
	var itemIDs []string
	for _, change := range page.Changes {
		if change.Item.Type == enum.Binary.String() {
			itemIDs = append(itemIDs, change.Item.ID)
		}
	}

	if len(itemIDs) <= 0 {
		s.logger.Debug("Nothing download")
		return nil
	}

	// Для каждого item будем получать чанки с данными
	for _, itemID := range itemIDs {
		// Скачиваем чанки до тех пор пока чанков более не останется
		// для выбранного item (вернется 404)
		chunkNum := 1

		for {
			s.logger.Debug("Download chunk for iten", "chunk", chunkNum, "item", itemID)

			path := fmt.Sprintf("/items/%s/chunks/%d", itemID, chunkNum)

			url, err := url.JoinPath(s.host, path)
			if err != nil {
				return fmt.Errorf("cannot join in downloads data: %w", err)
			}

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("cannot create request for download data operation: %w", err)
			}

			r.Header.Set("Authorization", token)

			s.logger.Debug("Download...")

			resp, err := s.client.Do(r)
			if err != nil {
				return fmt.Errorf("cannot do request: %w", err)
			}

			defer func() {
				if err := resp.Body.Close(); err != nil {
					s.logger.Warn("error closing body", "error", err)
				}
			}()

			if resp.StatusCode == http.StatusNotFound {
				// Более нет чанков, выходим
				return nil
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("status download is not OK: %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("cannot read download response: %w", err)
			}

			var chunkRes shrModel.ChunkGetResp

			if err := json.Unmarshal(body, &chunkRes); err != nil {
				return err
			}

			if err := s.repo.GetCommonRepo().InsertItemChunk(ctx, itemID, chunkNum, chunkRes.Ciphertext); err != nil {
				return fmt.Errorf("cannot put chunk to client db: %w", err)
			}

			chunkNum++
		}
	}

	return nil
}
