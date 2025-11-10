package database

import (
	"testing"

	"gin-crud-api/internal/graph/model"
	"gin-crud-api/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func TestEntEmployeeRepo_Save(t *testing.T) {
	// Setup: Create test department and employee
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	emp := &model.Employee{
		ID:           uuid.New().String(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		DepartmentID: dept.ID.String(),
	}

	// Test: Save employee
	err := repo.Save(emp)

	// Assert: No error
	require.NoError(t, err)

	// Verify: Employee can be retrieved
	saved, err := repo.FindByID(emp.ID)
	require.NoError(t, err)
	assert.Equal(t, emp.ID, saved.ID)
	assert.Equal(t, emp.Name, saved.Name)
	assert.Equal(t, emp.Email, saved.Email)
	assert.Equal(t, emp.DepartmentID, saved.DepartmentID)
}

func TestEntEmployeeRepo_Save_InvalidEmployeeID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	emp := &model.Employee{
		ID:           "invalid-uuid",
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		DepartmentID: dept.ID.String(),
	}

	// Test: Save with invalid employee ID
	err := repo.Save(emp)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid employee ID")
}

func TestEntEmployeeRepo_Save_InvalidDepartmentID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	emp := &model.Employee{
		ID:           uuid.New().String(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		DepartmentID: "invalid-uuid",
	}

	// Test: Save with invalid department ID
	err := repo.Save(emp)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
}

func TestEntEmployeeRepo_FindByID(t *testing.T) {
	// Setup: Create test employee
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")
	emp := testutil.SeedTestEmployee(t, client, "John Doe", "john.doe@example.com", dept.ID)

	// Test: Find employee by ID
	found, err := repo.FindByID(emp.ID.String())

	// Assert: Employee found with correct data
	require.NoError(t, err)
	assert.Equal(t, emp.ID.String(), found.ID)
	assert.Equal(t, emp.Name, found.Name)
	assert.Equal(t, emp.Email, found.Email)
	assert.Equal(t, emp.DepartmentID.String(), found.DepartmentID)
}

func TestEntEmployeeRepo_FindByID_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Find non-existent employee
	nonExistentID := uuid.New().String()
	found, err := repo.FindByID(nonExistentID)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, found)
}

func TestEntEmployeeRepo_FindByID_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Find with invalid UUID
	found, err := repo.FindByID("invalid-uuid")

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid employee ID")
	assert.Nil(t, found)
}

func TestEntEmployeeRepo_FindAll(t *testing.T) {
	// Setup: Create test department and multiple employees
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")
	employees := testutil.SeedMultipleEmployees(t, client, dept.ID, 3)

	// Test: Find all employees
	found, err := repo.FindAll()

	// Assert: All employees found
	require.NoError(t, err)
	assert.Len(t, found, 3)

	// Verify: IDs match
	foundIDs := make(map[string]bool)
	for _, emp := range found {
		foundIDs[emp.ID] = true
	}
	for _, emp := range employees {
		assert.True(t, foundIDs[emp.ID.String()], "Employee %s should be found", emp.ID)
	}
}

func TestEntEmployeeRepo_FindAll_Empty(t *testing.T) {
	// Setup: Empty database
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Find all in empty database
	found, err := repo.FindAll()

	// Assert: Empty slice returned
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Len(t, found, 0)
}

func TestEntEmployeeRepo_Update(t *testing.T) {
	// Setup: Create test employee
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept1 := testutil.SeedTestDepartment(t, client, "Engineering")
	dept2 := testutil.SeedTestDepartment(t, client, "Sales")
	emp := testutil.SeedTestEmployee(t, client, "John Doe", "john.doe@example.com", dept1.ID)

	// Test: Update employee
	updated := &model.Employee{
		ID:           emp.ID.String(),
		Name:         "John Smith",
		Email:        "john.smith@example.com",
		DepartmentID: dept2.ID.String(),
	}
	err := repo.Update(updated)

	// Assert: No error
	require.NoError(t, err)

	// Verify: Employee was updated
	found, err := repo.FindByID(emp.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "John Smith", found.Name)
	assert.Equal(t, "john.smith@example.com", found.Email)
	assert.Equal(t, dept2.ID.String(), found.DepartmentID)
}

