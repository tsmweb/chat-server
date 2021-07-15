package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/config"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"syscall"
)

var (
	debug = flag.String("pprof", "", "address for pprof http")
)

func main() {
	log.Println("[>] starting server")
	flag.Parse()

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

	// Working directory
	//workDir, _ := os.Getwd()
	config.Load("../../")

	// starts API server
	log.Println("[>] starting chat server")
	initWebServer()
}

func initWebServer() {
	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	executor := executor.New(config.GoPoolSize())
	defer executor.Shutdown()

	chatRouter, err := InitChatRouter(executor)
	if err != nil {
		log.Fatalf("[!] error when starting chat: %s\n", err.Error())
	}

	router := mux.NewRouter()
	chatRouter.MakeRouters(router)

	log.Println("[>] chat server started")
	log.Println("[>] starting entrypoint")

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	serverPort := config.ServerPort()
	nr.Run(fmt.Sprintf(":%d", serverPort))
}