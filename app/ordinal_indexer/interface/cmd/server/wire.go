//go:build wireinject
// +build wireinject

package main

import (
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/biz"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/data"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/server"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/service"
	"github.com/google/wire"
)

// InitServer Inject app's component
func InitServer(cfg string) (*server.Server, error) {
	panic(wire.Build(config.NewConfig, server.SvrProviderSet, data.DtProviderSet, biz.BzProviderSet, service.SvcProviderSet))
	return &server.Server{}, nil
}
