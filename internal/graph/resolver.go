package graph

import (
	"gin-crud-api/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DeptRepo database.DepartmentRepository
	EmpRepo  database.EmployeeRepository
}

// NewResolver creates a new resolver with injected dependencies
func NewResolver(deptRepo database.DepartmentRepository, empRepo database.EmployeeRepository) *Resolver {
	return &Resolver{
		DeptRepo: deptRepo,
		EmpRepo:  empRepo,
	}
}
