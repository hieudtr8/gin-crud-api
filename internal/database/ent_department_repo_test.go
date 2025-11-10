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

func TestEntDepartmentRepo_Save(t *testing.T) {
	// Setup: Create in-memory database
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Save a new department
	dept := &model.Department{
		ID:   uuid.New().String(),
		Name: "Engineering",
	}

	err := repo.Save(dept)

	// Assert: No error and department was saved
	require.NoError(t, err)

	// Verify: Department can be retrieved
	saved, err := repo.FindByID(dept.ID)
	require.NoError(t, err)
	assert.Equal(t, dept.ID, saved.ID)
	assert.Equal(t, dept.Name, saved.Name)
}

func TestEntDepartmentRepo_Save_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Save department with invalid UUID
	dept := &model.Department{
		ID:   "invalid-uuid",
		Name: "Engineering",
	}

	err := repo.Save(dept)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
}

func TestEntDepartmentRepo_FindByID(t *testing.T) {
	// Setup: Create test department
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Find department by ID
	found, err := repo.FindByID(dept.ID.String())

	// Assert: Department found with correct data
	require.NoError(t, err)
	assert.Equal(t, dept.ID.String(), found.ID)
	assert.Equal(t, dept.Name, found.Name)
}

func TestEntDepartmentRepo_FindByID_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Find non-existent department
	nonExistentID := uuid.New().String()
	found, err := repo.FindByID(nonExistentID)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, found)
}

func TestEntDepartmentRepo_FindByID_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Find with invalid UUID
	found, err := repo.FindByID("invalid-uuid")

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
	assert.Nil(t, found)
}

func TestEntDepartmentRepo_FindAll(t *testing.T) {
	// Setup: Create multiple departments
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	departments := testutil.SeedMultipleDepartments(t, client, []string{
		"Engineering",
		"Sales",
		"Marketing",
	})

	// Test: Find all departments
	found, err := repo.FindAll()

	// Assert: All departments found
	require.NoError(t, err)
	assert.Len(t, found, 3)

	// Verify: Department names are correct
	names := make(map[string]bool)
	for _, d := range found {
		names[d.Name] = true
	}
	assert.True(t, names["Engineering"])
	assert.True(t, names["Sales"])
	assert.True(t, names["Marketing"])

	// Verify: IDs match
	for i, dept := range departments {
		assert.Contains(t, found, &model.Department{
			ID:   dept.ID.String(),
			Name: dept.Name,
		})
		_ = i // avoid unused variable
	}
}

func TestEntDepartmentRepo_FindAll_Empty(t *testing.T) {
	// Setup: Empty database
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Find all in empty database
	found, err := repo.FindAll()

	// Assert: Empty slice returned
	require.NoError(t, err)
	assert.Empty(t, found)
	assert.Len(t, found, 0)
}

func TestEntDepartmentRepo_Update(t *testing.T) {
	// Setup: Create test department
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Update department name
	updated := &model.Department{
		ID:   dept.ID.String(),
		Name: "Engineering & Technology",
	}
	err := repo.Update(updated)

	// Assert: No error
	require.NoError(t, err)

	// Verify: Name was updated
	found, err := repo.FindByID(dept.ID.String())
	require.NoError(t, err)
	assert.Equal(t, "Engineering & Technology", found.Name)
}

func TestEntDepartmentRepo_Update_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Update non-existent department
	nonExistent := &model.Department{
		ID:   uuid.New().String(),
		Name: "Non-existent",
	}
	err := repo.Update(nonExistent)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEntDepartmentRepo_Update_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Update with invalid UUID
	invalid := &model.Department{
		ID:   "invalid-uuid",
		Name: "Invalid",
	}
	err := repo.Update(invalid)

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
}

func TestEntDepartmentRepo_Delete(t *testing.T) {
	// Setup: Create test department
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	dept := testutil.SeedTestDepartment(t, client, "Engineering")

	// Test: Delete department
	err := repo.Delete(dept.ID.String())

	// Assert: No error
	require.NoError(t, err)

	// Verify: Department no longer exists
	found, err := repo.FindByID(dept.ID.String())
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, found)
}

func TestEntDepartmentRepo_Delete_NotFound(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Delete non-existent department
	nonExistentID := uuid.New().String()
	err := repo.Delete(nonExistentID)

	// Assert: Should return ErrNotFound
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestEntDepartmentRepo_Delete_InvalidID(t *testing.T) {
	// Setup
	client := testutil.NewTestEntClient(t)
	defer client.Close()
	repo := NewEntDepartmentRepo(client)

	// Test: Delete with invalid UUID
	err := repo.Delete("invalid-uuid")

	// Assert: Should return error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid department ID")
}
