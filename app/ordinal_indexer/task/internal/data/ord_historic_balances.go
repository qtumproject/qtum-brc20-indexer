package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data/db"
	"gorm.io/gorm"
	"strings"
	"time"
)

var _ biz.IOrdHistoricBalancesDAO = (*OrdHistoricBalancesDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var OrdHistoricBalancesKey = func(taskCode string) string {
	return "reward_item_cache_key_" + taskCode
}

type OrdHistoricBalancesDAO struct {
	db            *db.ServerDatabase
	ordEventDao   biz.IOrdinalEventDAO
	collectionDao biz.ICollectionDAO
	ordHolderDao  biz.IHolderDAO
}

func NewOrdHistoricBalancesDAO(
	db *db.ServerDatabase,
	ordEventDao biz.IOrdinalEventDAO,
	ordHolderDao biz.IHolderDAO,
	collectionDao biz.ICollectionDAO,
) biz.IOrdHistoricBalancesDAO {
	return &OrdHistoricBalancesDAO{
		db:            db,
		ordEventDao:   ordEventDao,
		ordHolderDao:  ordHolderDao,
		collectionDao: collectionDao,
	}
}

func (dao *OrdHistoricBalancesDAO) ListHistoricBalances(ctx context.Context, collectionId int64) ([]biz.OrdHistoricBalancesBO, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *OrdHistoricBalancesDAO) GetLastBalance(ctx context.Context, chainId string, walletAddress string, collectionId int64) (biz.OrdHistoricBalancesBO, error) {
	var newestRecord biz.OrdHistoricBalancesBO
	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&db.OrdHistoricBalancesPO{}).
			Where("chain_id = ?", chainId).
			Where("collection_id = ?", collectionId).
			Where("wallet_address = ?", walletAddress).
			Order("block_height desc, id desc").
			Last(&newestRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var collectionInfo db.CollectionPO
			err = tx.Model(&db.CollectionPO{}).Where("id = ?", collectionId).First(&collectionInfo).Error
			if err != nil {
				return err
			}
			newestRecord.WalletAddress = walletAddress
			newestRecord.CollectionId = collectionId
			newestRecord.Tick = collectionInfo.Tick
			newestRecord.ChainId = chainId
			newestRecord.OverallBalance = 0
			newestRecord.AvailableBalance = 0
			return nil
		} else if err != nil {
			return err
		} else {
			return nil
		}
	})
	return newestRecord, err
}

func (dao *OrdHistoricBalancesDAO) AddHistoricBalanceRecord(ctx context.Context, bo biz.OrdHistoricBalancesBO) (biz.OrdHistoricBalancesBO, error) {

	callFunction := func(tx *gorm.DB) error {
		if bo.CreateTime == 0 {
			bo.CreateTime = time.Now().Unix()
			bo.UpdateTime = bo.CreateTime
		}
		var newestRecord biz.OrdHistoricBalancesBO
		err := tx.Model(&db.OrdHistoricBalancesPO{}).
			Where("chain_id = ?", bo.ChainId).
			Where("collection_id = ?", bo.CollectionId).
			Where("wallet_address = ?", bo.WalletAddress).
			Where("pkscript = ?", bo.Pkscript).
			Order("block_height desc, id desc").
			Last(&newestRecord).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bo.OverallBalance = 0
		} else if err != nil {
			return err
		} else {
			bo.OverallBalance = newestRecord.OverallBalance
			bo.AvailableBalance = newestRecord.AvailableBalance
		}
		if bo.ChangeType == 0 {
			bo.OverallBalance -= bo.OverallAmount
			bo.AvailableBalance -= bo.AvailableAmount
			if bo.OverallBalance < 0 || bo.AvailableBalance < 0 {
				return errors.New("add HistoricBalanceRecord fail, insufficient number of balance")
			}

		} else if bo.ChangeType == 1 {
			bo.OverallBalance += bo.OverallAmount
			bo.AvailableBalance += bo.AvailableAmount
		}

		err = tx.Model(&db.OrdHistoricBalancesPO{}).Create(&bo).Error
		if err != nil {
			return err
		}
		var holderInfo biz.HolderBO
		err = tx.Model(db.HolderPO{}).
			Where("address = ?", bo.WalletAddress).
			Where("pkscript = ?", bo.Pkscript).
			Where("chain_id = ?", bo.ChainId).
			Where("collection_id = ?", bo.CollectionId).First(&holderInfo).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			holderInfo = biz.HolderBO{
				Address:      bo.WalletAddress,
				Pkscript:     bo.Pkscript,
				ChainId:      bo.ChainId,
				CollectionId: bo.CollectionId,
				Tick:         bo.Tick,
				CreateTime:   time.Now().Unix(),
				UpdateTime:   time.Now().Unix(),
			}
			err = tx.Model(db.HolderPO{}).Create(&holderInfo).Error
			if err != nil {
				return err
			}
		}
		err = tx.Model(db.HolderPO{}).
			Where("id = ?", holderInfo.ID).
			Updates(map[string]interface{}{
				"overall_balance":   bo.OverallBalance,
				"available_balance": bo.AvailableBalance,
				"update_time":       time.Now().Unix(),
			}).Error
		return err
	}

	if tx, ok := ctx.Value(fmt.Sprintf("tx_%d", bo.EventId)).(*gorm.DB); ok {
		err := callFunction(tx)
		return bo, err
	} else {
		err := dao.db.SqlDB.Transaction(callFunction)
		return bo, err
	}

}

