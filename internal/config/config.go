package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type DatabaseConfig struct {
	DatabaseName     string
	DatabasePort     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseUsername: os.Getenv("DATABASE_USERNAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
	}
}

func BuildDatabaseDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DatabaseUsername,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)
}
