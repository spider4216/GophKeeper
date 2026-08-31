package services

import (
	"github.com/spider4216/GophKeeper/internal/enum"
	shrModel "github.com/spider4216/GophKeeper/internal/model"
	"github.com/spider4216/GophKeeper/internal/server/models"
)

func (s *Service) mapSyncResponse(
	changes []syncChange,
	items []models.ItemRepo,
	payloads []models.ItemPayloadRepo,
	metadata []shrModel.MetadataRepo,
	nextRev int64,
	hasMore bool,
) *shrModel.SyncGet {
	itemByID := make(map[string]models.ItemRepo, len(items))
	for _, item := range items {
		itemByID[item.ID] = item
	}

	payloadByItemID := make(map[string]models.ItemPayloadRepo, len(payloads))
	for _, payload := range payloads {
		payloadByItemID[payload.ItemID] = payload
	}

	metadataByItemID := make(map[string]map[string]string)

	for _, m := range metadata {
		if metadataByItemID[m.ItemID] == nil {
			metadataByItemID[m.ItemID] = make(map[string]string)
		}

		metadataByItemID[m.ItemID][m.Key] = m.Value
	}

	result := shrModel.SyncGet{
		Changes: make([]shrModel.SyncGetChange, 0, len(changes)),
		NextRev: nextRev,
		HasMore: hasMore,
	}

	for _, change := range changes {
		syncChange := shrModel.SyncGetChange{
			Operation: change.operation,
			Metadata:  metadataByItemID[change.itemID],
		}

		// DELETE: item уже может отсутствовать.
		if change.operation == enum.OpDelete {
			syncChange.Item.ID = change.itemID
			result.Changes = append(result.Changes, syncChange)
			continue
		}

		item := itemByID[change.itemID]
		payload := payloadByItemID[change.itemID]

		syncChange.Item = shrModel.ItemSyncGet{
			ID:         item.ID,
			Type:       item.Type,
			Ciphertext: payload.Ciphertext,
		}

		result.Changes = append(result.Changes, syncChange)
	}

	return &result
}