func (dao *OrdHistoricBalancesDAO) HandleDeployInscribeEvent(
	ctx context.Context,
	eventInfo biz.OrdEventBO,
	collectionInfo biz.CollectionBO) error {
	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		newStatus := "processed"
		err := tx.Model(&db.CollectionPO{}).Create(&collectionInfo).Error
		if err != nil {
			newStatus = "invalid"
			fmt.Printf("CreateCollection failed, error: %s\n", err.Error())
		}
		return dao.ordEventDao.UpdateEventStatus(ctx, eventInfo.ID, collectionInfo.ID, eventInfo.EventStatus, newStatus)
	})
	return err
}

func (dao *OrdHistoricBalancesDAO) HandleMintInscribeEvent(
	ctx context.Context,
	eventInfo biz.OrdEventBO,
	collectionInfo biz.CollectionBO,
	historicBalancesRecord biz.OrdHistoricBalancesBO) error {

	historicBalancesRecord.WalletAddress = eventInfo.TargetAddress
	historicBalancesRecord.Pkscript = eventInfo.TargetPkscript

	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {

		updateValues := map[string]interface{}{
			"total_minted":        collectionInfo.TotalMinted,
			"mint_times":          collectionInfo.MintTimes,
			"last_block_height":   collectionInfo.LastBlockHeight,
			"last_transaction_id": collectionInfo.LastTransactionId,
			"update_time":         time.Now().Unix(),
		}
		err := tx.Model(&db.CollectionPO{}).Where("id = ?", collectionInfo.ID).
			Updates(updateValues).Error
		if err != nil {
			fmt.Printf("MintBRC20OrdinalCollection event %v, but failed to UpdateCollectionInfo, error: %s\n", eventInfo, err.Error())
			return err
		}
		ctxWithTxHandle := context.WithValue(ctx, fmt.Sprintf("tx_%d", eventInfo.ID), tx)
		err = dao.ordEventDao.UpdateEventStatus(ctxWithTxHandle, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "processed")
		if err != nil {
			return err
		}
		err = tx.Model(&db.OrdEventPO{}).Where("id = ?", eventInfo.ID).
			Updates(map[string]interface{}{
				"collection_id": eventInfo.CollectionId,
				"update_time":   time.Now().Unix(),
			}).Error
		if err != nil {
			return err
		}

		_, err = dao.AddHistoricBalanceRecord(ctxWithTxHandle, historicBalancesRecord)
		if err != nil {
			fmt.Printf("MintBRC20OrdinalCollection event %v, but failed to AddHistoricBalanceRecord, error: %s\n", eventInfo, err.Error())
		}
		return err
	})
	return err
}

