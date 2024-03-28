package data

import (
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data/db"
	"github.com/6block/fox_ordinal/pkg/cache"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	MysqlDriver = "mysql"
	MaxPageSize = 1000
	MinPageSize = 1
)

// DtProviderSet is data providers.
var DtProviderSet = wire.NewSet(
	NewFoxRedis, NewGormSqlDB, NewServerDatabase, NewHolderDAO,
	NewCollectionDAO, NewOrdHistoricBalancesDAO, NewData, NewOrdinalEventDAO, NewBlockDAO)

type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}
type Data struct {
	db *db.ServerDatabase
}

func NewData(db *db.ServerDatabase) (*Data, error) {
	d := &Data{
		db: db,
	}
	return d, nil
}

func NewServerDatabase(sqlDB *gorm.DB, cacheCli *cache.FoxRedis) *db.ServerDatabase {
	d := &db.ServerDatabase{
		SqlDB:    sqlDB,
		CacheCli: cacheCli,
	}
	return d
}

func NewGormSqlDB(conf *config.Config) *gorm.DB {
	switch conf.Database.ServerDB.Driver {
	case MysqlDriver:
		sqlDB, err := gorm.Open(mysql.Open(conf.Database.ServerDB.Source), &gorm.Config{})
		if err != nil {
			panic("sql connect failed [ERROR]=> " + err.Error())
		}
		return sqlDB
	}
	panic(fmt.Sprintf("sql driver %s not support", conf.Database.ServerDB.Driver))
}

func NewFoxRedis(conf *config.Config) *cache.FoxRedis {
	return cache.NewFoxRedis(conf.Database.Redis.Host+":"+conf.Database.Redis.Port, conf.Database.Redis.Password)
}
