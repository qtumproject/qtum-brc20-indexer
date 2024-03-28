package service

import (
	"context"
	task_v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/task/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"github.com/google/wire"
)

// SvcProviderSet is serviceServer providers.
var SvcProviderSet = wire.NewSet(NewOrdinalIndexerServiceServer)

// OrdinalIndexerTaskServer Server struct
type OrdinalIndexerTaskServer struct {
	task_v1.UnimplementedOrdinalIndexerTaskServer
	collectionService          *biz.CollectionService
	ordHistoricBalancesService *biz.OrdHistoricBalancesService
	holderService              *biz.HolderService
	ordinalEventService        *biz.OrdinalEventService
	conf                       *config.Config
}

// NewOrdinalIndexerServiceServer New app grpc server
func NewOrdinalIndexerServiceServer(
	conf *config.Config,
	cc *biz.CollectionService,
	tc *biz.OrdHistoricBalancesService,
	hc *biz.HolderService,
	ec *biz.OrdinalEventService,
) *OrdinalIndexerTaskServer {
	return &OrdinalIndexerTaskServer{
		collectionService:          cc,
		ordHistoricBalancesService: tc,
		holderService:              hc,
		conf:                       conf,
		ordinalEventService:        ec,
	}
}

func (svc OrdinalIndexerTaskServer) Ping(context.Context, *task_v1.PingReq) (*task_v1.PingReply, error) {
	return &task_v1.PingReply{}, nil
}
