package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gin-crud-api/internal/config"
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional)
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize EntGo client (PostgreSQL) with auto-migrations
	log.Println("Connecting to PostgreSQL database...")
	entClient, err := database.NewEntClient(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseEntClient(entClient)

	// Create repositories
	deptRepo := database.NewEntDepartmentRepo(entClient)
	empRepo := database.NewEntEmployeeRepo(entClient)

	// Create GraphQL resolver with injected dependencies
	resolver := graph.NewResolver(deptRepo, empRepo)

	// Create GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// GraphQL Playground at root path "/"
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	// GraphQL endpoint at "/query"
	http.Handle("/query", srv)

	// Determine GraphQL port (default: 8081)
	graphqlPort := os.Getenv("GRAPHQL_PORT")
	if graphqlPort == "" {
		graphqlPort = "8081"
	}

	serverAddr := fmt.Sprintf(":%s", graphqlPort)

	log.Printf("╔═══════════════════════════════════════════════════════════╗")
	log.Printf("║  GraphQL Server is running!                               ║")
	log.Printf("╠═══════════════════════════════════════════════════════════╣")
	log.Printf("║  Playground:  http://localhost:%s/                     ║", graphqlPort)
	log.Printf("║  GraphQL API: http://localhost:%s/query                ║", graphqlPort)
	log.Printf("╠═══════════════════════════════════════════════════════════╣")
	log.Printf("║  Database:    PostgreSQL (via EntGo)                      ║")
	log.Printf("║  Schema:      internal/graph/schema.graphql               ║")
	log.Printf("╚═══════════════════════════════════════════════════════════╝")

	// Start server
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
