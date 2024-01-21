package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spacelift-io/homework-object-storage/api"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)
	appCtx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		zap.L().Info("Received signal to terminate the server. Attempting graceful shutdown...")
		cancel()
	}()

	server := api.Server{Addr: ":3000"}
	server.Run(appCtx)
}
