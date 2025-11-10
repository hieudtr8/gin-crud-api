package testutil

import (
	"context"
	"testing"

	"entgo.io/ent/dialect"
	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/ent/enttest"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver for in-memory tests
)

// NewTestEntClient creates a new EntGo client with in-memory SQLite database
// This is perfect for fast unit and integration tests
func NewTestEntClient(t *testing.T) *ent.Client {
	// Create in-memory SQLite database
	// enttest.Open automatically:
	// 1. Creates the database
	// 2. Runs migrations
	// 3. Cleans up on test completion
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")

	return client
}

// SeedTestDepartment creates a test department with given name
// Returns the created department for use in tests
func SeedTestDepartment(t *testing.T, client *ent.Client, name string) *ent.Department {
	dept, err := client.Department.
		Create().
		SetID(uuid.New()).
		SetName(name).
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to seed test department: %v", err)
	}
	return dept
}

// SeedTestEmployee creates a test employee with given name, email, and department
// Returns the created employee for use in tests
func SeedTestEmployee(t *testing.T, client *ent.Client, name, email string, deptID uuid.UUID) *ent.Employee {
	emp, err := client.Employee.
		Create().
		SetID(uuid.New()).
		SetName(name).
		SetEmail(email).
		SetDepartmentID(deptID).
		Save(context.Background())
	if err != nil {
		t.Fatalf("Failed to seed test employee: %v", err)
	}
	return emp
}

// SeedMultipleDepartments creates multiple test departments
// Useful for testing FindAll and other list operations
func SeedMultipleDepartments(t *testing.T, client *ent.Client, names []string) []*ent.Department {
	var departments []*ent.Department
	for _, name := range names {
		dept := SeedTestDepartment(t, client, name)
		departments = append(departments, dept)
	}
	return departments
}

// SeedMultipleEmployees creates multiple test employees for a department
// Useful for testing cascade delete and FindByDepartmentID
func SeedMultipleEmployees(t *testing.T, client *ent.Client, deptID uuid.UUID, count int) []*ent.Employee {
	var employees []*ent.Employee
	for i := 0; i < count; i++ {
		email := uuid.New().String() + "@test.com" // Unique email
		emp := SeedTestEmployee(t, client, "Test Employee", email, deptID)
		employees = append(employees, emp)
	}
	return employees
}
