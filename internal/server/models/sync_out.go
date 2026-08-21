package models

// todo объйдинить модели клиента и сервера, они одинаковые

type SyncOutReq struct {
	Changes []SyncOutChange `json:"changes"`
	NextRev int64           `json:"next_revision"`
}

type SyncOutChange struct {
	Operation string            `json:"operation"`
	Item      ItemSyncOut       `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncOut struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
