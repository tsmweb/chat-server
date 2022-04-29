package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"os"
)

func main() {
	log.Println("[>] starting server")

	// Working directory
	workDir, _ := os.Getwd()
	if err := config.Load(workDir); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	provider := CreateProvider()
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
