package database

import (
	"context"
	"fmt"

	"gin-crud-api/internal/ent"
	"gin-crud-api/internal/ent/employee"
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

	// Parse UUID strings
	empID, err := uuid.Parse(emp.ID)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	deptID, err := uuid.Parse(emp.DepartmentID)
	if err != nil {
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
		return fmt.Errorf("failed to save employee: %w", err)
	}

	return nil
}

// FindByID retrieves an employee by their ID
func (r *EntEmployeeRepo) FindByID(id string) (*models.Employee, error) {
	ctx := context.Background()

	// Parse UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid employee ID: %w", err)
	}

	// Query employee using EntGo
	entEmp, err := r.client.Employee.Get(ctx, uid)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}

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

	// Query all employees
	entEmps, err := r.client.Employee.Query().All(ctx)
	if err != nil {
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

	return employees, nil
}

// Update updates an existing employee
func (r *EntEmployeeRepo) Update(emp *models.Employee) error {
	ctx := context.Background()

	// Parse UUID strings
	empID, err := uuid.Parse(emp.ID)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	deptID, err := uuid.Parse(emp.DepartmentID)
	if err != nil {
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
			return models.ErrNotFound
		}
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

// Delete removes an employee from the database
func (r *EntEmployeeRepo) Delete(id string) error {
	ctx := context.Background()

	// Parse UUID string
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid employee ID: %w", err)
	}

	// Delete employee using EntGo
	err = r.client.Employee.DeleteOneID(uid).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return models.ErrNotFound
		}
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}

// FindByDepartmentID retrieves all employees in a specific department
func (r *EntEmployeeRepo) FindByDepartmentID(deptID string) ([]*models.Employee, error) {
	ctx := context.Background()

	// Parse UUID string
	uid, err := uuid.Parse(deptID)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID: %w", err)
	}

	// Query employees by department ID using EntGo's type-safe predicate
	entEmps, err := r.client.Employee.
		Query().
		Where(employee.DepartmentID(uid)).
		All(ctx)

	if err != nil {
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

	return employees, nil
}
