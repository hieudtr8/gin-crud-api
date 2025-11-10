package database

import (
	"context"
	"fmt"

	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/logger"
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
	log := logger.WithComponent("DepartmentRepo")

	log.Debug().
		Str("department_id", dept.ID).
		Str("name", dept.Name).
		Msg("Saving department to database")

	// Parse the UUID string to UUID type
	id, err := uuid.Parse(dept.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", dept.ID).
			Msg("Invalid department ID format")
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Create department using EntGo's type-safe builder
	_, err = r.client.Department.
		Create().
		SetID(id).
		SetName(dept.Name).
		Save(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", dept.ID).
			Msg("Failed to save department to database")
		return fmt.Errorf("failed to save department: %w", err)
	}

	log.Debug().
		Str("department_id", dept.ID).
		Str("name", dept.Name).
		Msg("Department saved successfully")

	return nil
}

// FindByID retrieves a department by its ID
func (r *EntDepartmentRepo) FindByID(id string) (*models.Department, error) {
	ctx := context.Background()
	log := logger.WithComponent("DepartmentRepo")

	log.Debug().
		Str("department_id", id).
		Msg("Finding department by ID")

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", id).
			Msg("Invalid department ID format")
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Query department using EntGo
	entDept, err := r.client.Department.Get(ctx, uid)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("department_id", id).
				Msg("Department not found in database")
			return nil, models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("department_id", id).
			Msg("Database error while finding department")
		return nil, fmt.Errorf("failed to find department: %w", err)
	}

	log.Debug().
		Str("department_id", entDept.ID.String()).
		Str("name", entDept.Name).
		Msg("Department found successfully")

	// Convert EntGo entity to domain model
	return &models.Department{
		ID:   entDept.ID.String(),
		Name: entDept.Name,
	}, nil
}

// FindAll retrieves all departments from the database
func (r *EntDepartmentRepo) FindAll() ([]*models.Department, error) {
	ctx := context.Background()
	log := logger.WithComponent("DepartmentRepo")

	log.Debug().Msg("Finding all departments")

	// Query all departments
	entDepts, err := r.client.Department.Query().All(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Database error while finding all departments")
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

	log.Debug().
		Int("count", len(departments)).
		Msg("All departments found successfully")

	return departments, nil
}

// Update updates an existing department
func (r *EntDepartmentRepo) Update(dept *models.Department) error {
	ctx := context.Background()
	log := logger.WithComponent("DepartmentRepo")

	log.Debug().
		Str("department_id", dept.ID).
		Str("name", dept.Name).
		Msg("Updating department")

	// Parse the UUID string
	id, err := uuid.Parse(dept.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", dept.ID).
			Msg("Invalid department ID format")
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Update department using EntGo's type-safe builder
	err = r.client.Department.
		UpdateOneID(id).
		SetName(dept.Name).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("department_id", dept.ID).
				Msg("Department not found for update")
			return models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("department_id", dept.ID).
			Msg("Database error while updating department")
		return fmt.Errorf("failed to update department: %w", err)
	}

	log.Debug().
		Str("department_id", dept.ID).
		Str("name", dept.Name).
		Msg("Department updated successfully")

	return nil
}

// Delete removes a department from the database
func (r *EntDepartmentRepo) Delete(id string) error {
	ctx := context.Background()
	log := logger.WithComponent("DepartmentRepo")

	log.Debug().
		Str("department_id", id).
		Msg("Deleting department")

	// Parse the UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", id).
			Msg("Invalid department ID format")
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Delete department using EntGo
	err = r.client.Department.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("department_id", id).
				Msg("Department not found for deletion")
			return models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("department_id", id).
			Msg("Database error while deleting department")
		return fmt.Errorf("failed to delete department: %w", err)
	}

	log.Debug().
		Str("department_id", id).
		Msg("Department deleted successfully")

	return nil
}
