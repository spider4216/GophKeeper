package models

import "time"

type SyncChangesRepo struct {
	ID        int64
	Revision  int64
	ItemID    int64
	Operation string // todo operation custom type
	CreatedAt time.Time
}
