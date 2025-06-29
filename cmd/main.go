package main

import (
	"eticket-api/config"
	"eticket-api/internal/app"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	server, err := app.New(config)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		address := fmt.Sprintf(":%d", config.Server.Port)
		if err := server.App().Run(address); err != nil {
			log.Fatalf("server run failed: %v", err)
		}
	}()

	<-quit
}