func (dao *OrdHistoricBalancesDAO) HandleTransferInscribeEvent(
	ctx context.Context, eventInfo biz.OrdEventBO, bo biz.OrdHistoricBalancesBO) error {

	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		ctxWithTxHandle := context.WithValue(ctx, fmt.Sprintf("tx_%d", eventInfo.ID), tx)
		_, err := dao.AddHistoricBalanceRecord(ctxWithTxHandle, bo)
		if err != nil {
			fmt.Printf("HandleTransferInscribeEvent event %v, but failed to AddHistoricBalanceRecord, error: %s\n", eventInfo, err.Error())
			return err
		}
		err = tx.Model(&db.OrdEventPO{}).Where("id = ?", eventInfo.ID).
			Updates(map[string]interface{}{
				"collection_id": eventInfo.CollectionId,
				"update_time":   time.Now().Unix(),
			}).Error
		if err != nil {
			return err
		}
		return dao.ordEventDao.UpdateEventStatus(ctxWithTxHandle, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "processed")
	})
	return err

}

func (dao *OrdHistoricBalancesDAO) HandleSentAsFeeTransferEvent(ctx context.Context, eventInfo biz.OrdEventBO, bo biz.OrdHistoricBalancesBO) error {

	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		ctxWithTxHandle := context.WithValue(ctx, fmt.Sprintf("tx_%d", eventInfo.ID), tx)
		_, err := dao.AddHistoricBalanceRecord(ctxWithTxHandle, bo)
		if err != nil {
			fmt.Printf("HandleSentAsFeeTransferEvent event %v, but failed to AddHistoricBalanceRecord, error: %s\n", eventInfo, err.Error())
			return err
		}
		err = tx.Model(&db.OrdEventPO{}).Where("id = ?", eventInfo.ID).
			Updates(map[string]interface{}{
				"collection_id": eventInfo.CollectionId,
				"update_time":   time.Now().Unix(),
			}).Error
		if err != nil {
			return err
		}
		return dao.ordEventDao.UpdateEventStatus(ctxWithTxHandle, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "processed")

	})
	return err
}

func (dao *OrdHistoricBalancesDAO) HandleTransferTransferEvent(ctx context.Context, eventInfo biz.OrdEventBO, bo biz.OrdHistoricBalancesBO) error {

	err := dao.db.SqlDB.Transaction(func(tx *gorm.DB) error {
		ctxWithTransferInscribeTxHandle := context.WithValue(ctx, fmt.Sprintf("tx_%d", eventInfo.IsUsed), tx)
		ctxWithTxHandle := context.WithValue(ctx, fmt.Sprintf("tx_%d", eventInfo.ID), tx)

		err := dao.ordEventDao.MarkTransferInscribeEventAsUsed(ctxWithTransferInscribeTxHandle, eventInfo.IsUsed)
		if err != nil {
			if strings.HasPrefix(err.Error(), "MarkTransferInscribeEventAsUsed failed") {
				return dao.ordEventDao.UpdateEventStatus(
					ctxWithTxHandle, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "invalid")
			}
			return err
		}

		_, err = dao.AddHistoricBalanceRecord(ctxWithTxHandle, bo)
		if err != nil {
			fmt.Printf("TransferBRC20OrdinalCollection event %v, but failed to AddHistoricBalanceRecord, error: %s\n", eventInfo, err.Error())
			return err
		}
		//transfer-transfer, targetAddress overall_balance += amount, available_balance += amount
		bo.AvailableAmount = bo.OverallAmount
		bo.ChangeType = 1 //收入
		bo.WalletAddress = eventInfo.TargetAddress
		bo.Pkscript = eventInfo.TargetPkscript
		_, err = dao.AddHistoricBalanceRecord(ctxWithTxHandle, bo)
		if err != nil {
			fmt.Printf("TransferBRC20OrdinalCollection event %v, but failed to AddHistoricBalanceRecord, error: %s\n", eventInfo, err.Error())
			return err
		}
		err = tx.Model(&db.OrdEventPO{}).Where("id = ?", eventInfo.ID).
			Updates(map[string]interface{}{
				"collection_id": eventInfo.CollectionId,
				"update_time":   time.Now().Unix(),
			}).Error
		if err != nil {
			return err
		}
		return dao.ordEventDao.UpdateEventStatus(ctxWithTxHandle, eventInfo.ID, eventInfo.CollectionId, eventInfo.EventStatus, "processed")

	})

	return err

}
