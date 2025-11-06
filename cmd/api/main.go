package main

import (
	"gin-crud-api/internal/config"
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/department"
	"gin-crud-api/internal/employee"
	"gin-crud-api/internal/router"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting server with storage type: %s", cfg.Storage)

	var deptRepo database.DepartmentRepository
	var empRepo database.EmployeeRepository
	var cleanup func()

	// Initialize repositories based on storage type
	switch cfg.Storage {
	case "postgres":
		// Initialize PostgreSQL
		db, err := database.NewPostgresDB(
			cfg.Database.DSN(),
			cfg.Database.MaxConns,
			cfg.Database.MinConns,
		)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		cleanup = db.Close
		log.Println("Connected to PostgreSQL successfully")

		// Run migrations
		migrationsPath := filepath.Join(".", "migrations")
		if err := database.RunMigrations(cfg.Database.DSN(), migrationsPath); err != nil {
			log.Printf("Warning: Migration failed: %v", err)
		} else {
			log.Println("Database migrations completed successfully")
		}

		// Create PostgreSQL repositories
		deptRepo = database.NewPostgresDepartmentRepository(db)
		empRepo = database.NewPostgresEmployeeRepository(db)

	case "memory":
		fallthrough
	default:
		// Initialize in-memory storage
		store := database.NewInMemoryStore()
		deptRepo = database.NewDepartmentRepository(store)
		empRepo = database.NewEmployeeRepository(store)
		cleanup = func() { log.Println("Shutting down in-memory storage") }
		log.Println("Using in-memory storage")
	}

	defer cleanup()

	// Initialize handlers (same for both storage types)
	deptHandler := department.NewHandler(deptRepo, empRepo)
	empHandler := employee.NewHandler(empRepo, deptRepo)

	// Setup router
	r := router.Setup(deptHandler, empHandler)

	// Start server
	serverAddr := ":" + cfg.Port
	log.Printf("Starting server on %s...", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}