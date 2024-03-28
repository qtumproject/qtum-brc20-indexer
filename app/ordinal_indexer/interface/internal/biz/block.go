package biz

import (
	"context"
)

type IBlockDAO interface {
	//FetchOrdinalEvent 从指定block获取OrdinalEvent信息
	GetBlockInfoByChainId(ctx context.Context, chainId string) (BlockBO, error)
	UpdateBlockInfo(ctx context.Context, bo BlockBO) error
	UpdateSyncedBlockNumber(ctx context.Context, chainId string, blockHeight int64, blockHash string) error
	ListSupportedChainBlock(ctx context.Context) ([]BlockBO, error)
}

type BlockBO struct {
	ID                      int64  `gorm:"AUTO_INCREMENT" json:"id,omitempty"`
	BlockHash               string `gorm:"block_hash" json:"block_hash"`
	ChainId                 string `gorm:"chain_id" json:"chain_id"`
	Status                  int64  `gorm:"status" json:"status"`
	SyncCountPerTime        int64  `gorm:"sync_count_per_time" json:"sync_count_per_time"`
	LatestSyncedBlockHeight int64  `gorm:"Latest_synced_block_height" json:"Latest_synced_block_height"` //最新已同步的区块
	TotalBlockHeight        int64  `gorm:"total_block_height" json:"total_block_height"`                 //总区块数量
	CreateTime              int64  `gorm:"create_time" json:"create_time"`
	UpdateTime              int64  `gorm:"update_time" json:"update_time"`
}
