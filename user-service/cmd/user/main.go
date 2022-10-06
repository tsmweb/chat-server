package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/go-helper-api/observability/event"
	"github.com/tsmweb/user-service/config"
	"github.com/urfave/negroni"
)

func main() {
	log.Println("[INFO] Starting server")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	provider := CreateProvider()

	// Initializes the service's event producer.
	producerEvents := provider.NewKafkaProducer(config.KafkaEventsTopic())
	if err := event.Init(producerEvents); err != nil {
		log.Fatalf("[ERROR] Could not start events collects. Error: %s", err.Error())
	}
	defer event.Close()

	// Configure the routes.
	router := mux.NewRouter()
	provider.ContactRouter(router)
	provider.GroupRouter(router)

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
