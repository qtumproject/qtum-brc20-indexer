package biz

import (
	"context"
)

type IDataSourceDAO interface {
	//BatchFetchOrdinalEvent 批量block获取OrdinalEvent信息。 startBlockHeight: 开始高度, blockCount: 区块数量
	BatchFetchOrdinalEvent(ctx context.Context, startBlockHeight, blockCount int64) ([]OrdEventBO, int64, string, error)
	GetBlockNumber(ctx context.Context) (blockHeight int64, err error)
	GetChainId() string
}
