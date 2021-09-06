package broker

import (
	"context"
	"github.com/tsmweb/broker-service/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("[>] starting server")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func(ctx context.Context, fn context.CancelFunc) {
		<-ctx.Done()
		fn()
	}(ctx, stop)

	// Working directory
	//workDir, _ := os.Getwd()
	//config.Load(workDir)
	config.Load("../../")

	//TODO
}
