package main

import (
	"context"
	"github.com/tsmweb/broker-service/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("[>] start broker service")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func(ctx context.Context, fn context.CancelFunc) {
		<-ctx.Done()
		log.Println("[>] stop broker service")
		fn()
	}(ctx, stop)

	// start broker service
	provider := CreateProvider(ctx)
	brokerService := provider.BrokerProvider()
	brokerService.Start()
}
