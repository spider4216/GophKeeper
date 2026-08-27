package model

import "github.com/spider4216/GophKeeper/internal/enum"

type SyncGet struct {
	Changes []SyncGetChange `json:"changes"`
	NextRev int64           `json:"next_revision"`
	HasMore bool            `json:"has_more"`
}

type SyncGetChange struct {
	Operation enum.OpType       `json:"operation"`
	Item      ItemSyncGet       `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncGet struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
