package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig   DBConfig
	ServerPort string
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
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

	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	serverPort, err := getEnv("SERVER_PORT")
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfig: DBConfig{
			Host:     dbHost,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
			Port:     dbPort,
		},
		ServerPort: serverPort,
	}, nil
}
