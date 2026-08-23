package models

import shrMocdel "github.com/spider4216/GophKeeper/internal/model"

type ItemWithMeta struct {
	ItemRepo
	Metadata []shrMocdel.MetadataRepo
}
