package database

import (
	"gin-crud-api/internal/models"
	"sync"
)

// Chúng ta định nghĩa 'interface' trước.
// Điều này cho phép chúng ta dễ dàng "swap" (hoán đổi)
// In-Memory DB này với một DB thật (ví dụ: Postgres)
// mà không cần thay đổi code của 'handler'.
type DepartmentRepository interface {
	Save(dept *models.Department) error
	FindByID(id string) (*models.Department, error)
	FindAll() ([]*models.Department, error)
	Update(dept *models.Department) error
	Delete(id string) error
}

type EmployeeRepository interface {
	Save(emp *models.Employee) error
	FindByID(id string) (*models.Employee, error)
	FindAll() ([]*models.Employee, error)
	Update(emp *models.Employee) error
	Delete(id string) error
	FindByDepartmentID(deptID string) ([]*models.Employee, error)
}

// --- In-Memory Store (Shared Storage) ---
// Đây là DB giả, dùng map và bảo vệ bằng RWMutex
// để đảm bảo an toàn khi chạy song song (concurrent-safe).
type InMemoryStore struct {
	deptMu      sync.RWMutex
	departments map[string]*models.Department

	empMu      sync.RWMutex
	employees  map[string]*models.Employee
}

// NewInMemoryStore là 'constructor'
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		departments: make(map[string]*models.Department),
		employees:   make(map[string]*models.Employee),
	}
}

// --- Separate Repository Implementations ---
type InMemoryDepartmentRepo struct {
	store *InMemoryStore
}

type InMemoryEmployeeRepo struct {
	store *InMemoryStore
}

// Constructor functions for repositories
func NewDepartmentRepository(store *InMemoryStore) DepartmentRepository {
	return &InMemoryDepartmentRepo{store: store}
}

func NewEmployeeRepository(store *InMemoryStore) EmployeeRepository {
	return &InMemoryEmployeeRepo{store: store}
}

// --- Department Repository Methods ---
func (r *InMemoryDepartmentRepo) Save(dept *models.Department) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	r.store.departments[dept.ID] = dept
	return nil
}

func (r *InMemoryDepartmentRepo) FindByID(id string) (*models.Department, error) {
	r.store.deptMu.RLock()
	defer r.store.deptMu.RUnlock()
	dept, ok := r.store.departments[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return dept, nil
}

func (r *InMemoryDepartmentRepo) FindAll() ([]*models.Department, error) {
	r.store.deptMu.RLock()
	defer r.store.deptMu.RUnlock()

	result := make([]*models.Department, 0, len(r.store.departments))
	for _, dept := range r.store.departments {
		result = append(result, dept)
	}
	return result, nil
}

func (r *InMemoryDepartmentRepo) Update(dept *models.Department) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	if _, exists := r.store.departments[dept.ID]; !exists {
		return models.ErrNotFound
	}
	r.store.departments[dept.ID] = dept
	return nil
}

func (r *InMemoryDepartmentRepo) Delete(id string) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	if _, exists := r.store.departments[id]; !exists {
		return models.ErrNotFound
	}
	delete(r.store.departments, id)
	return nil
}

// --- Employee Repository Methods ---
func (r *InMemoryEmployeeRepo) Save(emp *models.Employee) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	r.store.employees[emp.ID] = emp
	return nil
}

func (r *InMemoryEmployeeRepo) FindByID(id string) (*models.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()
	emp, ok := r.store.employees[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return emp, nil
}

func (r *InMemoryEmployeeRepo) FindAll() ([]*models.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()

	result := make([]*models.Employee, 0, len(r.store.employees))
	for _, emp := range r.store.employees {
		result = append(result, emp)
	}
	return result, nil
}

func (r *InMemoryEmployeeRepo) Update(emp *models.Employee) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	if _, exists := r.store.employees[emp.ID]; !exists {
		return models.ErrNotFound
	}
	r.store.employees[emp.ID] = emp
	return nil
}

func (r *InMemoryEmployeeRepo) Delete(id string) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	if _, exists := r.store.employees[id]; !exists {
		return models.ErrNotFound
	}
	delete(r.store.employees, id)
	return nil
}

func (r *InMemoryEmployeeRepo) FindByDepartmentID(deptID string) ([]*models.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()

	var result []*models.Employee
	for _, emp := range r.store.employees {
		if emp.DepartmentID == deptID {
			result = append(result, emp)
		}
	}
	return result, nil
}