package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go-backend-valos-id/core/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	connStr := cfg.GetConnectionString()

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Configure connection pool
	pool.Config().MaxConns = 25
	pool.Config().MinConns = 5
	pool.Config().MaxConnLifetime = 5 * time.Minute
	pool.Config().MaxConnIdleTime = 2 * time.Minute
	pool.Config().HealthCheckPeriod = 1 * time.Minute

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return &Database{Pool: pool}, nil
}

func (d *Database) Close() error {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection pool closed")
	}
	return nil
}

// Health check method
func (d *Database) Health() error {
	return d.Pool.Ping(context.Background())
}
