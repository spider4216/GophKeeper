package models

// todo объеденить с сервером

type SyncSendResp struct {
	LastRev int64 `json:"last_rev"`
}
