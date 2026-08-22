package model

type SyncGet struct {
	Changes []SyncGetChange `json:"changes"`
	NextRev int64           `json:"next_revision"`
}

type SyncGetChange struct {
	Operation string            `json:"operation"`
	Item      ItemSyncGet       `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncGet struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
