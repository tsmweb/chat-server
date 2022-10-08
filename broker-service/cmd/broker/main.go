package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/go-helper-api/observability/event"
)

func main() {
	log.Println("[INFO] start broker service")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		log.Println("[INFO] stopping broker service...")
		cancelFunc()
	}()

	provider := CreateProvider(ctx)

	// Initializes the service's event producer.
	producerEvents := provider.NewKafkaProducer(config.KafkaEventsTopic())
	if err := event.Init(producerEvents); err != nil {
		log.Fatalf("[ERROR] Could not start events collects. Error: %s\n", err.Error())
	}
	defer event.Close()

	// Start broker service
	brokerService := provider.BrokerProvider()
	brokerService.Start()
}
