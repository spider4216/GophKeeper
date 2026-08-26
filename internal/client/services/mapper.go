package services

import (
	"github.com/spider4216/GophKeeper/internal/client/models"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

func (s *Service) buildItemsWithMeta(items []models.ItemRepo, metadata []shrModel.MetadataRepo) []models.ItemWithMeta {
	metadataByItemID := make(map[int64][]shrModel.MetadataRepo)

	for _, m := range metadata {
		metadataByItemID[m.ItemID] = append(metadataByItemID[m.ItemID], m)
	}

	var res []models.ItemWithMeta

	for _, item := range items {
		meta := metadataByItemID[item.ID]

		one := models.ItemWithMeta{
			ItemRepo: models.ItemRepo{
				ID:         item.ID,
				Type:       item.Type,
				Ciphertext: item.Ciphertext,
				UserID:     item.UserID,
				CreatedAt:  item.CreatedAt,
			},
		}

		for _, m := range meta {
			one.Metadata = append(one.Metadata, shrModel.MetadataRepo{
				ID:    m.ID,
				Key:   m.Key,
				Value: m.Value,
			})
		}

		res = append(res, one)
	}

	return res
}

func (s *Service) buildSyncRequest(pendingChanges []models.PendChangesRepo, items []models.ItemRepo, metadata []shrModel.MetadataRepo) shrModel.SyncPutReq {
	itemByID := make(map[int64]models.ItemRepo, len(items))

	for _, item := range items {
		itemByID[item.ID] = item
	}

	metadataByItemID := make(map[int64]map[string]string)

	for _, m := range metadata {
		if metadataByItemID[m.ItemID] == nil {
			metadataByItemID[m.ItemID] = make(map[string]string)
		}

		metadataByItemID[m.ItemID][m.Key] = m.Value
	}

	req := shrModel.SyncPutReq{
		Changes: make([]shrModel.SyncPutChange, 0, len(pendingChanges)),
	}

	for _, pending := range pendingChanges {
		item, ok := itemByID[pending.ItemID]
		if !ok {
			// Для DELETE item может уже отсутствовать.
			req.Changes = append(req.Changes, shrModel.SyncPutChange{
				Operation: pending.Operation,
				Item: shrModel.ItemSyncPut{
					ID: int(pending.ItemID),
				},
			})

			continue
		}

		req.Changes = append(req.Changes, shrModel.SyncPutChange{
			Operation: pending.Operation,
			Item: shrModel.ItemSyncPut{
				ID:         int(item.ID),
				Type:       item.Type.String(),
				Ciphertext: item.Ciphertext,
			},
			Metadata: metadataByItemID[item.ID],
		})
	}

	return req
}
