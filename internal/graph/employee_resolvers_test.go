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

func setupEmployeeResolverTest(t *testing.T) (*Resolver, context.Context, *model.Department) {
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

	// Create a test department for employees
	deptInput := model.CreateDepartmentInput{Name: "Engineering"}
	dept, err := resolver.Mutation().CreateDepartment(ctx, deptInput)
	require.NoError(t, err)

	return resolver, ctx, dept
}

// TestCreateEmployee_Success tests successful employee creation
func TestCreateEmployee_Success(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	input := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}

	// Create employee
	emp, err := resolver.Mutation().CreateEmployee(ctx, input)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, emp)
	assert.Equal(t, "John Doe", emp.Name)
	assert.Equal(t, "john@example.com", emp.Email)
	assert.Equal(t, dept.ID, emp.DepartmentID)
	assert.NotEmpty(t, emp.ID)

	// Verify UUID is valid
	_, err = uuid.Parse(emp.ID)
	assert.NoError(t, err, "Employee ID should be a valid UUID")
}

// TestCreateEmployee_EmptyName tests validation for empty employee name
func TestCreateEmployee_EmptyName(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	input := model.CreateEmployeeInput{
		Name:         "",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}

	// Attempt to create employee with empty name
	emp, err := resolver.Mutation().CreateEmployee(ctx, input)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, emp)
	assert.Contains(t, err.Error(), "name is required")
}

// TestCreateEmployee_EmptyEmail tests validation for empty employee email
func TestCreateEmployee_EmptyEmail(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	input := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "",
		DepartmentID: dept.ID,
	}

	// Attempt to create employee with empty email
	emp, err := resolver.Mutation().CreateEmployee(ctx, input)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, emp)
	assert.Contains(t, err.Error(), "email is required")
}

// TestCreateEmployee_InvalidEmail tests validation for invalid email format
func TestCreateEmployee_InvalidEmail(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	testCases := []struct {
		name  string
		email string
	}{
		{"Missing @", "invalidemail.com"},
		{"Missing domain", "user@"},
		{"Missing username", "@example.com"},
		{"Spaces", "user @example.com"},
		{"Invalid TLD", "user@example.c"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := model.CreateEmployeeInput{
				Name:         "John Doe",
				Email:        tc.email,
				DepartmentID: dept.ID,
			}

			// Attempt to create employee with invalid email
			emp, err := resolver.Mutation().CreateEmployee(ctx, input)

			// Assert validation error
			require.Error(t, err)
			assert.Nil(t, emp)
			assert.Contains(t, err.Error(), "invalid email format")
		})
	}
}

// TestCreateEmployee_DepartmentNotFound tests creating employee with non-existent department
func TestCreateEmployee_DepartmentNotFound(t *testing.T) {
	resolver, ctx, _ := setupEmployeeResolverTest(t)

	input := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: uuid.New().String(),
	}

	// Attempt to create employee with non-existent department
	emp, err := resolver.Mutation().CreateEmployee(ctx, input)

	// Assert error
	require.Error(t, err)
	assert.Nil(t, emp)
	assert.Contains(t, err.Error(), "department not found")
}

// TestUpdateEmployee_Success tests successful employee update
func TestUpdateEmployee_Success(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create initial employee
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Update employee
	updateInput := model.UpdateEmployeeInput{
		Name:         "Jane Doe",
		Email:        "jane@example.com",
		DepartmentID: dept.ID,
	}
	updated, err := resolver.Mutation().UpdateEmployee(ctx, created.ID, updateInput)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Jane Doe", updated.Name)
	assert.Equal(t, "jane@example.com", updated.Email)
	assert.Equal(t, dept.ID, updated.DepartmentID)
}

// TestUpdateEmployee_NotFound tests updating non-existent employee
func TestUpdateEmployee_NotFound(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	nonExistentID := uuid.New().String()
	input := model.UpdateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}

	// Attempt to update non-existent employee
	emp, err := resolver.Mutation().UpdateEmployee(ctx, nonExistentID, input)

	// Assert not found error
	require.Error(t, err)
	assert.Nil(t, emp)
	assert.Contains(t, err.Error(), "not found")
}

// TestUpdateEmployee_EmptyName tests validation for empty name on update
func TestUpdateEmployee_EmptyName(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create initial employee
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Attempt to update with empty name
	updateInput := model.UpdateEmployeeInput{
		Name:         "",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	updated, err := resolver.Mutation().UpdateEmployee(ctx, created.ID, updateInput)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "name is required")
}

// TestUpdateEmployee_InvalidEmail tests validation for invalid email on update
func TestUpdateEmployee_InvalidEmail(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create initial employee
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Attempt to update with invalid email
	updateInput := model.UpdateEmployeeInput{
		Name:         "John Doe",
		Email:        "invalid-email",
		DepartmentID: dept.ID,
	}
	updated, err := resolver.Mutation().UpdateEmployee(ctx, created.ID, updateInput)

	// Assert validation error
	require.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "invalid email format")
}

