// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/server"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/service"
)

// Injectors from wire.go:

// InitServer Inject app's component
func InitServer(cfg2 string) (*server.Server, error) {
	configConfig := config.NewConfig(cfg2)
	db := data.NewGormSqlDB(configConfig)
	foxRedis := data.NewFoxRedis(configConfig)
	serverDatabase := data.NewServerDatabase(db, foxRedis)
	dataData, err := data.NewData(serverDatabase)
	if err != nil {
		return nil, err
	}
	iCollectionDAO := data.NewCollectionDAO(dataData)
	collectionService := biz.NewCollectionService(iCollectionDAO)
	iOrdinalEventDAO := data.NewOrdinalEventDAO(dataData)
	iHolderDAO := data.NewHolderDAO(dataData)
	iOrdHistoricBalancesDAO := data.NewOrdHistoricBalancesDAO(serverDatabase, iOrdinalEventDAO, iHolderDAO, iCollectionDAO)
	ordHistoricBalancesService := biz.NewOrdHistoricBalancesService(iOrdHistoricBalancesDAO)
	holderService := biz.NewHolderService(iHolderDAO)
	iDataSourceDAO := data.NewDataSourceDAO(configConfig, dataData)
	iBlockDAO := data.NewBlockDAO(serverDatabase)
	ordinalEventService := biz.NewOrdinalEventService(iOrdinalEventDAO, iDataSourceDAO, iBlockDAO)
	ordinalIndexerTaskServer := service.NewOrdinalIndexerServiceServer(configConfig, collectionService, ordHistoricBalancesService, holderService, ordinalEventService)
	grpcServer := server.NewGrpcServer(ordinalIndexerTaskServer)
	serverServer := server.NewServer(ordinalIndexerTaskServer, grpcServer, configConfig)
	return serverServer, nil
}