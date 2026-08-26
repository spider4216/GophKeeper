package models

import (
	"time"

	"github.com/spider4216/GophKeeper/internal/enum"
)

type SyncChangesRepo struct {
	ID        int64
	Revision  int64
	ItemID    int64
	Operation enum.OpType
	CreatedAt time.Time
}
