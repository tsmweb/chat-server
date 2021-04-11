package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/tsmweb/user-service/helper/setting"
	"github.com/urfave/negroni"
	"log"
	"os"
)

func main() {
	log.Println("[>] Starting server")
	workDir, _ := os.Getwd() // working directory
	setting.Load(workDir)

	router := mux.NewRouter()
	// Contact
	contactRouter := InitContactRouter()
	contactRouter.MakeRouters(router)
	// Group
	groupRouter := InitGroupRouter()
	groupRouter.MakeRouters(router)

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	serverPort := setting.ServerPort()
	nr.Run(fmt.Sprintf(":%d", serverPort))
}
