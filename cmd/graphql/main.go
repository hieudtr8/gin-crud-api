package main

import (
	"fmt"
	"net/http"
	"os"

	"gin-crud-api/internal/config"
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/graph"
	"gin-crud-api/internal/logger"
	"gin-crud-api/internal/middleware"

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
		// Can't use logger yet, use panic
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// Initialize logger first
	logger.Init(cfg.LogLevel, cfg.LogPretty)
	log := logger.GetLogger()

	log.Info().
		Str("log_level", cfg.LogLevel).
		Bool("pretty", cfg.LogPretty).
		Msg("Application starting")

	// Initialize EntGo client (PostgreSQL) with auto-migrations
	log.Info().
		Str("host", cfg.Database.Host).
		Int("port", cfg.Database.Port).
		Str("database", cfg.Database.DBName).
		Msg("Connecting to PostgreSQL database")

	entClient, err := database.NewEntClient(&cfg.Database)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to database")
	}
	defer database.CloseEntClient(entClient)

	log.Info().Msg("Database connection established successfully")

	// Create repositories
	deptRepo := database.NewEntDepartmentRepo(entClient)
	empRepo := database.NewEntEmployeeRepo(entClient)

	log.Info().Msg("Repositories initialized")

	// Create GraphQL resolver with injected dependencies
	resolver := graph.NewResolver(deptRepo, empRepo)

	// Create GraphQL server with logging middleware
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AroundOperations(middleware.LoggingMiddleware())

	log.Info().Msg("GraphQL server configured with logging middleware")

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

	// Log server startup info
	log.Info().
		Str("playground_url", fmt.Sprintf("http://localhost:%s/", graphqlPort)).
		Str("graphql_url", fmt.Sprintf("http://localhost:%s/query", graphqlPort)).
		Str("schema", "internal/graph/schema.graphql").
		Str("database", "PostgreSQL (via EntGo)").
		Msg("╔════════════════════════════════════════════════════════╗")

	log.Info().Msg("║  GraphQL Server is running!                        ║")
	log.Info().Msgf("║  Playground:  http://localhost:%s/              ║", graphqlPort)
	log.Info().Msgf("║  GraphQL API: http://localhost:%s/query         ║", graphqlPort)
	log.Info().Msg("╚════════════════════════════════════════════════════════╝")

	// Start server
	log.Info().
		Str("address", serverAddr).
		Msg("Starting HTTP server")

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start server")
	}
}
