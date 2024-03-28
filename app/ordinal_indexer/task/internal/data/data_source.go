package data

import (
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
)

var _ biz.IDataSourceDAO = (*QtumDataSourceDAO)(nil) //编译时接口检查，确保userRepo实现了biz.UserRepo定义的接口

var QtumDataSourceCacheKey = func(taskCode string) string {
	return "reward_item_cache_key_" + taskCode
}

func NewDataSourceDAO(conf *config.Config, data *Data) biz.IDataSourceDAO {
	//TODO: 根据配置决定加载哪一个dataSource
	if true {
		return NewQordDataSourceDAO(data, conf)
	} else {
		return NewSimpleDataSourceDAO(conf)
	}
}
