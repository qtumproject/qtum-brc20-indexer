package service

import (
	indexer_v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/config"
	"github.com/google/wire"
)

// SvcProviderSet is serviceServer providers.
var SvcProviderSet = wire.NewSet(NewOrdinalIndexerServiceServer)

// OrdinalIndexerInterfaceServer Server struct
type OrdinalIndexerInterfaceServer struct {
	indexer_v1.UnimplementedOrdinalIndexerInterfaceServer
	collectionService          *biz.CollectionService
	ordHistoricBalancesService *biz.OrdHistoricBalancesService
	holderService              *biz.HolderService
	conf                       *config.Config
}

// NewOrdinalIndexerServiceServer New app grpc server
func NewOrdinalIndexerServiceServer(conf *config.Config, cc *biz.CollectionService, tc *biz.OrdHistoricBalancesService, hc *biz.HolderService) indexer_v1.OrdinalIndexerInterfaceServer {
	return &OrdinalIndexerInterfaceServer{
		collectionService:          cc,
		ordHistoricBalancesService: tc,
		holderService:              hc,
		conf:                       conf,
	}
}