func TestEntEmployeeRepo_Update_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Update non-existent employee
	nonExistent := &model.Employee{
		ID:           uuid.New().String(),
		Name:         "Non-existent",
		Email:        "non@example.com",
		DepartmentID: dept.ID.String(),
	}
	err := repo.Update(nonExistent)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEntEmployeeRepo_Update_InvalidEmployeeID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Update with invalid employee UUID
	invalid := &model.Employee{
		ID:           "invalid-uuid",
		Name:         "Invalid",
		Email:        "invalid@example.com",
		DepartmentID: dept.ID.String(),
	}
	err := repo.Update(invalid)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid employee ID")
}

func TestEntEmployeeRepo_Update_InvalidDepartmentID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")
	emp := testutil.SeedTestEmployee(t, client, "John Doe", "john.doe@example.com", dept.ID)

	// Test: Update with invalid department UUID
	invalid := &model.Employee{
		ID:           emp.ID.String(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		DepartmentID: "invalid-uuid",
	}
	err := repo.Update(invalid)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
}

func TestEntEmployeeRepo_Delete(t *testing.T) {
	// Setup: Create test employee
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")
	emp := testutil.SeedTestEmployee(t, client, "John Doe", "john.doe@example.com", dept.ID)

	// Test: Delete employee
	err := repo.Delete(emp.ID.String())

	// Assert: No error
	require.NoError(t, err)

	// Verify: Employee no longer exists
	found, err := repo.FindByID(emp.ID.String())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, found)
}

func TestEntEmployeeRepo_Delete_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Delete non-existent employee
	nonExistentID := uuid.New().String()
	err := repo.Delete(nonExistentID)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEntEmployeeRepo_Delete_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Delete with invalid UUID
	err := repo.Delete("invalid-uuid")

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid employee ID")
}

func TestEntEmployeeRepo_FindByDepartmentID(t *testing.T) {
	// Setup: Create department and multiple employees
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept1 := testutil.SeedTestDepartment(t, client, "Engineering")
	dept2 := testutil.SeedTestDepartment(t, client, "Sales")

	// Create 3 employees in dept1
	emp1 := testutil.SeedTestEmployee(t, client, "John Doe", "john@example.com", dept1.ID)
	emp2 := testutil.SeedTestEmployee(t, client, "Jane Doe", "jane@example.com", dept1.ID)
	emp3 := testutil.SeedTestEmployee(t, client, "Bob Smith", "bob@example.com", dept1.ID)

	// Create 1 employee in dept2
	_ = testutil.SeedTestEmployee(t, client, "Alice Johnson", "alice@example.com", dept2.ID)

	// Test: Find employees in dept1
	found, err := repo.FindByDepartmentID(dept1.ID.String())

	// Assert: 3 employees found
	require.NoError(t, err)
	assert.Len(t, found, 3)

	// Verify: All employees belong to dept1
	for _, emp := range found {
		assert.Equal(t, dept1.ID.String(), emp.DepartmentID)
	}

	// Verify: Correct employees found
	foundIDs := make(map[string]bool)
	for _, emp := range found {
		foundIDs[emp.ID] = true
	}
	assert.True(t, foundIDs[emp1.ID.String()])
	assert.True(t, foundIDs[emp2.ID.String()])
	assert.True(t, foundIDs[emp3.ID.String()])
}

func TestEntEmployeeRepo_FindByDepartmentID_Empty(t *testing.T) {
	// Setup: Create department with no employees
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Find employees in empty department
	found, err := repo.FindByDepartmentID(dept.ID.String())

	// Assert: Empty slice returned
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Len(t, found, 0)
}

func TestEntEmployeeRepo_FindByDepartmentID_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntEmployeeRepo(client)

	// Test: Find with invalid department UUID
	found, err := repo.FindByDepartmentID("invalid-uuid")

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
	assert.Nil(t, found)
}
