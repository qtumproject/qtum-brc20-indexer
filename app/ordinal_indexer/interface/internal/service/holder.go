package service

import (
	"context"
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
)

func (svc OrdinalIndexerInterfaceServer) ListAddressTickerBalances(
	ctx context.Context, req *v1.ListAddressTickerBalancesReq) (reply *v1.ListAddressTickerBalancesReply, err error) {
	list, err := svc.holderService.GetHolderCollectionList(ctx, req.GetAddress(), req.GetChainId())
	if err != nil {
		return reply, err
	}
	reply = &v1.ListAddressTickerBalancesReply{
		AddressTickerBalanceList: make([]*v1.AddressTickerBalanceInfo, 0),
	}
	for _, record := range list {
		vo := v1.AddressTickerBalanceInfo{
			TokenName:     record.Tick,
			Balance:       record.OverallBalance,
			WalletAddress: record.Address,
			Transferable:  record.OverallBalance - record.AvailableBalance,
			Available:     record.AvailableBalance,
		}
		reply.AddressTickerBalanceList = append(reply.AddressTickerBalanceList, &vo)
	}
	return reply, nil
}

func (svc OrdinalIndexerInterfaceServer) ListHolders(
	ctx context.Context, req *v1.ListHoldersReq) (reply *v1.ListHoldersReply, err error) {
	tickInfo, err := svc.collectionService.GetCollection(ctx, req.GetChainId(), req.GetProtocol(), req.GetTick())
	if err != nil {
		return reply, err
	}
	//pagination := biz.ValidatePaginate(req.GetPage(), req.GetPageSize())
	list, err := svc.holderService.GetHolderListByCollectionId(ctx, req.GetChainId(), tickInfo.ID)
	if err != nil {
		return reply, err
	}
	reply = &v1.ListHoldersReply{
		AddressTickerBalanceList: make([]*v1.AddressTickerBalanceInfo, 0),
	}
	for _, record := range list {
		vo := v1.AddressTickerBalanceInfo{
			TokenName:     record.Tick,
			WalletAddress: record.Address,
			Balance:       record.OverallBalance,
			Transferable:  record.OverallBalance - record.AvailableBalance,
			Available:     record.AvailableBalance,
		}
		reply.AddressTickerBalanceList = append(reply.AddressTickerBalanceList, &vo)
	}
	return reply, nil
}
