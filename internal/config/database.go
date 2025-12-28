package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	Pool        *pgxpool.Pool
	DbConfig    *DatabaseConfig
	DatabaseUrl string
)

func init() {
	DbConfig = LoadDatabaseConfig()
	DatabaseUrl = BuildDatabaseDSN(DbConfig)
}

func Connect() {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(DatabaseUrl)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	// Configure connection pool settings from environment
	config.MaxConns = DbConfig.MaxConnections
	config.MinConns = DbConfig.MinConnections
	config.MaxConnLifetime = DbConfig.MaxConnLifetime
	config.MaxConnIdleTime = DbConfig.MaxConnIdleTime

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Test the connection
	err = Pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Database connection established successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("Database connection pool closed")
	}
}

func RunMigrations() {
	log.Println("Running database migrations...")

	m, err := migrate.New(
		"file://internal/migrations",
		DatabaseUrl,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}

	log.Println("Migrations applied successfully")
}
