package main

import (
	"log"
	"time-tracker/internal/config"
	"time-tracker/internal/handler"
	"time-tracker/internal/repository"
	"time-tracker/internal/server"
	"time-tracker/internal/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	srv := new(server.Server)
	if err := srv.Run(cfg.ServerPort, handler); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
