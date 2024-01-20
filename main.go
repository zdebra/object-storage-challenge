package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spacelift-io/homework-object-storage/api"
)

func main() {
	dddd()
	appCtx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Printf("received signal to terminate the server\nattempting graceful shutdown...\n")
		cancel()
	}()

	server := api.Server{Addr: ":3000"}
	server.Run(appCtx)
}

func dddd() {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(context.Background())

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, ctr := range containers {
		fmt.Printf("%s %s\n", ctr.ID, ctr.Image)
	}
}
