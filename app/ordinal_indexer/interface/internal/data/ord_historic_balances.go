package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data/db"
	"gorm.io/gorm"
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
			newestRecord.WalletAddress = walletAddress
			newestRecord.CollectionId = collectionId
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
		err = tx.Model(db.HolderPO{}).
			Where("address = ?", bo.WalletAddress).
			Where("chain_id = ?", bo.ChainId).
			Where("collection_id = ?", bo.CollectionId).
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
