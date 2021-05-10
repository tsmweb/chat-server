package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/chat-service/helper/setting"
	"github.com/tsmweb/go-helper-api/concurrent/executor"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
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
	workDir, _ := os.Getwd()
	setting.Load(workDir)

	log.Println("[>] starting chat server")

	// Executor to perform background processing,
	// limiting resource consumption when executing a collection of jobs.
	exe := executor.New(setting.GoPollSize())
	defer exe.Shutdown()

	chatServer, err := InitChat(setting.Localhost(), exe)
	if err != nil {
		log.Fatalf("[!] error when starting chat: %s\n", err.Error())
	}

	router := mux.NewRouter()
	chatServer.MakeRouters(router)

	log.Println("[>] chat server started")
	log.Println("[>] starting entrypoint")

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	serverPort := setting.ServerPort()
	nr.Run(fmt.Sprintf(":%d", serverPort))
}