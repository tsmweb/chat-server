package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	router := mux.NewRouter()

	// starts API server
	providers := CreateProvider(ctx)
	if err := providers.ChatRouter(router); err != nil {
		log.Fatalf("[!] error when starting server: %s\n", err.Error())
	}

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	//nr.Run(fmt.Sprintf(":%d", config.ServerPort()))

	log.Fatal(http.ListenAndServeTLS(
		fmt.Sprintf(":%d", config.ServerPort()),
		config.CertSecureFile(),
		config.KeySecureFile(),
		nr,
	))
}
