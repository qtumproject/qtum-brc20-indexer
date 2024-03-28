package biz

import (
	"errors"
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/google/wire"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

const (
	MaxPageSize = 1000
	MinPageSize = 1
)

// BzProviderSet is biz providers.
var BzProviderSet = wire.NewSet(NewOrdHistoricBalancesService, NewHolderService, NewCollectionService, NewOrdinalEventService)

func ValidatePaginate(page, pageSize int64) v1.Pagination {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	switch {
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	case pageSize < MinPageSize:
		pageSize = MinPageSize
	}
	return v1.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
