package db

import (
	"github.com/6block/fox_ordinal/pkg/cache"
	"gorm.io/gorm"
)

const (
	MysqlDriver = "mysql"
	MaxPageSize = 1000
	MinPageSize = 1
)

type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type ServerDatabase struct {
	SqlDB    *gorm.DB
	CacheCli *cache.FoxRedis
}

func ValidatePaginate(page, pageSize int) Pagination {
	if page == 0 {
		page = 1
	}
	switch {
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	case pageSize < MinPageSize:
		pageSize = MinPageSize
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func Paginate(pagination Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		//page, pageSize = ValidatePaginate(page, pageSize)
		offset := (pagination.Page - 1) * pagination.PageSize
		return db.Offset(offset).Limit(pagination.PageSize)
	}
}
