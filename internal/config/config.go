package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
	MaxConnections   int32
	MinConnections   int32
	MaxConnLifetime  time.Duration
	MaxConnIdleTime  time.Duration
}

type ServerConfig struct {
	Port string
	Host string
}

func LoadServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: getEnvOrDefault("SERVER_PORT", "8080"),
		Host: getEnvOrDefault("SERVER_HOST", "localhost"),
	}
}

func LoadDatabaseConfig() *DatabaseConfig {
	maxConns, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_CONNECTIONS", "30"))
	minConns, _ := strconv.Atoi(getEnvOrDefault("DB_MIN_CONNECTIONS", "5"))
	maxLifetimeHours, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_CONNECTION_LIFETIME_HOURS", "1"))
	maxIdleMinutes, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_CONNECTION_IDLE_TIME_MINUTES", "30"))

	return &DatabaseConfig{
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabasePort:     os.Getenv("DATABASE_PORT"),
		DatabaseUsername: os.Getenv("DATABASE_USERNAME"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		MaxConnections:   int32(maxConns),
		MinConnections:   int32(minConns),
		MaxConnLifetime:  time.Duration(maxLifetimeHours) * time.Hour,
		MaxConnIdleTime:  time.Duration(maxIdleMinutes) * time.Minute,
	}
}

func BuildDatabaseDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUsername,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
