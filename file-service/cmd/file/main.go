package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/file-service/config"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
)

func main() {
	log.Println("[>] starting server")

	// Working directory
	//workDir, _ := os.Getwd()
	//config.Load(workDir)
	if err := config.Load("../../"); err != nil {
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

	serverPort := config.ServerPort()
	nr.Run(fmt.Sprintf(":%d", serverPort))
}