package models

import shrModel "github.com/spider4216/GophKeeper/internal/model"

type ItemWithMeta struct {
	ItemRepo
	Metadata []shrModel.MetadataRepo
}
