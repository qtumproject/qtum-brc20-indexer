package main

import (
	"flag"
	"fmt"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/server"
)

var cfg = flag.String("conf", "app/ordinal_indexer/interface/configs/config.yaml", "configs file location")

// main
func main() {
	flag.Parse()
	fmt.Println("[debug] conf file path: " + *cfg)
	Run(*cfg)
}

func Run(cfg string) {

	// 获取实例化服务
	server := initServer(cfg)
	//// 上报链路 trace 数据
	//defer func() {
	//	if err = server.trace.Shutdown(context.Background()); err != nil {
	//		_ = level.Info(server.logger).Log("msg", "shutdown trace provider failed", "err", err)
	//	}
	//}()
	//server.GetCheckInReminderStatus()
	//
	// 启动 http 以及 grpc 服务
	server.RunServer()
	// listen exit server event
	server.HandleExitServer()

}

// SetServer Wire inject app's component
func initServer(cfg string) *server.Server {
	server, err := InitServer(cfg)
	if err != nil {
		panic("run server failed.[ERROR]=>" + err.Error())
	}
	return server
}
