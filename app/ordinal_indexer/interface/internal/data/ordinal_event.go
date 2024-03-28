package data

import (
	"context"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data/db"
)

var _ biz.IOrdinalEventDAO = (*OrdEventDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var OrdinalEventCacheKey = func(taskCode string) string {
	return "reward_item_cache_key_" + taskCode
}

type OrdEventDAO struct {
	data *Data
}

func (r *OrdEventDAO) GetEventList(
	ctx context.Context,
	chainId string,
	collectionId int64,
	inscriptionId string,
	eventStatus string,
	maxBlockHeight int64,
	fromAddressList []string,
	toAddressList []string) (eventList []biz.OrdEventBO, err error) {

	query := r.data.db.SqlDB.Model(&db.OrdEventPO{})
	if chainId != "" {
		query = query.Where("chain_id = ?", chainId)
	}
	if collectionId != 0 {
		query = query.Where("collection_id = ?", collectionId)
	}
	if inscriptionId != "" {
		query = query.Where("inscription_id = ?", inscriptionId)
	}
	if eventStatus != "" {
		query = query.Where("event_status = ?", eventStatus)
	}
	if maxBlockHeight > 0 {
		query = query.Where("block_height <= ?", maxBlockHeight)
	}
	if len(fromAddressList) > 0 && len(toAddressList) > 0 {
		query = query.Where("from_address IN ? OR to_address IN ?", fromAddressList, toAddressList)
	} else if len(fromAddressList) > 0 {
		query = query.Where("from_address IN ?", fromAddressList)
	} else if len(toAddressList) > 0 {
		query = query.Where("to_address IN ?", toAddressList)
	}
	query = query.Order("id ASC")
	//query = query.Or("operation = ?", biz.BRC20_OP_DEPLOY)
	err = query.Find(&eventList).Error
	return
}

//GetRelateEventList 查询事件列表

func NewOrdinalEventDAO(data *Data) biz.IOrdinalEventDAO {
	return &OrdEventDAO{
		data: data,
	}
}
