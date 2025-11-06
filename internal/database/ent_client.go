package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gin-crud-api/internal/config"
	"gin-crud-api/internal/ent"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// NewEntClient creates a new EntGo client with PostgreSQL connection
// and automatically runs database migrations
func NewEntClient(cfg *config.DatabaseConfig) (*ent.Client, error) {
	// Build PostgreSQL connection string for lib/pq driver
	// Format: "host=localhost port=5432 user=postgres password=postgres dbname=gin_crud_api sslmode=disable"
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Open database connection using postgres driver
	drv, err := entsql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db := drv.DB()
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetConnMaxLifetime(time.Hour)

	// Create EntGo client with PostgreSQL dialect
	client := ent.NewClient(ent.Driver(drv))

	// Run automatic migrations
	// This creates/updates tables based on schema definitions
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ EntGo client initialized successfully")
	log.Println("✅ Database migrations completed")

	return client, nil
}

// CloseEntClient closes the EntGo client and database connection
func CloseEntClient(client *ent.Client) error {
	if err := client.Close(); err != nil {
		return fmt.Errorf("failed to close EntGo client: %w", err)
	}
	log.Println("✅ EntGo client closed successfully")
	return nil
}
