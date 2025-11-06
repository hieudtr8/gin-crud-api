package main

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/department"
	"gin-crud-api/internal/employee"
	"gin-crud-api/internal/router"
	"log"
)

func main() {
	// Initialize shared in-memory storage
	store := database.NewInMemoryStore()

	// Create repositories with shared store
	deptRepo := database.NewDepartmentRepository(store)
	empRepo := database.NewEmployeeRepository(store)

	// Initialize handlers with dependency injection
	deptHandler := department.NewHandler(deptRepo, empRepo) // Needs empRepo for cascade delete
	empHandler := employee.NewHandler(empRepo, deptRepo)    // Needs deptRepo for validation

	// Setup router with handlers
	r := router.Setup(deptHandler, empHandler)

	// Start server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}