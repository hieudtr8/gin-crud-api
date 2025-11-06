package database

import (
	"context"
	"fmt"

	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/models"

	"github.com/google/uuid"
)

// EntDepartmentRepo implements DepartmentRepository using EntGo
type EntDepartmentRepo struct {
	client *ent.Client
}

// NewEntDepartmentRepo creates a new department repository using EntGo
func NewEntDepartmentRepo(client *ent.Client) DepartmentRepository {
	return &EntDepartmentRepo{client: client}
}

// Save creates a new department in the database
func (r *EntDepartmentRepo) Save(dept *models.Department) error {
	ctx := context.Background()

	// Parse the UUID string to UUID type
	id, err := uuid.Parse(dept.ID)
	if err != nil {
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Create department using EntGo's type-safe builder
	_, err = r.client.Department.
		Create().
		SetID(id).
		SetName(dept.Name).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to save department: %w", err)
	}

	return nil
}

// FindByID retrieves a department by its ID
func (r *EntDepartmentRepo) FindByID(id string) (*models.Department, error) {
	ctx := context.Background()

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Query department using EntGo
	entDept, err := r.client.Department.Get(ctx, uid)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find department: %w", err)
	}

	// Convert EntGo entity to domain model
	return &models.Department{
		ID:   entDept.ID.String(),
		Name: entDept.Name,
	}, nil
}

// FindAll retrieves all departments from the database
func (r *EntDepartmentRepo) FindAll() ([]*models.Department, error) {
	ctx := context.Background()

	// Query all departments
	entDepts, err := r.client.Department.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all departments: %w", err)
	}

	// Convert EntGo entities to domain models
	departments := make([]*models.Department, len(entDepts))
	for i, entDept := range entDepts {
		departments[i] = &models.Department{
			ID:   entDept.ID.String(),
			Name: entDept.Name,
		}
	}

	return departments, nil
}

// Update updates an existing department
func (r *EntDepartmentRepo) Update(dept *models.Department) error {
	ctx := context.Background()

	// Parse the UUID string
	id, err := uuid.Parse(dept.ID)
	if err != nil {
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Update department using EntGo's type-safe builder
	err = r.client.Department.
		UpdateOneID(id).
		SetName(dept.Name).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return models.ErrNotFound
		}
		return fmt.Errorf("failed to update department: %w", err)
	}

	return nil
}

// Delete removes a department from the database
func (r *EntDepartmentRepo) Delete(id string) error {
	ctx := context.Background()

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Delete department using EntGo
	err = r.client.Department.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return models.ErrNotFound
		}
		return fmt.Errorf("failed to delete department: %w", err)
	}

	return nil
}
