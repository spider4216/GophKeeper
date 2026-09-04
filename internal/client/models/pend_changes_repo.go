package models

import "github.com/spider4216/GophKeeper/internal/enum"

type PendChangesRepo struct {
	ItemID    string
	Operation enum.OpType
	UserID    int64
}
