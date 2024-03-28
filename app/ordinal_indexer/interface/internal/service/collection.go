package service

import (
	"context"
	"fmt"
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/pkg/utils"
)

func (svc OrdinalIndexerInterfaceServer) ListTickers(
	ctx context.Context, req *v1.ListTickersReq) (*v1.ListTickersReply, error) {
	pagination := biz.ValidatePaginate(req.GetPage(), req.GetPageSize())

	list, pagination, err := svc.collectionService.GetCollectionList(ctx, req.GetChainId(), req.GetTick(), req.GetStatus(), pagination)
	if err != nil {
		return nil, err
	}
	tickerList := make([]*v1.TickerInfo, 0)

	for _, record := range list {
		holderList, _ := svc.holderService.GetHolderListByCollectionId(ctx, req.GetChainId(), record.ID)
		process := (record.TotalMinted / record.MaxAmount) * 100
		vo := v1.TickerInfo{
			TokenName:  record.Tick,
			DeployTime: utils.TimestampUnixToFormat(record.CreateTime, nil),
			Progress:   fmt.Sprintf("%.03f%%", process),
			Holders:    int64(len(holderList)),
			MintTimes:  record.MintTimes,
		}
		if record.DeployTime > 0 {
			vo.DeployTime = utils.TimestampUnixToFormat(record.DeployTime, nil)
		}
		tickerList = append(tickerList, &vo)
	}
	reply := &v1.ListTickersReply{
		TickerList:     tickerList,
		PaginationInfo: &pagination,
	}
	return reply, nil
}
