package data

import (
	"context"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
	"time"
)

var _ biz.IBlockDAO = (*BlockDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var BlockCacheKey = func(taskCode string) string {
	return "reward_item_cache_key_" + taskCode
}

type BlockDAO struct {
	db *db.ServerDatabase
}

func (r *BlockDAO) GetBlockInfoByChainId(ctx context.Context, chainId string) (bo biz.BlockBO, err error) {
	err = r.db.SqlDB.Model(&db.BlockPO{}).
		Where("chain_id = ?", chainId).
		Where("status = ?", 1).
		First(&bo).Error
	return
}
func (r *BlockDAO) UpdateBlockInfo(ctx context.Context, bo biz.BlockBO) error {
	//po := db.BlockPO{
	//	ID: bo.ID,
	//	BlockHash: bo.BlockHash,
	//	ChainId: bo.ChainId,
	//	Status: bo.Status,
	//	SyncCountPerTime: bo.SyncCountPerTime,
	//	LatestSyncedBlockHeight: bo.LatestSyncedBlockHeight,
	//	TotalBlockHeight: bo.TotalBlockHeight,
	//	CreateTime: bo.CreateTime,
	//	UpdateTime: time.Now().Unix(),
	//}
	return r.db.SqlDB.Model(&db.BlockPO{}).Where("id = ?", bo.ID).Updates(map[string]interface{}{
		"block_hash":                 bo.BlockHash,
		"chain_id":                   bo.ChainId,
		"status":                     bo.Status,
		"sync_count_per_time":        bo.SyncCountPerTime,
		"Latest_synced_block_height": bo.LatestSyncedBlockHeight,
		"total_block_height":         bo.TotalBlockHeight,
		"create_time":                bo.CreateTime,
		"update_time":                time.Now().Unix(),
	}).Error
}

func (r *BlockDAO) UpdateSyncedBlockNumber(ctx context.Context, chainId string, blockHeight int64, blockHash string) error {
	return r.db.SqlDB.Model(&db.BlockPO{}).Where("chain_id = ?", chainId).Updates(
		map[string]interface{}{
			"block_hash":                 blockHash,
			"Latest_synced_block_height": blockHeight,
			"update_time":                time.Now().Unix(),
		}).Error
}

// ListSupportedChainBlock 获取所有支持的链的block同步信息
func (r *BlockDAO) ListSupportedChainBlock(ctx context.Context) (boList []biz.BlockBO, err error) {
	err = r.db.SqlDB.Model(&db.BlockPO{}).Where("status = ?", 1).Find(&boList).Error
	return
}

func NewBlockDAO(db *db.ServerDatabase) biz.IBlockDAO {
	return &BlockDAO{
		db: db,
	}
}
