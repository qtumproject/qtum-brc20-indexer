package server

import (
	"context"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/config"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var SvrProviderSet = wire.NewSet(NewGrpcServer, NewHttpServer, NewServer)

type HttpServerI interface {
	Run() (err error)
	ShutDown(ctx context.Context) (err error)
}

// Server 注入微服务所需要的所有依赖
type Server struct {
	Conf       *config.Config
	httpServer HttpServerI
	grpcServer *grpc.Server
}

// NewServer 实例化 Server
func NewServer(
	conf *config.Config,
	httpServer HttpServerI,
	grpcServer *grpc.Server,
) *Server {
	return &Server{
		Conf:       conf,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}

// RunServer 启动 http 以及 grpc 服务
func (s *Server) RunServer() {

	// 启动 grpc 服务
	fmt.Println("Listening grpc server on port: " + s.Conf.Grpc.Port)
	listen, err := net.Listen("tcp", ":"+s.Conf.Grpc.Port)
	if err != nil {
		panic("listen grpc tcp failed.[ERROR]=>" + err.Error())
	}
	go func() {
		if err = s.grpcServer.Serve(listen); err != nil {
			log.Fatal("grpc serve failed", err)
		}
	}()

	// 启动 http 服务
	go func() {
		err = s.httpServer.Run()
		if err != nil {
			log.Fatal("Start fox backend fail: %s", err)
		}
		fmt.Println("foxd start running")
	}()

}

// HandleExitServer Handle service exit event
func (s *Server) HandleExitServer() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.grpcServer.GracefulStop()
	if err := s.httpServer.ShutDown(ctx); err != nil {
		panic("shutdown service failed.[ERROR]=>" + err.Error())
	}
	<-ctx.Done()
	close(ch)
	fmt.Println("Graceful shutdown http & grpc server.")

}

func WaitForSignal(fn func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("receive signal: %s, closing foxd", <-c)
		fn()
		os.Exit(0)
	}()
}
