package data

import (
	"context"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
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

func (dao *HolderDAO) GetHolderListByCollectionId(ctx context.Context, collectionId int64) (holderList []biz.HolderBO, err error) {
	return
}

func (dao *HolderDAO) UpdateHolderBalance(ctx context.Context, id int64, overallBalance, availableBalance float64) error {
	err := dao.data.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		return tx.Model(db.HolderPO{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"overall_balance":   overallBalance,
				"available_balance": availableBalance,
				"update_time":       time.Now().Unix(),
			}).Error
	})
	return err
}

func (dao *HolderDAO) GetHolderInfo(ctx context.Context, address string, pkscript string, chainId string, tick string, collectionId int64) (holderInfo biz.HolderBO, err error) {
	err = dao.data.db.SqlDB.Model(db.HolderPO{}).
		Where("address = ?", address).
		Where("pkscript = ?", pkscript).
		Where("chain_id = ?", chainId).
		Where("collection_id = ?", collectionId).First(&holderInfo).Error
	return holderInfo, err
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	holder := biz.HolderBO{
	//		Address:      address,
	//		Pkscript:     pkscript,
	//		ChainId:      chainId,
	//		CollectionId: collectionId,
	//		Tick:         tick,
	//		CreateTime:   time.Now().Unix(),
	//		UpdateTime:   time.Now().Unix(),
	//	}
	//	err = dao.data.db.SqlDB.Model(db.HolderPO{}).Create(&holder).Error
	//	return holder, err
	//}
	//return holderInfo, err
}
