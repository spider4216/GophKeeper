package repository

import (
	"context"

	shrModel "github.com/spider4216/GophKeeper/internal/model"
)

type CommonRepositoryInterface interface {
	GetMetadataByItemIDs(ctx context.Context, itemIDs []int64) ([]shrModel.MetadataRepo, error)
}
