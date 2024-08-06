package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBSource   string `env:"DB_SOURCE,required"`
	ServerPort string `env:"SERVER_PORT,required"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := new(Config)

	err = env.Parse(cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return cfg, nil
}
