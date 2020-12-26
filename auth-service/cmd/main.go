package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tsmweb/auth-service/helper/setting"
	"github.com/tsmweb/go-helper-api/middleware"
	"github.com/urfave/negroni"
	"log"
	"os"
)

func main() {
	log.Println("[>] Starting server")
	workDir, _ := os.Getwd() // working directory
	setting.Load(workDir)

	router := mux.NewRouter()
	profileRouter := InitProfileRouter()
	profileRouter.MakeRouters(router)

	//loginRouter := InitLoginRouter()
	//loginRouter.MakeRouters(router)

	handler := middleware.GZIP(router)
	handler = middleware.CORS(handler)

	nr := negroni.New()
	nr.Use(negroni.NewLogger())
	nr.UseHandler(handler)

	serverPort := setting.ServerPort()
	nr.Run(fmt.Sprintf(":%d", serverPort))
}
