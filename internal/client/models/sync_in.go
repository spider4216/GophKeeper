package models

// todo объйдинить модели клиента и сервера, они одинаковые
// todo это тоже самое что SyncOut в сервере

type SyncReceiveReq struct {
	Changes []SyncReceiveChange `json:"changes"`
	NextRev int64               `json:"next_revision"`
}

type SyncReceiveChange struct {
	Operation string            `json:"operation"`
	Item      ItemSyncReceive   `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncReceive struct {
	ID         int64  `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
