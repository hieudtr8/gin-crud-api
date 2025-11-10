package database

import (
	"context"
	"fmt"

	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/ent/employee"
	"gin-crud-api/internal/logger"
	"gin-crud-api/internal/models"

	"github.com/google/uuid"
)

// EntEmployeeRepo implements EmployeeRepository using EntGo
type EntEmployeeRepo struct {
	client *ent.Client
}

// NewEntEmployeeRepo creates a new employee repository using EntGo
func NewEntEmployeeRepo(client *ent.Client) EmployeeRepository {
	return &EntEmployeeRepo{client: client}
}

// Save creates a new employee in the database
func (r *EntEmployeeRepo) Save(emp *models.Employee) error {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().
		Str("employee_id", emp.ID).
		Str("name", emp.Name).
		Str("email", emp.Email).
		Str("department_id", emp.DepartmentID).
		Msg("Saving employee to database")

	// Parse UUID strings
	empID, err := uuid.Parse(emp.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", emp.ID).
			Msg("Invalid employee ID format")
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	deptID, err := uuid.Parse(emp.DepartmentID)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", emp.DepartmentID).
			Msg("Invalid department ID format")
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Create employee using EntGo's type-safe builder
	_, err = r.client.Employee.
		Create().
		SetID(empID).
		SetName(emp.Name).
		SetEmail(emp.Email).
		SetDepartmentID(deptID).
		Save(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", emp.ID).
			Msg("Failed to save employee to database")
		return fmt.Errorf("failed to save employee: %w", err)
	}

	log.Debug().
		Str("employee_id", emp.ID).
		Str("name", emp.Name).
		Msg("Employee saved successfully")

	return nil
}

// FindByID retrieves an employee by their ID
func (r *EntEmployeeRepo) FindByID(id string) (*models.Employee, error) {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().
		Str("employee_id", id).
		Msg("Finding employee by ID")

	// Parse UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", id).
			Msg("Invalid employee ID format")
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	// Query employee using EntGo
	entEmp, err := r.client.Employee.Get(ctx, uid)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("employee_id", id).
				Msg("Employee not found in database")
			return nil, models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("employee_id", id).
			Msg("Database error while finding employee")
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

	log.Debug().
		Str("employee_id", entEmp.ID.String()).
		Str("name", entEmp.Name).
		Msg("Employee found successfully")

	// Convert EntGo entity to domain model
	return &models.Employee{
		ID:           entEmp.ID.String(),
		Name:         entEmp.Name,
		Email:        entEmp.Email,
		DepartmentID: entEmp.DepartmentID.String(),
	}, nil
}

// FindAll retrieves all employees from the database
func (r *EntEmployeeRepo) FindAll() ([]*models.Employee, error) {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().Msg("Finding all employees")

	// Query all employees
	entEmps, err := r.client.Employee.Query().All(ctx)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Database error while finding all employees")
		return nil, fmt.Errorf("failed to find all employees: %w", err)
	}

	// Convert EntGo entities to domain models
	employees := make([]*models.Employee, len(entEmps))
	for i, entEmp := range entEmps {
		employees[i] = &models.Employee{
			ID:           entEmp.ID.String(),
			Name:         entEmp.Name,
			Email:        entEmp.Email,
			DepartmentID: entEmp.DepartmentID.String(),
		}
	}

	log.Debug().
		Int("count", len(employees)).
		Msg("All employees found successfully")

	return employees, nil
}

// Update updates an existing employee
func (r *EntEmployeeRepo) Update(emp *models.Employee) error {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().
		Str("employee_id", emp.ID).
		Str("name", emp.Name).
		Str("email", emp.Email).
		Str("department_id", emp.DepartmentID).
		Msg("Updating employee")

	// Parse UUID strings
	empID, err := uuid.Parse(emp.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", emp.ID).
			Msg("Invalid employee ID format")
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	deptID, err := uuid.Parse(emp.DepartmentID)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", emp.DepartmentID).
			Msg("Invalid department ID format")
		return fmt.Errorf("invalid department ID: %w", err)
	}

	// Update employee using EntGo's type-safe builder
	err = r.client.Employee.
		UpdateOneID(empID).
		SetName(emp.Name).
		SetEmail(emp.Email).
		SetDepartmentID(deptID).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("employee_id", emp.ID).
				Msg("Employee not found for update")
			return models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("employee_id", emp.ID).
			Msg("Database error while updating employee")
		return fmt.Errorf("failed to update employee: %w", err)
	}

	log.Debug().
		Str("employee_id", emp.ID).
		Str("name", emp.Name).
		Msg("Employee updated successfully")

	return nil
}

// Delete removes an employee from the database
func (r *EntEmployeeRepo) Delete(id string) error {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().
		Str("employee_id", id).
		Msg("Deleting employee")

	// Parse UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("employee_id", id).
			Msg("Invalid employee ID format")
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	// Delete employee using EntGo
	err = r.client.Employee.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			log.Debug().
				Str("employee_id", id).
				Msg("Employee not found for deletion")
			return models.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("employee_id", id).
			Msg("Database error while deleting employee")
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	log.Debug().
		Str("employee_id", id).
		Msg("Employee deleted successfully")

	return nil
}

// FindByDepartmentID retrieves all employees in a specific department
func (r *EntEmployeeRepo) FindByDepartmentID(deptID string) ([]*models.Employee, error) {
	ctx := context.Background()
	log := logger.WithComponent("EmployeeRepo")

	log.Debug().
		Str("department_id", deptID).
		Msg("Finding employees by department ID")

	// Parse UUID string
	uid, err := uuid.Parse(deptID)
	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", deptID).
			Msg("Invalid department ID format")
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Query employees by department ID using EntGo's type-safe predicate
	entEmps, err := r.client.Employee.
		Query().
		Where(employee.DepartmentID(uid)).
		All(ctx)

	if err != nil {
		log.Error().
			Err(err).
			Str("department_id", deptID).
			Msg("Database error while finding employees by department")
		return nil, fmt.Errorf("failed to find employees by department: %w", err)
	}

	// Convert EntGo entities to domain models
	employees := make([]*models.Employee, len(entEmps))
	for i, entEmp := range entEmps {
		employees[i] = &models.Employee{
			ID:           entEmp.ID.String(),
			Name:         entEmp.Name,
			Email:        entEmp.Email,
			DepartmentID: entEmp.DepartmentID.String(),
		}
	}

	log.Debug().
		Str("department_id", deptID).
		Int("count", len(employees)).
		Msg("Employees found by department successfully")

	return employees, nil
}
