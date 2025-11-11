package main

import (
	"gin-crud-api/internal/config"
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/legacy/rest/department"
	"gin-crud-api/internal/legacy/rest/employee"
	"gin-crud-api/internal/legacy/rest/router"
	"log"
	"os"
)

func main() {
	// Determine environment (dev, prod, test)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev" // Default to development
	}

	// Load configuration from YAML and environment variables
	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("🚀 Starting Gin CRUD API with EntGo ORM")
	log.Printf("🌍 Environment: %s", env)
	log.Printf("📊 Database: PostgreSQL at %s:%d", cfg.Database.Host, cfg.Database.Port)

	// Initialize EntGo client (automatically runs migrations)
	entClient, err := database.NewEntClient(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize EntGo client: %v", err)
	}
	defer func() {
		if err := database.CloseEntClient(entClient); err != nil {
			log.Printf("Error closing EntGo client: %v", err)
		}
	}()

	// Create repositories using EntGo
	deptRepo := database.NewEntDepartmentRepo(entClient)
	empRepo := database.NewEntEmployeeRepo(entClient)

	// Initialize handlers
	deptHandler := department.NewHandler(deptRepo, empRepo)
	empHandler := employee.NewHandler(empRepo, deptRepo)

	// Setup router
	r := router.Setup(deptHandler, empHandler)

	// Start server
	serverAddr := ":" + cfg.Server.RESTPort
	log.Printf("🌐 Server starting on http://localhost%s", serverAddr)
	log.Printf("📝 API endpoints available at http://localhost%s/api/v1", serverAddr)
	log.Printf("💚 Health check at http://localhost%s/health", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}