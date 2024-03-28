package server

import (
	"context"
	"fmt"
	task_v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/task/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/config"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/task/internal/service"
	"github.com/google/wire"
	"github.com/robfig/cron"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var SvrProviderSet = wire.NewSet(NewServer, NewGrpcServer)

// Server 注入微服务所需要的所有依赖
type Server struct {
	ordinalIndexerTaskServer *service.OrdinalIndexerTaskServer
	grpcServer               *grpc.Server
	conf                     *config.Config
	cronTab                  *cron.Cron
}

// NewServer 实例化 Server
func NewServer(
	service *service.OrdinalIndexerTaskServer,
	grpcServer *grpc.Server,
	conf *config.Config,
) *Server {
	cronTab := cron.New()
	//TODO: 根据配置决定需要开启的定时任务
	cronTab.AddFunc("0 */2 * * * *", service.ProcessOrdinalEventJob) //run every 5 mins
	cronTab.AddFunc("0 */2 * * * *", service.SyncOrdinalEventJob)    //run every 5 mins
	return &Server{
		ordinalIndexerTaskServer: service,
		cronTab:                  cronTab,
		grpcServer:               grpcServer,
		conf:                     conf,
	}
}

// NewGrpcServer 实例化 Grpc 服务
func NewGrpcServer(server *service.OrdinalIndexerTaskServer) *grpc.Server {
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
	task_v1.RegisterOrdinalIndexerTaskServer(grpcServer, server)
	return grpcServer
}

// RunServer 启动
func (s *Server) RunServer() {
	//s.ordinalIndexerTaskServer.SyncOrdinalEventJob()
	//s.ordinalIndexerTaskServer.ProcessOrdinalEventJob()
	s.cronTab.Start()

	// 启动 grpc 服务
	fmt.Println("Listening grpc server on port: " + s.conf.Grpc.Port)
	listen, err := net.Listen("tcp", ":"+s.conf.Grpc.Port)
	if err != nil {
		panic("listen grpc tcp failed.[ERROR]=>" + err.Error())
	}
	go func() {
		if err = s.grpcServer.Serve(listen); err != nil {
			log.Fatal("grpc serve failed", err)
		}
	}()
}

// HandleExitServer Handle service exit event
func (s *Server) HandleExitServer() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	s.cronTab.Stop()
	defer cancel()
	<-ctx.Done()
	close(ch)
	fmt.Println("Graceful shutdown http & grpc server.")

}
