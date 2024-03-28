package data

import (
	"context"
	"errors"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data/db"
	"gorm.io/gorm"
	"time"
)

var _ biz.IHolderDAO = (*HolderDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var HolderCacheKey = func(taskCode string) string {
	return "ordinal_indexer_cache_key_" + taskCode
}

type HolderDAO struct {
	data *Data
}

func NewHolderDAO(data *Data) biz.IHolderDAO {
	return &HolderDAO{
		data: data,
	}
}

func (dao *HolderDAO) GetHolderListByChainId(ctx context.Context, chainId string) (list []biz.HolderBO, err error) {
	return
}

func (dao *HolderDAO) GetHolderListByCollectionId(ctx context.Context, chainId string, collectionId int64) (holderList []biz.HolderBO, err error) {
	holderList = make([]biz.HolderBO, 0)
	err = dao.data.db.SqlDB.Model(db.HolderPO{}).
		Where("chain_id = ?", chainId).
		Where("collection_id = ?", collectionId).
		Find(&holderList).Error
	return
}

func (dao *HolderDAO) GetHolderCollectionList(ctx context.Context, address, chainId string) (holderList []biz.HolderBO, err error) {
	holderList = make([]biz.HolderBO, 0)
	err = dao.data.db.SqlDB.Model(db.HolderPO{}).
		Where("address = ?", address).
		Where("chain_id = ?", chainId).
		Find(&holderList).Error
	return
}

func (dao *HolderDAO) UpdateAmount(ctx context.Context, id int64, amount float64) error {
	return dao.data.db.SqlDB.Model(db.HolderPO{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"amount":      amount,
			"update_time": time.Now().Unix(),
		}).Error
}

func (dao *HolderDAO) GetHolderInfo(ctx context.Context, address string, chainId string, collectionId int64) (holderInfo biz.HolderBO, err error) {
	err = dao.data.db.SqlDB.Model(db.HolderPO{}).
		Where("address = ?", address).
		Where("chain_id = ?", chainId).
		Where("collection_id = ?", collectionId).First(&holderInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		holder := biz.HolderBO{
			Address:      address,
			ChainId:      chainId,
			CollectionId: collectionId,
			CreateTime:   time.Now().Unix(),
			UpdateTime:   time.Now().Unix(),
		}
		err = dao.data.db.SqlDB.Model(db.HolderPO{}).Create(&holder).Error
		return holder, err
	}
	return holderInfo, err
}
