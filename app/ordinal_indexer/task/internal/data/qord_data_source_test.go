package data

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"testing"
)

func TestBatchFetchOrdinalEvent(t *testing.T) {
	configConfig := config.NewConfig("../../configs/config.yaml")
	db := NewGormSqlDB(configConfig)
	foxRedis := NewFoxRedis(configConfig)
	serverDatabase := NewServerDatabase(db, foxRedis)
	data, _ := NewData(serverDatabase)
	dao := &OrdEventDAO{
		data: data,
	}
	res := dao.ReorgEvent(context.Background(), "qtum", 3696520)

	fmt.Println(res.Error())
}
