package biz

import (
	"errors"
	"github.com/google/wire"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// BzProviderSet is biz providers.
var BzProviderSet = wire.NewSet(NewOrdHistoricBalancesService, NewHolderService, NewCollectionService, NewOrdinalEventService)
