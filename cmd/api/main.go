package main

import (
	"log"
	_ "time-tracker/docs"
	"time-tracker/internal/config"
	repository "time-tracker/internal/db/sqlc"
	"time-tracker/internal/handler"
	"time-tracker/internal/server"
	"time-tracker/internal/service"
	"time-tracker/pkg/database"
)

// @title Time Tracker API
// @version 1.0
// @description API Server for Time Tracker Application
// @host localhost:8000
// @BasePath /api
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	pgxPool, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	defer pgxPool.Close()

	newRepository := repository.New(pgxPool)
	newService := service.NewService(newRepository)
	newHandler := handler.NewHandler(newService)

	srv := new(server.Server)
	if err := srv.Run(cfg.ServerPort, newHandler); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
