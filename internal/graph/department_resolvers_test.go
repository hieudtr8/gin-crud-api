package graph

import (
	"context"
	"testing"

	"gin-crud-api/internal/database"
	"gin-crud-api/internal/graph/model"
	"gin-crud-api/internal/middleware"
	"gin-crud-api/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDepartmentResolverTest(t *testing.T) (*Resolver, context.Context) {
	// Create test database client
	client := testutil.NewTestEntClient(t)
	t.Cleanup(func() { client.Close() })

	// Create repositories
	deptRepo := database.NewEntDepartmentRepo(client)
	empRepo := database.NewEntEmployeeRepo(client)
	projRepo := database.NewEntProjectRepo(client)

	// Create resolver with dependencies
	resolver := NewResolver(deptRepo, empRepo, projRepo)

	// Create context with request ID (for logging)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	return resolver, ctx
}

// TestCreateDepartment_Success tests successful department creation
func TestCreateDepartment_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	input := model.CreateDepartmentInput{
		Name: "Engineering",
	}

	// Create department
	dept, err := resolver.Mutation().CreateDepartment(ctx, input)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, dept)
	assert.Equal(t, "Engineering", dept.Name)
	assert.NotEmpty(t, dept.ID)

	// Verify UUID is valid
	_, err = uuid.Parse(dept.ID)
	assert.NoError(t, err, "Department ID should be a valid UUID")
}

// TestCreateDepartment_EmptyName tests validation for empty department name
func TestCreateDepartment_EmptyName(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	input := model.CreateDepartmentInput{
		Name: "",
	}

	// Attempt to create department with empty name
	dept, err := resolver.Mutation().CreateDepartment(ctx, input)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, dept)
	assert.Contains(t, err.Error(), "department name is required")
}

// TestUpdateDepartment_Success tests successful department update
func TestUpdateDepartment_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create initial department
	createInput := model.CreateDepartmentInput{Name: "Engineering"}
	created, err := resolver.Mutation().CreateDepartment(ctx, createInput)
	require.NoError(t, err)

	// Update department
	updateInput := model.UpdateDepartmentInput{Name: "Product Engineering"}
	updated, err := resolver.Mutation().UpdateDepartment(ctx, created.ID, updateInput)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Product Engineering", updated.Name)
}

// TestUpdateDepartment_NotFound tests updating non-existent department
func TestUpdateDepartment_NotFound(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	nonExistentID := uuid.New().String()
	input := model.UpdateDepartmentInput{Name: "Engineering"}

	// Attempt to update non-existent department
	dept, err := resolver.Mutation().UpdateDepartment(ctx, nonExistentID, input)

	// Assert not found error
	require.Error(t, err)
	assert.Nil(t, dept)
	assert.Contains(t, err.Error(), "not found")
}

// TestUpdateDepartment_EmptyName tests validation for empty name on update
func TestUpdateDepartment_EmptyName(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create initial department
	createInput := model.CreateDepartmentInput{Name: "Engineering"}
	created, err := resolver.Mutation().CreateDepartment(ctx, createInput)
	require.NoError(t, err)

	// Attempt to update with empty name
	updateInput := model.UpdateDepartmentInput{Name: ""}
	updated, err := resolver.Mutation().UpdateDepartment(ctx, created.ID, updateInput)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "department name is required")
}

// TestDeleteDepartment_Success tests successful department deletion
func TestDeleteDepartment_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create department
	createInput := model.CreateDepartmentInput{Name: "Engineering"}
	created, err := resolver.Mutation().CreateDepartment(ctx, createInput)
	require.NoError(t, err)

	// Delete department
	success, err := resolver.Mutation().DeleteDepartment(ctx, created.ID)

	// Assert success
	require.NoError(t, err)
	assert.True(t, success)

	// Verify department is deleted (GraphQL returns nil for not found)
	deleted, err := resolver.Query().Department(ctx, created.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted, "Deleted department should return nil")
}

// TestDeleteDepartment_NotFound tests deleting non-existent department
func TestDeleteDepartment_NotFound(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	nonExistentID := uuid.New().String()

	// Attempt to delete non-existent department
	success, err := resolver.Mutation().DeleteDepartment(ctx, nonExistentID)

	// Assert not found error
	require.Error(t, err)
	assert.False(t, success)
	assert.Contains(t, err.Error(), "not found")
}

