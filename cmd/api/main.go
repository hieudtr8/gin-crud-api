package main

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/department"
	"gin-crud-api/internal/employee"
	"gin-crud-api/internal/router"
	"log"
)

func main() {
	// 1. Khởi tạo Database Store
	// InMemoryStore is now just the shared storage backend
	store := database.NewInMemoryStore()

	// 2. Khởi tạo Repositories
	// Create separate repository instances that share the same store
	deptRepo := database.NewDepartmentRepository(store)
	empRepo := database.NewEmployeeRepository(store)

	// 3. Khởi tạo Handlers (Inject Repositories)
	deptHandler := department.NewHandler(deptRepo)
	empHandler := employee.NewHandler(empRepo, deptRepo) // (cần cả 2 repo)

	// 4. Khởi tạo Router (Inject Handlers)
	r := router.Setup(deptHandler, empHandler)

	// 5. Chạy Server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}