package services

import (
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

func (s *Service) mapSyncResponse(
	changes []models.SyncChangesRepo,
	items []models.ItemRepo,
	payloads []models.ItemPayloadRepo,
	metadata []shrModel.MetadataRepo,
	since int64,
) *shrModel.SyncGet {
	itemByID := make(map[int64]models.ItemRepo, len(items))
	for _, item := range items {
		itemByID[item.ID] = item
	}

	payloadByItemID := make(map[int64]models.ItemPayloadRepo, len(payloads))
	for _, payload := range payloads {
		payloadByItemID[payload.ItemID] = payload
	}

	metadataByItemID := make(map[int64]map[string]string)

	for _, m := range metadata {
		if metadataByItemID[m.ItemID] == nil {
			metadataByItemID[m.ItemID] = make(map[string]string)
		}

		metadataByItemID[m.ItemID][m.Key] = m.Value
	}

	result := shrModel.SyncGet{
		Changes: make([]shrModel.SyncGetChange, 0, len(changes)),
		NextRev: since,
	}

	for _, change := range changes {
		syncChange := shrModel.SyncGetChange{
			Operation: change.Operation,
			Metadata:  metadataByItemID[change.ItemID],
		}

		// DELETE: item уже может отсутствовать.
		if change.Operation == "DELETE" {
			syncChange.Item.ID = change.ItemID
			result.Changes = append(result.Changes, syncChange)
			result.NextRev = change.Revision
			continue
		}

		item := itemByID[change.ItemID]
		payload := payloadByItemID[change.ItemID]

		syncChange.Item = shrModel.ItemSyncGet{
			ID:         item.ID,
			Type:       item.Type,
			Ciphertext: payload.Ciphertext,
		}

		result.Changes = append(result.Changes, syncChange)
		result.NextRev = change.Revision
	}

	return &result
}
