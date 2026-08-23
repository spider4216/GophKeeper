package model

type SyncPutReq struct {
	Changes []SyncPutChange `json:"changes"`
}

type SyncPutChange struct {
	Operation string            `json:"operation"`
	Item      ItemSyncPut       `json:"item"`
	Metadata  map[string]string `json:"metadata"`
}

type ItemSyncPut struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Ciphertext string `json:"ciphertext"`
}