// TestUpdateEmployee_ChangeDepartment tests updating employee's department
func TestUpdateEmployee_ChangeDepartment(t *testing.T) {
	resolver, ctx, dept1 := setupEmployeeResolverTest(t)

	// Create second department
	dept2Input := model.CreateDepartmentInput{Name: "Sales"}
	dept2, err := resolver.Mutation().CreateDepartment(ctx, dept2Input)
	require.NoError(t, err)

	// Create employee in first department
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept1.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Update employee to second department
	updateInput := model.UpdateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept2.ID,
	}
	updated, err := resolver.Mutation().UpdateEmployee(ctx, created.ID, updateInput)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, dept2.ID, updated.DepartmentID)

	// Verify employee is in new department
	emps, err := resolver.Query().EmployeesByDepartment(ctx, dept2.ID)
	require.NoError(t, err)
	assert.Len(t, emps, 1)
	assert.Equal(t, updated.ID, emps[0].ID)
}

// TestDeleteEmployee_Success tests successful employee deletion
func TestDeleteEmployee_Success(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create employee
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Delete employee
	success, err := resolver.Mutation().DeleteEmployee(ctx, created.ID)

	// Assert success
	require.NoError(t, err)
	assert.True(t, success)

	// Verify employee is deleted (GraphQL returns nil for not found)
	deleted, err := resolver.Query().Employee(ctx, created.ID)
	require.NoError(t, err)
	assert.Nil(t, deleted, "Deleted employee should return nil")
}

// TestDeleteEmployee_NotFound tests deleting non-existent employee
func TestDeleteEmployee_NotFound(t *testing.T) {
	resolver, ctx, _ := setupEmployeeResolverTest(t)

	nonExistentID := uuid.New().String()

	// Attempt to delete non-existent employee
	success, err := resolver.Mutation().DeleteEmployee(ctx, nonExistentID)

	// Assert not found error
	require.Error(t, err)
	assert.False(t, success)
	assert.Contains(t, err.Error(), "not found")
}

// TestEmployeeQuery_Success tests querying a single employee
func TestEmployeeQuery_Success(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create employee
	createInput := model.CreateEmployeeInput{
		Name:         "John Doe",
		Email:        "john@example.com",
		DepartmentID: dept.ID,
	}
	created, err := resolver.Mutation().CreateEmployee(ctx, createInput)
	require.NoError(t, err)

	// Query employee
	found, err := resolver.Query().Employee(ctx, created.ID)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Name, found.Name)
	assert.Equal(t, created.Email, found.Email)
	assert.Equal(t, created.DepartmentID, found.DepartmentID)
}

// TestEmployeeQuery_NotFound tests querying non-existent employee
func TestEmployeeQuery_NotFound(t *testing.T) {
	resolver, ctx, _ := setupEmployeeResolverTest(t)

	nonExistentID := uuid.New().String()

	// Query non-existent employee
	emp, err := resolver.Query().Employee(ctx, nonExistentID)

	// Assert returns nil for not found (GraphQL convention)
	require.NoError(t, err)
	assert.Nil(t, emp, "Non-existent employee should return nil")
}

// TestEmployeesQuery_Success tests querying all employees
func TestEmployeesQuery_Success(t *testing.T) {
	resolver, ctx, dept := setupEmployeeResolverTest(t)

	// Create multiple employees
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

	// Query all employees
	employees, err := resolver.Query().Employees(ctx)

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, employees)
	assert.GreaterOrEqual(t, len(employees), 2)

	// Verify our employees are in the list
	ids := make(map[string]bool)
	for _, emp := range employees {
		ids[emp.ID] = true
	}
	assert.True(t, ids[emp1.ID], "emp1 should be in the list")
	assert.True(t, ids[emp2.ID], "emp2 should be in the list")
}

// TestEmployeesQuery_Empty tests querying when no employees exist
func TestEmployeesQuery_Empty(t *testing.T) {
	// Create fresh test setup without pre-created department
	client := testutil.NewTestEntClient(t)
	defer client.Close()

	deptRepo := database.NewEntDepartmentRepo(client)
	empRepo := database.NewEntEmployeeRepo(client)
	projRepo := database.NewEntProjectRepo(client)
	resolver := NewResolver(deptRepo, empRepo, projRepo)
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "test-request-id")

	// Query all employees (empty database)
	employees, err := resolver.Query().Employees(ctx)

	// Assert success with empty list
	require.NoError(t, err)
	require.NotNil(t, employees)
	assert.Empty(t, employees)
}
