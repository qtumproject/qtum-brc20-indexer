//go:build wireinject
// +build wireinject

package main

import (
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/data"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/server"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/service"
	"github.com/google/wire"
)

// InitServer Inject app's component
func InitServer(cfg string) (*server.Server, error) {
	panic(wire.Build(config.NewConfig, server.SvrProviderSet, data.DtProviderSet, biz.BzProviderSet, service.SvcProviderSet))
	return &server.Server{}, nil
}
