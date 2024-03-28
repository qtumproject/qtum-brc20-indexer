package db

import (
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/6block/fox_ordinal/pkg/cache"
	"gorm.io/gorm"
)

const (
	MysqlDriver = "mysql"
)

type ServerDatabase struct {
	SqlDB    *gorm.DB
	CacheCli *cache.FoxRedis
}

func Paginate(pagination v1.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//page, pageSize = ValidatePaginate(page, pageSize)
		offset := (pagination.Page - 1) * pagination.PageSize
		return db.Offset(int(offset)).Limit(int(pagination.PageSize))
	}
}
