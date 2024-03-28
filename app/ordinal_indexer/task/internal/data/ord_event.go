package data

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var _ biz.IOrdinalEventDAO = (*OrdEventDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var OrdinalEventCacheKey = func(taskCode string) string {
	return "reward_item_cache_key_" + taskCode
}

type OrdEventDAO struct {
	data         *Data
	ordHolderDao biz.IHolderDAO
}

// ReorgEvent
// 获取大于同步高度的所有事件，按高度和id升序排
// 从前往后遍历要删的事件
// //如果是deploy，加入delete列表，直接删collection跟所有holder
// 如果是mint并且tick不在delete列表，把mint的tick加入到collection待更新列表
// 记录每个受影响的holder
// 删除以上事件
// 删除block_height大于newBlockHeight的所有余额变更历史
// 更新受影响的holder的余额（在余额变更历史中找到最新一条变更历史，更新余额到holder表）
// 更新受影响的collection的mint次数，剩余可mint余额，上次mint高度和txid（从事件表中查找指定collection的所有合法mint事件，统计次数和总额度）
func (r *OrdEventDAO) ReorgEvent(ctx context.Context, chainId string, newBlockHeight int64) error {
	var eventInfoList []biz.OrdEventBO

	callFunction := func(tx *gorm.DB) error {
		err := tx.Model(&db.OrdEventPO{}).Where("chain_id = ?", chainId).
			Where("block_height > ?", newBlockHeight).
			Order("block_height asc, id asc").
			Find(&eventInfoList).Error
		if err != nil {
			return err
		}
		if len(eventInfoList) == 0 {
			return nil
		}
		deleteCollection := make(map[int64]bool)
		waitForUpdateCollection := make(map[int64]bool)
		CollectionIdToAddressList := make(map[int64]map[string]bool)
		eventIds := make([]int64, 0)
		for _, event := range eventInfoList {
			eventIds = append(eventIds, event.ID)

			if event.EventStatus != "processed" || event.CollectionId == 0 {
				continue
			}

			if _, ok := CollectionIdToAddressList[event.CollectionId]; !ok {
				CollectionIdToAddressList[event.CollectionId] = make(map[string]bool)
			}

			switch event.EventType {
			case biz.DeployInscribe:
				err = tx.Where("id = ?", event.CollectionId).
					Delete(&db.CollectionPO{}).Error
				if err != nil {
					return err
				}
				err = tx.Where("chain_id = ?", chainId).
					Where("collection_id = ?", event.CollectionId).
					Delete(&db.HolderPO{}).Error
				if err != nil {
					return err
				}
				deleteCollection[event.CollectionId] = true

			case biz.MintInscribe:
				//如果是mint并且tick不在delete列表，把mint的tick加入到collection待更新列表
				if _, ok := deleteCollection[event.CollectionId]; !ok {
					waitForUpdateCollection[event.CollectionId] = true
				}
				if event.TargetAddress != "" {
					CollectionIdToAddressList[event.CollectionId][event.TargetAddress] = true
				}
			case biz.TransferInscribe, biz.TransferTransfer:
				if event.TargetAddress != "" {
					CollectionIdToAddressList[event.CollectionId][event.TargetAddress] = true
				}
				if event.SourceAddress != "" {
					CollectionIdToAddressList[event.CollectionId][event.SourceAddress] = true
				}
			}
		}
		//删除事件
		err = tx.Where("id in (?)", eventIds).Delete(&db.OrdEventPO{}).Error
		if err != nil {
			return err
		}
		//删除block_height大于newBlockHeight的所有余额变更历史
		err = tx.Where("block_height > ?", newBlockHeight).Delete(&db.OrdHistoricBalancesPO{}).Error
		if err != nil {
			return err
		}
		queryTemplate :=
			`WITH tempp AS (
    			SELECT MAX(id) AS id
    			FROM ord_historic_balances
				Where collection_id = {{CollectionID}} AND
				wallet_address IN (?)
   			 	GROUP BY wallet_address
			)
			SELECT
				ohb.wallet_address,
				ohb.id,
				ohb.tick,
				ohb.overall_balance,
				ohb.available_balance,
				ohb.block_height
			FROM tempp t
			LEFT JOIN ord_historic_balances ohb ON ohb.id = t.id;`

		//更新受影响的holder的余额（在余额变更历史中找到最新一条变更历史，更新余额到holder表）
		for collectionId, holders := range CollectionIdToAddressList {
			addressList := make([]string, 0)
			for address, _ := range holders {
				addressList = append(addressList, address)
			}
			if len(addressList) == 0 {
				continue
			}
			var records []biz.OrdHistoricBalancesBO
			query := strings.Replace(queryTemplate, "{{CollectionID}}", strconv.Itoa(int(collectionId)), 1)
			err = tx.Raw(query, addressList).Find(&records).Error
			if err != nil {
				return err
			}
			for _, record := range records {
				err = tx.Model(db.HolderPO{}).
					Where("address = ?", record.WalletAddress).
					Where("chain_id = ?", chainId).
					Where("collection_id = ?", collectionId).
					Updates(map[string]interface{}{
						"overall_balance":     record.OverallBalance,
						"available_balance":   record.AvailableBalance,
						"update_time":         time.Now().Unix(),
						"last_block_height":   record.BlockHeight,
						"last_transaction_id": record.TransactionId,
					}).Error
				if err != nil {
					return err
				}
				//updated
				holders[record.WalletAddress] = false
			}

			//未找到余额历史说明holder在回滚后无余额，更新余额为0
			for address, ok := range holders {
				if !ok {
					continue
				}
				err = tx.Model(db.HolderPO{}).
					Where("address = ?", address).
					Where("chain_id = ?", chainId).
					Where("collection_id = ?", collectionId).
					Updates(map[string]interface{}{
						"overall_balance":     0,
						"available_balance":   0,
						"update_time":         time.Now().Unix(),
						"last_block_height":   0,
						"last_transaction_id": 0,
					}).Error
				if err != nil {
					return err
				}

			}

		}

		collectionList := make([]int64, 0)
		for collectionId, _ := range waitForUpdateCollection {
			collectionList = append(collectionList, collectionId)
		}
		// 更新受影响的collection的mint次数，剩余可mint余额，上次mint高度和txid
		query :=
			`SELECT
			collection_id,
			MAX(block_height) AS block_height,
			SUM(amount) AS amount,
			COUNT(*) AS is_used
		FROM ord_event 
		WHERE event_status = 'processed'
    		AND event_type = 1
    		AND collection_id IN (?)
		GROUP BY collection_id;`
		var groupRes []db.OrdEventPO
		err = tx.Raw(query, collectionList).Find(&groupRes).Error
		if err != nil {
			return err
		}
		for _, res := range groupRes {
			err = tx.Model(db.CollectionPO{}).
				Where("id = ?", res.CollectionId).
				Updates(map[string]interface{}{
					"total_minted":        res.Amount,
					"mint_times":          res.IsUsed,
					"last_block_height":   res.BlockHeight,
					"last_transaction_id": "",
					"update_time":         time.Now().Unix(),
				}).Error
			if err != nil {
				return err
			}
		}
		return nil
	}
	return r.data.db.SqlDB.Transaction(callFunction)
}

func (r *OrdEventDAO) UpdateEventStatus(ctx context.Context, eventId int64, collectionId int64, oldStatus, newStatus string) error {
	var eventInfo biz.OrdEventBO
	callFunction := func(tx *gorm.DB) error {
		err := tx.Model(&db.OrdEventPO{}).Where("id = ?", eventId).First(&eventInfo).Error
		if err != nil {
			return err
		}
		if eventInfo.EventStatus != oldStatus {
			return fmt.Errorf("update event %v status failed, old status %s inconsistent", eventInfo, oldStatus)
		}
		return tx.Model(&db.OrdEventPO{}).Where("id = ?", eventId).
			Updates(map[string]interface{}{
				"event_status":  newStatus,
				"collection_id": collectionId,
				"update_time":   time.Now().Unix(),
			}).Error
	}

	if tx, ok := ctx.Value(fmt.Sprintf("tx_%d", eventId)).(*gorm.DB); ok {
		err := callFunction(tx)
		return err
	} else {
		return r.data.db.SqlDB.Transaction(callFunction)
	}
}

func (r *OrdEventDAO) MarkTransferInscribeEventAsUsed(ctx context.Context, relatedOrdTransferId int64) error {
	var eventInfo biz.OrdEventBO
	callFunction := func(tx *gorm.DB) error {
		err := tx.Model(&db.OrdEventPO{}).Where("related_ord_transfer_id = ?", relatedOrdTransferId).First(&eventInfo).Error
		if err != nil {
			return err
		}
		if eventInfo.EventType != biz.TransferInscribe {
			return fmt.Errorf("MarkTransferInscribeEventAsUsed failed, event is not Transfer-Inscribe")
		}
		if eventInfo.EventStatus != "processed" {
			return fmt.Errorf("MarkTransferInscribeEventAsUsed failed, event is not valid")
		}
		if eventInfo.IsUsed == 1 {
			return fmt.Errorf("MarkTransferInscribeEventAsUsed failed, event is used")
		}
		return tx.Model(&db.OrdEventPO{}).Where("related_ord_transfer_id = ?", relatedOrdTransferId).
			Updates(map[string]interface{}{
				"is_used":     1,
				"update_time": time.Now().Unix(),
			}).Error
	}

	if tx, ok := ctx.Value(fmt.Sprintf("tx_%d", relatedOrdTransferId)).(*gorm.DB); ok {
		err := callFunction(tx)
		return err
	} else {
		return r.data.db.SqlDB.Transaction(callFunction)
	}
}

func (r *OrdEventDAO) BatchInsertEvent(ctx context.Context, eventList []biz.OrdEventBO) error {
	return r.data.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&db.OrdEventPO{}).CreateInBatches(eventList, len(eventList)).Error
		return err
	})
}

func (r *OrdEventDAO) GetEventList(
	ctx context.Context,
	chainId string,
	collectionId int64,
	inscriptionId string,
	eventStatus string,
	minBlockHeight int64,
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
	if maxBlockHeight > 0 || minBlockHeight > 0 {
		if maxBlockHeight > 0 && minBlockHeight > 0 {
			query = query.Where("block_height BETWEEN ? AND ?", minBlockHeight, maxBlockHeight)
		} else if maxBlockHeight > 0 {
			query = query.Where("block_height <= ?", maxBlockHeight)
		} else {
			query = query.Where("block_height >= ?", minBlockHeight)
		}
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
