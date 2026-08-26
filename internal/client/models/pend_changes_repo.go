package models

import "github.com/spider4216/GophKeeper/internal/enum"

type PendChangesRepo struct {
	ItemID    int64
	Operation enum.OpType
	UserID    int64
}
