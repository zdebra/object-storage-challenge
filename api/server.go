package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	genserver "github.com/spacelift-io/homework-object-storage/gen-server"
	storagegateway "github.com/spacelift-io/homework-object-storage/storage_gateway"
)

type Server struct {
	Addr string
}

type MinioInstanceCfg struct {
	AccessKey string
	SecretKey string
	Endpoint  string
}

func (s *Server) Run(ctx context.Context) {
	// todo: put this elsewhere
	cfgs := storagegateway.DiscoverMinioInstancesInDocker(ctx)

	instances := []*storagegateway.StorageInstance{}
	for _, cfg := range cfgs {
		instances = append(instances, storagegateway.NewStorageInstance(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey))
	}

	service := storagegateway.NewService(instances...)
	storageGatewayAPI := StorageGatewayAPI{
		service: service,
	}

	httpServer := http.Server{
		Addr:    s.Addr,
		Handler: genserver.Handler(&storageGatewayAPI),
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

	fmt.Println("server is running on", s.Addr)
	<-ctx.Done()

	// attempt to gracefully shutdown the server
	if ctx.Err() == context.Canceled {
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctxShutdown)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("server gracefully stopped")
	}

}
