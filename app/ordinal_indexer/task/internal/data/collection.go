package data

import (
	"context"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
	"gorm.io/gorm"
	"time"
)

var _ biz.ICollectionDAO = (*CollectionDAO)(nil) //编译时接口检查，确保userRepo实现了biz.CollectionRepo定义的接口

var CollectionCacheKey = func(taskCode string) string {
	return "ordinal_indexer_cache_key_" + taskCode
}

type CollectionDAO struct {
	data *Data
}

func NewCollectionDAO(data *Data) biz.ICollectionDAO {
	return &CollectionDAO{
		data: data,
	}
}

func (dao *CollectionDAO) GetCollectionInfo(ctx context.Context, chainId, protocol, tick string) (bo biz.CollectionBO, err error) {
	err = dao.data.db.SqlDB.Model(&db.CollectionPO{}).
		Where("chain_id = ?", chainId).
		Where("protocol = ?", protocol).
		Where("tick = ?", tick).
		Find(&bo).Error
	return
}

func (dao *CollectionDAO) CreateCollection(ctx context.Context, collectionInfo *biz.CollectionBO) error {
	err := dao.data.db.SqlDB.Model(&db.CollectionPO{}).Create(collectionInfo).Error
	return err
}

func (dao *CollectionDAO) UpdateCollectionInfo(ctx context.Context, id int64, bo biz.CollectionBO) error {
	return dao.data.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		updateValues := map[string]interface{}{
			"total_minted":        bo.TotalMinted,
			"mint_times":          bo.MintTimes,
			"last_block_height":   bo.LastBlockHeight,
			"last_transaction_id": bo.LastTransactionId,
			"update_time":         time.Now().Unix(),
		}
		return tx.Model(&db.CollectionPO{}).Where("id = ?", id).
			Updates(updateValues).Error
	})
}
