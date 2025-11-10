package database

import (
	"fmt"

	"gin-crud-api/internal/graph/model"
)

// Repository interfaces define the contract for data access
// These interfaces use GraphQL-generated models as the single source of truth
// This simplifies the codebase by eliminating conversion layers

// ErrNotFound is returned when a record is not found in the database
var ErrNotFound = fmt.Errorf("record not found")

// DepartmentRepository defines all operations for managing departments
type DepartmentRepository interface {
	Save(dept *model.Department) error
	FindByID(id string) (*model.Department, error)
	FindAll() ([]*model.Department, error)
	Update(dept *model.Department) error
	Delete(id string) error
}

// EmployeeRepository defines all operations for managing employees
type EmployeeRepository interface {
	Save(emp *model.Employee) error
	FindByID(id string) (*model.Employee, error)
	FindAll() ([]*model.Employee, error)
	Update(emp *model.Employee) error
	Delete(id string) error
	FindByDepartmentID(deptID string) ([]*model.Employee, error)
}