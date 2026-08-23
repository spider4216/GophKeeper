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

func (s *Service) SyncSend(ctx context.Context, userID int64, token string) error {
	// todo transavtion

	s.logger.Debug("Get pendings...")

	pends, err := s.repo.GetPendingUserChanges(ctx, int(userID))

	if err != nil {
		return err
	}

	s.logger.Debugf("Pendings count: %d", len(pends))

	var itemIDs []int64

	for _, pend := range pends {
		itemIDs = append(itemIDs, pend.ItemID)
	}

	items, err := s.repo.GetItemsByIDs(ctx, itemIDs)

	s.logger.Debugf("User items sync: %d", len(items))

	if err != nil {
		return err
	}

	// todo моэно в горутину отправитьт вместе с предыдущим получением
	meta, err := s.repo.GetCommonRepo().GetMetadataByItemIDs(ctx, itemIDs)

	if err != nil {
		return err
	}

	s.logger.Debug(pends)
	s.logger.Debug(items)
	s.logger.Debug(meta)

	req := s.buildSyncRequest(pends, items, meta)

	data, err := json.Marshal(req)

	if err != nil {
		return err
	}

	// todo endpoint to const
	url, err := url.JoinPath(s.host, "/sync")

	if err != nil {
		return err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", token)

	s.logger.Debug("Send sync...")
	s.logger.Debug(string(data))
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

	var res shrModel.SyncPutResp

	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	s.logger.Debugf("Latest revision: %d, for user: %d", res.LastRev, userID)

	s.logger.Debug("Delete pendings")
	if err := s.repo.DeletePendingByItemIDs(ctx, itemIDs); err != nil {
		return fmt.Errorf("Cannot delete pending changes: %s", err)
	}

	s.logger.Debugf("Update latest user revision to %d", res.LastRev)

	if err := s.repo.UpdateLastUserRev(ctx, userID, res.LastRev); err != nil {
		return fmt.Errorf("cannot update to latest rev version: %s", err)
	}

	s.logger.Debug("Sync done")

	return nil
}

func (s *Service) SyncGet(ctx context.Context, userID int64, token string) error {
	// Получаем последнюю версию синхронизации
	rev, err := s.GetLatestUserRev(ctx, userID)

	if err != nil {
		return err
	}

	// todo endpoint to const
	url, err := url.JoinPath(s.host, "/sync")

	if err != nil {
		return err
	}

	revStr := strconv.FormatInt(rev, 10)

	url = url + "?since=" + revStr

	r, err := http.NewRequest(http.MethodGet, url, nil)
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
		return err
	}

	s.logger.Debugf("Changes: %d", len(res.Changes))

	// todo transaction
	for _, change := range res.Changes {
		// todo в зависимости от операции метод будет разный, пока просто create
		// todo mybe gorutine

		if change.Operation == "CREATE" {
			item := models.ItemRepo{
				Type:       change.Item.Type,
				Ciphertext: change.Item.Ciphertext,
				UserID:     userID,
			}

			itemID, err := s.repo.CreateItem(ctx, item)

			if err != nil {
				return err
			}

			for k, v := range change.Metadata {
				if _, err := s.repo.CreateMeta(ctx, itemID, k, v); err != nil {
					return err
				}
			}
		}

		if change.Operation == "DELETE" {
			// todo transaction
			if err := s.repo.DeleteUserMetaByItemID(ctx, change.Item.ID, userID); err != nil {
				return err
			}

			if err := s.repo.DeleteUserItemByID(ctx, change.Item.ID, userID); err != nil {
				return err
			}
		}

		if change.Operation == "UPDATE" {
			// todo update
			if err := s.repo.UpdateUserItem(ctx, change.Item.ID, userID, change.Item.Ciphertext); err != nil {
				return err
			}
		}

	}

	if err := s.UpdateLastUserRev(ctx, userID, res.NextRev); err != nil {
		return err
	}

	return nil
}
