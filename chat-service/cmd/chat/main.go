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
	"github.com/tsmweb/go-helper-api/observability/metric"
	"github.com/urfave/negroni"
)

var (
	debug = flag.String("pprof", "", "address for pprof http")
)

func main() {
	log.Println("[>] starting server")
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
		log.Printf("[>] starting pprof server on %s", x)
		go func() {
			log.Printf("[!] pprof server error: %v", http.ListenAndServe(x, nil))
		}()
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func(ctx context.Context, fn context.CancelFunc) {
		<-ctx.Done()
		fn()
	}(ctx, stop)

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

	// starts API server
	router := mux.NewRouter()
	if err := provider.ChatRouter(router); err != nil {
		log.Fatalf("[!] error when starting server: %s\n", err.Error())
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
