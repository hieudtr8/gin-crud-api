package legacy

import (
	"context"
	"errors"
	"fmt"
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/graph/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresDepartmentRepo implements DepartmentRepository for PostgreSQL
type PostgresDepartmentRepo struct {
	pool *pgxpool.Pool
}

// PostgresEmployeeRepo implements EmployeeRepository for PostgreSQL
type PostgresEmployeeRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresDepartmentRepository(db *PostgresDB) database.DepartmentRepository {
	return &PostgresDepartmentRepo{pool: db.Pool}
}

func NewPostgresEmployeeRepository(db *PostgresDB) database.EmployeeRepository {
	return &PostgresEmployeeRepo{pool: db.Pool}
}

// Department Repository Implementation

func (r *PostgresDepartmentRepo) Save(dept *model.Department) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO departments (id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`
	_, err := r.pool.Exec(ctx, query, dept.ID, dept.Name)
	if err != nil {
		return fmt.Errorf("failed to save department: %w", err)
	}
	return nil
}

func (r *PostgresDepartmentRepo) FindByID(id string) (*model.Department, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name FROM departments WHERE id = $1`

	var dept model.Department
	err := r.pool.QueryRow(ctx, query, id).Scan(&dept.ID, &dept.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find department: %w", err)
	}
	return &dept, nil
}

func (r *PostgresDepartmentRepo) FindAll() ([]*model.Department, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name FROM departments ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %w", err)
	}
	defer rows.Close()

	var departments []*model.Department
	for rows.Next() {
		var dept model.Department
		if err := rows.Scan(&dept.ID, &dept.Name); err != nil {
			return nil, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, &dept)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating departments: %w", err)
	}

	// Return empty slice instead of nil to match in-memory behavior
	if departments == nil {
		departments = make([]*model.Department, 0)
	}

	return departments, nil
}

func (r *PostgresDepartmentRepo) Update(dept *model.Department) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE departments
		SET name = $2, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, dept.ID, dept.Name)
	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *PostgresDepartmentRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Note: CASCADE DELETE is handled by database foreign key constraint
	query := `DELETE FROM departments WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

// Employee Repository Implementation

func (r *PostgresEmployeeRepo) Save(emp *model.Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO employees (id, name, email, department_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`
	_, err := r.pool.Exec(ctx, query, emp.ID, emp.Name, emp.Email, emp.DepartmentID)
	if err != nil {
		return fmt.Errorf("failed to save employee: %w", err)
	}
	return nil
}

func (r *PostgresEmployeeRepo) FindByID(id string) (*model.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, department_id
		FROM employees
		WHERE id = $1
	`

	var emp model.Employee
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&emp.ID, &emp.Name, &emp.Email, &emp.DepartmentID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find employee: %w", err)
	}
	return &emp, nil
}

func (r *PostgresEmployeeRepo) FindAll() ([]*model.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, department_id
		FROM employees
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %w", err)
	}
	defer rows.Close()

	var employees []*model.Employee
	for rows.Next() {
		var emp model.Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.DepartmentID); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, &emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %w", err)
	}

	// Return empty slice instead of nil to match in-memory behavior
	if employees == nil {
		employees = make([]*model.Employee, 0)
	}

	return employees, nil
}

func (r *PostgresEmployeeRepo) Update(emp *model.Employee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE employees
		SET name = $2, email = $3, department_id = $4, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, emp.ID, emp.Name, emp.Email, emp.DepartmentID)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *PostgresEmployeeRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	if result.RowsAffected() == 0 {
		return database.ErrNotFound
	}

	return nil
}

func (r *PostgresEmployeeRepo) FindByDepartmentID(deptID string) ([]*model.Employee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, name, email, department_id
		FROM employees
		WHERE department_id = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch employees by department: %w", err)
	}
	defer rows.Close()

	var employees []*model.Employee
	for rows.Next() {
		var emp model.Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Email, &emp.DepartmentID); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, &emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %w", err)
	}

	// Return empty slice instead of nil
	if employees == nil {
		employees = make([]*model.Employee, 0)
	}

	return employees, nil
}