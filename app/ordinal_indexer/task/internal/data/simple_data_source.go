package data

import (
	"context"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
)

type SimpleDataSourceDAO struct{}

func (dao *SimpleDataSourceDAO) BatchFetchOrdinalEvent(
	ctx context.Context,
	startBlockHeight,
	blockCount int64) (eventList []biz.OrdEventBO, lastBlockHeight int64, lastBlockHash string, err error) {
	return
}
func (dao *SimpleDataSourceDAO) GetChainId() string {
	return "simple"
}

func (dao *SimpleDataSourceDAO) GetBlockNumber(ctx context.Context) (blockHeight int64, err error) {
	return
}

func NewSimpleDataSourceDAO(config *config.Config) biz.IDataSourceDAO {
	return &SimpleDataSourceDAO{}
}
