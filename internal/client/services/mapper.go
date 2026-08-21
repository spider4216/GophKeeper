package services

import "github.com/spider4216/GophKeeper/internal/client/models"

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
