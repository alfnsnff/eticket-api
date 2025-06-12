package main

import (
	"eticket-api/internal/app"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serv, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%d", serv.Config().Server.Port)
		if err := serv.App().Run(addr); err != nil {
			log.Fatalf("server run failed: %v", err)
		}
	}()

	<-quit
}
