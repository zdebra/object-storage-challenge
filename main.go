package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spacelift-io/homework-object-storage/api"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Printf("received signal to terminate the server\nattempting graceful shutdown...\n")
		cancel()
	}()

	server := api.Server{Addr: ":8080"}
	server.Run(appCtx)
}
