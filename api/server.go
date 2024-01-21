package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	genserver "github.com/spacelift-io/homework-object-storage/gen-server"
	storagegateway "github.com/spacelift-io/homework-object-storage/storage_gateway"
	"go.uber.org/zap"
)

type Server struct {
	Addr string
}

func (s *Server) Run(ctx context.Context) {
	discoveredConfigs := storagegateway.DiscoverMinioInstancesInDocker(ctx)
	minioInstances := storagegateway.InitInstances(discoveredConfigs)
	service := storagegateway.NewService(minioInstances...)
	storageGatewayAPI := StorageGatewayAPI{
		service: service,
	}

	httpServer := http.Server{
		Addr:    s.Addr,
		Handler: handlers.LoggingHandler(os.Stdout, genserver.Handler(&storageGatewayAPI)),
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		}
		if err != nil {
			panic(err.Error())
		}
	}()

	zap.L().Info("server started", zap.String("addr", s.Addr))
	<-ctx.Done()

	// attempt to gracefully shutdown the server
	if ctx.Err() == context.Canceled {
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctxShutdown)
		if err != nil {
			panic(err.Error())
		}
		zap.L().Info("server gracefully stopped")
	}

}
