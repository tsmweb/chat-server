package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/go-helper-api/observability/event"
	"github.com/urfave/negroni"
)

var (
	debug = flag.String("pprof", "", "address for pprof http")
)

func main() {
	log.Println("[INFO] starting server")
	flag.Parse()

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	if x := *debug; x != "" {
		log.Printf("[INFO] starting pprof server on %s\n", x)
		go func() {
			log.Printf("[ERROR] pprof server error: %v\n", http.ListenAndServe(x, nil))
		}()
	}

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		log.Println("[INFO] stopping chat service...")
		cancelFunc()
	}()

	provider := CreateProvider(ctx)

	// Initializes the service's event producer.
	producerEvents := provider.NewKafkaProducer(config.KafkaEventsTopic())
	if err := event.Init(producerEvents); err != nil {
		log.Fatalf("[ERROR] Could not start events collects. Error: %s\n", err.Error())
	}
	defer event.Close()

	// starts API server
	router := mux.NewRouter()
	if err := provider.ChatRouter(router); err != nil {
		log.Fatalf("[ERROR] error when starting server: %s\n", err.Error())
	}

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
