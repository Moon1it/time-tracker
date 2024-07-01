package main

import (
	"log"
	"time-tracker/internal/config"
	"time-tracker/internal/handlers"
	"time-tracker/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	handler := handlers.NewHandler()

	srv := new(server.Server)
	if err := srv.Run(cfg.ServerPort, handler); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
