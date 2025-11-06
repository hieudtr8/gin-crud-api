package database

import (
	"gin-crud-api/internal/models"
)

// Repository interfaces define the contract for data access
// These interfaces allow swapping between different storage implementations
// (e.g., in-memory, PostgreSQL, MongoDB, etc.) without changing handler code

// DepartmentRepository defines all operations for managing departments
type DepartmentRepository interface {
	Save(dept *models.Department) error
	FindByID(id string) (*models.Department, error)
	FindAll() ([]*models.Department, error)
	Update(dept *models.Department) error
	Delete(id string) error
}

// EmployeeRepository defines all operations for managing employees
type EmployeeRepository interface {
	Save(emp *models.Employee) error
	FindByID(id string) (*models.Employee, error)
	FindAll() ([]*models.Employee, error)
	Update(emp *models.Employee) error
	Delete(id string) error
	FindByDepartmentID(deptID string) ([]*models.Employee, error)
}