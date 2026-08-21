package services

import "github.com/spider4216/GophKeeper/internal/client/models"

func (s *Service) buildItemsWithMeta(items []models.ItemRepo, metadata []models.MetadataRepo) []models.ItemWitjMeta {
	// todo move to func because dry in buildSyncRequest

	metadataByItemID := make(map[int64]map[string]string)

	for _, m := range metadata {
		if metadataByItemID[m.ItemID] == nil {
			metadataByItemID[m.ItemID] = make(map[string]string)
		}

		metadataByItemID[m.ItemID][m.Key] = m.Value
	}

	// todo rename model - misatke
	var res []models.ItemWitjMeta

	for _, item := range items {
		meta := metadataByItemID[item.ID]

		one := models.ItemWitjMeta{
			ItemRepo: models.ItemRepo{
				ID:         item.ID,
				Type:       item.Type,
				Ciphertext: item.Ciphertext,
				UserID:     item.UserID,
				CreatedAt:  item.CreatedAt,
			},
		}

		for k, v := range meta {
			one.Metadata = append(one.Metadata, models.MetadataRepo{
				Key:   k,
				Value: v,
			})
		}

		res = append(res, one)
	}

	return res
}

func (s *Service) buildSyncRequest(pendingChanges []models.PendChangesRepo, items []models.ItemRepo, metadata []models.MetadataRepo) models.SyncInReq {
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

	req := models.SyncInReq{
		Changes: make([]models.SyncInChange, 0, len(pendingChanges)),
	}

	for _, pending := range pendingChanges {
		item, ok := itemByID[pending.ItemID]
		if !ok {
			// Для DELETE item может уже отсутствовать.
			req.Changes = append(req.Changes, models.SyncInChange{
				Operation: pending.Operation,
				Item: models.ItemSyncIn{
					ID: int(pending.ItemID),
				},
			})

			continue
		}

		req.Changes = append(req.Changes, models.SyncInChange{
			Operation: pending.Operation,
			Item: models.ItemSyncIn{
				ID:         int(item.ID),
				Type:       item.Type,
				Ciphertext: item.Ciphertext,
			},
			Metadata: metadataByItemID[item.ID],
		})
	}

	return req
}