// TestDeleteDepartment_WithEmployees tests cascade deletion
func TestDeleteDepartment_WithEmployees(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create department
	deptInput := model.CreateDepartmentInput{Name: "Engineering"}
	dept, err := resolver.Mutation().CreateDepartment(ctx, deptInput)
	require.NoError(t, err)

	// Create employee in department
	empInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	_, err = resolver.Mutation().CreateEmployee(ctx, empInput)
	require.NoError(t, err)

	// Delete department (should cascade delete employees)
	success, err := resolver.Mutation().DeleteDepartment(ctx, dept.ID)

	// Assert success
	require.NoError(t, err)
	assert.True(t, success)
}

// TestDepartmentQuery_Success tests querying a single department
func TestDepartmentQuery_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create department
	createInput := model.CreateDepartmentInput{Name: "Engineering"}
	created, err := resolver.Mutation().CreateDepartment(ctx, createInput)
	require.NoError(t, err)

	// Query department
	found, err := resolver.Query().Department(ctx, created.ID)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Name, found.Name)
}

// TestDepartmentQuery_NotFound tests querying non-existent department
func TestDepartmentQuery_NotFound(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	nonExistentID := uuid.New().String()

	// Query non-existent department
	dept, err := resolver.Query().Department(ctx, nonExistentID)

	// Assert returns nil for not found (GraphQL convention)
	require.NoError(t, err)
	assert.Nil(t, dept, "Non-existent department should return nil")
}

// TestDepartmentsQuery_Success tests querying all departments
func TestDepartmentsQuery_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create multiple departments
	dept1Input := model.CreateDepartmentInput{Name: "Engineering"}
	dept1, err := resolver.Mutation().CreateDepartment(ctx, dept1Input)
	require.NoError(t, err)

	dept2Input := model.CreateDepartmentInput{Name: "Sales"}
	dept2, err := resolver.Mutation().CreateDepartment(ctx, dept2Input)
	require.NoError(t, err)

	// Query all departments
	departments, err := resolver.Query().Departments(ctx)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, departments)
	assert.GreaterOrEqual(t, len(departments), 2)

	// Verify our departments are in the list
	ids := make(map[string]bool)
	for _, dept := range departments {
		ids[dept.ID] = true
	}
	assert.True(t, ids[dept1.ID], "dept1 should be in the list")
	assert.True(t, ids[dept2.ID], "dept2 should be in the list")
}

// TestDepartmentsQuery_Empty tests querying when no departments exist
func TestDepartmentsQuery_Empty(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Query all departments (empty database)
	departments, err := resolver.Query().Departments(ctx)

	// Assert success with empty list
	require.NoError(t, err)
	require.NotNil(t, departments)
	assert.Empty(t, departments)
}

// TestEmployeesByDepartment_Success tests querying employees by department
func TestEmployeesByDepartment_Success(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create department
	deptInput := model.CreateDepartmentInput{Name: "Engineering"}
	dept, err := resolver.Mutation().CreateDepartment(ctx, deptInput)
	require.NoError(t, err)

	// Create employees
	emp1Input := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	emp1, err := resolver.Mutation().CreateEmployee(ctx, emp1Input)
	require.NoError(t, err)

	emp2Input := model.CreateEmployeeInput{
		Name:         "Jane Smith",
		Email:        "jane@example.com",
		DepartmentID: dept.ID,
	}
	emp2, err := resolver.Mutation().CreateEmployee(ctx, emp2Input)
	require.NoError(t, err)

	// Query department employees
	employees, err := resolver.Query().EmployeesByDepartment(ctx, dept.ID)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, employees)
	assert.Len(t, employees, 2)

	// Verify employee IDs
	ids := make(map[string]bool)
	for _, emp := range employees {
		ids[emp.ID] = true
	}
	assert.True(t, ids[emp1.ID], "emp1 should be in the list")
	assert.True(t, ids[emp2.ID], "emp2 should be in the list")
}

// TestEmployeesByDepartment_Empty tests department with no employees
func TestEmployeesByDepartment_Empty(t *testing.T) {
	resolver, ctx := setupDepartmentResolverTest(t)

	// Create department without employees
	deptInput := model.CreateDepartmentInput{Name: "Engineering"}
	dept, err := resolver.Mutation().CreateDepartment(ctx, deptInput)
	require.NoError(t, err)

	// Query department employees
	employees, err := resolver.Query().EmployeesByDepartment(ctx, dept.ID)

	// Assert success with empty list
	require.NoError(t, err)
	require.NotNil(t, employees)
	assert.Empty(t, employees)
}
