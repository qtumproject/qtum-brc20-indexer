package server

import (
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// NewGrpcServer 实例化 Grpc 服务
func NewGrpcServer(server v1.OrdinalIndexerInterfaceServer) *grpc.Server {
	grpcServer := grpc.NewServer(
	//grpc.ChainStreamInterceptor(
	//	// otel 链路追踪
	//	otelgrpc.StreamServerInterceptor(),
	//),
	//grpc.ChainUnaryInterceptor(
	//	// otel 链路追踪
	//	otelgrpc.UnaryServerInterceptor(),
	//	// PGV 中间件
	//	Jgrpc_pgv_interceptor.ValidationUnaryInterceptor,
	//),
	)
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())
	v1.RegisterOrdinalIndexerInterfaceServer(grpcServer, server)
	return grpcServer

}
