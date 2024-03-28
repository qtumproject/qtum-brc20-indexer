package server

import (
	"context"
	"fmt"
	v1 "github.com/6block/fox_ordinal/api/ordinal_indexer/interface/v1"
	"github.com/6block/fox_ordinal/app/ordinal_indexer/interface/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

type GinHttpServer struct {
	engin   *gin.Engine
	port    string
	service v1.OrdinalIndexerInterfaceServer
}

func (server *GinHttpServer) Run() (err error) {
	return server.engin.Run(server.port)
}

func (server *GinHttpServer) ShutDown(ctx context.Context) (err error) {
	return nil
}

func (server *GinHttpServer) HandleCheckInReminderStatus(c *gin.Context) {
	//did := c.DefaultQuery("did", "0")
	//registerId := c.DefaultQuery("registerId", "0")
	//
	//req := &v1.GetCheckInReminderStatusReq{
	//	RegisterId: registerId,
	//}
	//reply, err := server.service.GetCheckInReminderStatus(c, req)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, "error")
	//	return
	//}
	//tz := time.FixedZone("CST", 8*3600)
	//time := time.Now().In(tz).Format("2006-01-02 15:04:05")
	//
	//fmt.Printf("【debug】【%s】call GetCustomCheckInReminderConfig result: %t\n", time, status)
	c.JSON(http.StatusOK, nil)
}

// NewHttpServer 实例化 Http 服务
func NewHttpServer(conf *config.Config, service v1.OrdinalIndexerInterfaceServer) HttpServerI {
	engin := gin.Default()
	//auth := engin.Group("/api/v1")
	v1.RegisterOrdinalIndexerInterfaceHTTPServer(engin, service)
	server := &GinHttpServer{
		engin:   engin,
		port:    ":" + conf.Http.Port,
		service: service,
	}
	//auth.GET("/task/checkInReminderStatus", server.HandleCheckInReminderStatus)
	return server
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in the allowed header, don't send the header
	return s, false
}

// RunHttpServer Run http server
func RunHttpServer(conf *config.Config) {
	//mux := runtime.NewServeMux(
	//	// convert header in response(going from gateway) from metadata received.
	//	runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
	//	runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
	//		header := request.Header.Get("Authorization")
	//		// send all the headers received from the client
	//		md := metadata.Pairs("auth", header)
	//		return md
	//	}),
	//	runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
	//		//creating a new HTTTPStatusError with a custom status, and passing error
	//		newError := runtime.HTTPStatusError{
	//			HTTPStatus: 400,
	//			Err:        err,
	//		}
	//		// using default handler to do the rest of heavy lifting of marshaling error and adding headers
	//		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
	//	}))

	server := gin.Default()
	server.Use(gin.Logger())
	// additonal route
	//server.GET("/test", func(c *gin.Context) {
	//	c.String(http.StatusOK, "Ok")
	//})

	//httpServer := &http.Server{
	//	Addr:    config.Http.Port,
	//	Handler: mux,
	//}
	//go func() {
	//	if err := httpServer.ListenAndServe(); err != nil {
	//		fmt.Println("listen http server failed.[ERROR]=>" + err.Error())
	//	}
	//}()
	fmt.Println("Listening http server on port: " + conf.Http.Port)
	go func() {
		err := server.Run(conf.Http.Port)
		if err != nil {
			log.Fatal(err)
		}
	}()

}
