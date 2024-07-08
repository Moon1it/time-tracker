package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource   string
	ServerPort string
}

func getEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	dbSource, err := getEnv("DB_SOURCE")
	if err != nil {
		return nil, err
	}

	serverPort, err := getEnv("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		DBSource:   dbSource,
		ServerPort: serverPort,
	}, nil
}
