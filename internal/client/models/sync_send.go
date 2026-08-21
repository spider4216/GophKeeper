package models

// todo объйдинить модели клиента и сервера, они одинаковые

type SyncInReq struct {
	Changes []SyncInChange `json:"changes"`
}

type SyncInChange struct {
	Operation string            `json:"operation"`
	Item      ItemSyncIn        `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncIn struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
