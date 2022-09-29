package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/tsmweb/broker-service/config"
	"github.com/tsmweb/go-helper-api/observability/event"
	"github.com/tsmweb/go-helper-api/observability/metric"
)

func main() {
	fmt.Println("[i] start broker service")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		fmt.Println("[i] stopping broker service...")
		cancelFunc()
	}()

	provider := CreateProvider(ctx)

	// Collect service metrics.
	producerMetrics := provider.NewKafkaProducer(config.KafkaMetricsTopic())
	err := metric.Start(config.HostID(), config.MetricsSendInterval(), producerMetrics)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[!] Could not start metrics collects. Error: %v", err)
	} else {
		defer metric.Stop()
	}

	// Initializes the service's event producer.
	producerEvents := provider.NewKafkaProducer(config.KafkaEventsTopic())
	if err = event.Init(producerEvents); err != nil {
		fmt.Fprintf(os.Stderr, "[!] Could not start events collects. Error: %v", err)
	} else {
		defer event.Close()
	}

	// Start broker service
	brokerService := provider.BrokerProvider()
	brokerService.Start()
}
