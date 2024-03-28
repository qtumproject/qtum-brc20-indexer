package data

import (
	"context"
	"fmt"
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data/db"
	"gorm.io/gorm"
	"math"
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

func (dao *CollectionDAO) GetCollectionList(
	ctx context.Context,
	chainId string,
	tickSearchQuery string,
	status string,
	pagination v1.Pagination) ([]biz.CollectionBO, v1.Pagination, error) {
	var list []biz.CollectionBO
	if len(tickSearchQuery) > 4 {
		return list, pagination, fmt.Errorf("tickSearchQuery param too long")
	}

	query := dao.data.db.SqlDB.Model(&db.CollectionPO{}).
		Where("chain_id = ?", chainId)
	if tickSearchQuery != "" {
		query = query.Where("tick LIKE ?", tickSearchQuery+"%")
	}
	if status == "in-progress" {
		query = query.Where("total_minted != max_amount")
	} else if status == "completed" {
		query = query.Where("total_minted = max_amount")
	}
	query = query.Order("deploy_time DESC")
	query.Count(&pagination.Total)
	pagination.TotalPage = int64(math.Ceil(float64(pagination.Total) / float64(pagination.PageSize)))
	err := query.Scopes(db.Paginate(pagination)).Find(&list).Error
	return list, pagination, err
}

func (dao *CollectionDAO) GetCollectionInfo(ctx context.Context, chainId, protocol, tick string) (bo biz.CollectionBO, err error) {
	if protocol == "" {
		protocol = "brc-20"
	}
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
