package data

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"testing"
	"time"
)

func BenchmarkAdd(b *testing.B) {
	configConfig := config.NewConfig("../../configs/config.yaml")
	db := NewGormSqlDB(configConfig)
	foxRedis := NewFoxRedis(configConfig)
	serverDatabase := NewServerDatabase(db, foxRedis)
	data, _ := NewData(serverDatabase)
	ordEventDao := NewOrdinalEventDAO(data)
	ordHolderDao := NewHolderDAO(data)

	collectionDao := NewCollectionDAO(data)

	dao := NewOrdHistoricBalancesDAO(serverDatabase, ordEventDao, ordHolderDao, collectionDao)

	res, _ := dao.AddHistoricBalanceRecord(context.Background(), biz.OrdHistoricBalancesBO{
		ChainId:         "qtum",
		EventId:         60,
		ChangeType:      1,
		TransactionHash: "dfwq",
		TransactionId:   "dgww",
		BlockHeight:     3547000,
		CreateTime:      time.Now().Unix(),
		UpdateTime:      time.Now().Unix(),
		WalletAddress:   "qQvET4q6ojnxpEmu9M4dLaQPZ4idaVSYGT",
		Pkscript:        "76a91450ba59196c614b68b4810f2801c4d9ca9ab7c66988ac",
		CollectionId:    1,
		Tick:            "qtum",
		OverallAmount:   0.2,
		AvailableAmount: 0.2,
	})
	fmt.Println(res)

}
