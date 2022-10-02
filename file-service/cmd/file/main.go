package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/go-helper-api/observability/event"
	"github.com/tsmweb/go-helper-api/observability/metric"
	"github.com/urfave/negroni"
)

func main() {
	log.Println("[INFO] starting server")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	provider := CreateProvider()

	// Collect service metrics.
	producerMetrics := provider.NewKafkaProducer(config.KafkaMetricsTopic())
	err := metric.Start(config.HostID(), config.MetricsSendInterval(), producerMetrics)
	if err != nil {
		log.Fatalf("[ERROR] Could not start metrics collects. Error: %s\n", err.Error())
	}
	defer metric.Stop()

	// Initializes the service's event producer.
	producerEvents := provider.NewKafkaProducer(config.KafkaEventsTopic())
	if err = event.Init(producerEvents); err != nil {
		log.Fatalf("[ERROR] Could not start events collects. Error: %s\n", err.Error())
	}
	defer event.Close()

	// Configure the routes.
	router := mux.NewRouter()
	provider.UserRouter(router)
	provider.GroupRouter(router)
	provider.MediaRouter(router)

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	nr.Run(fmt.Sprintf(":%d", config.ServerPort()))

	//log.Fatal(http.ListenAndServeTLS(
	//	fmt.Sprintf(":%d", config.ServerPort()),
	//	config.CertSecureFile(),
	//	config.KeySecureFile(),
	//	nr,
	//))
}
