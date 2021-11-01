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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func(ctx context.Context, fn context.CancelFunc) {
		<-ctx.Done()
		log.Println("[>] stop broker service")
		fn()
	}(ctx, stop)

	// Working directory
	workDir, _ := os.Getwd()
	config.Load(workDir)
	//config.Load("../../")

	// start broker service
	provider := CreateProvider(ctx)
	brokerService := provider.BrokerProvider()
	brokerService.Start()
}
