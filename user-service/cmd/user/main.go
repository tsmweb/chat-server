package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/config"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("[>] Starting server")

	// Working directory
	workDir, _ := os.Getwd()
	config.Load(workDir)
	//config.Load("../../")

	router := mux.NewRouter()

	provider := CreateProvider()
	provider.ContactRouter(router)
	provider.GroupRouter(router)

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
