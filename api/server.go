package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	genserver "github.com/spacelift-io/homework-object-storage/gen-server"
)

type Server struct {
	Addr string
}

func (s *Server) Run(ctx context.Context) {
	storageGatewayAPI := StorageGatewayAPI{}

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
