package legacy

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/graph/model"
	"sync"
)

// InMemoryStore provides thread-safe in-memory storage using RWMutex
type InMemoryStore struct {
	deptMu      sync.RWMutex
	departments map[string]*model.Department

	empMu      sync.RWMutex
	employees  map[string]*model.Employee
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		departments: make(map[string]*model.Department),
		employees:   make(map[string]*model.Employee),
	}
}

type InMemoryDepartmentRepo struct {
	store *InMemoryStore
}

type InMemoryEmployeeRepo struct {
	store *InMemoryStore
}

func NewDepartmentRepository(store *InMemoryStore) database.DepartmentRepository {
	return &InMemoryDepartmentRepo{store: store}
}

func NewEmployeeRepository(store *InMemoryStore) database.EmployeeRepository {
	return &InMemoryEmployeeRepo{store: store}
}

// Department Repository Implementation

func (r *InMemoryDepartmentRepo) Save(dept *model.Department) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	r.store.departments[dept.ID] = dept
	return nil
}

func (r *InMemoryDepartmentRepo) FindByID(id string) (*model.Department, error) {
	r.store.deptMu.RLock()
	defer r.store.deptMu.RUnlock()
	dept, ok := r.store.departments[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	return dept, nil
}

func (r *InMemoryDepartmentRepo) FindAll() ([]*model.Department, error) {
	r.store.deptMu.RLock()
	defer r.store.deptMu.RUnlock()

	result := make([]*model.Department, 0, len(r.store.departments))
	for _, dept := range r.store.departments {
		result = append(result, dept)
	}
	return result, nil
}

func (r *InMemoryDepartmentRepo) Update(dept *model.Department) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	if _, exists := r.store.departments[dept.ID]; !exists {
		return database.ErrNotFound
	}
	r.store.departments[dept.ID] = dept
	return nil
}

func (r *InMemoryDepartmentRepo) Delete(id string) error {
	r.store.deptMu.Lock()
	defer r.store.deptMu.Unlock()
	if _, exists := r.store.departments[id]; !exists {
		return database.ErrNotFound
	}
	delete(r.store.departments, id)
	return nil
}

// Employee Repository Implementation

func (r *InMemoryEmployeeRepo) Save(emp *model.Employee) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	r.store.employees[emp.ID] = emp
	return nil
}

func (r *InMemoryEmployeeRepo) FindByID(id string) (*model.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()
	emp, ok := r.store.employees[id]
	if !ok {
		return nil, database.ErrNotFound
	}
	return emp, nil
}

func (r *InMemoryEmployeeRepo) FindAll() ([]*model.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()

	result := make([]*model.Employee, 0, len(r.store.employees))
	for _, emp := range r.store.employees {
		result = append(result, emp)
	}
	return result, nil
}

func (r *InMemoryEmployeeRepo) Update(emp *model.Employee) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	if _, exists := r.store.employees[emp.ID]; !exists {
		return database.ErrNotFound
	}
	r.store.employees[emp.ID] = emp
	return nil
}

func (r *InMemoryEmployeeRepo) Delete(id string) error {
	r.store.empMu.Lock()
	defer r.store.empMu.Unlock()
	if _, exists := r.store.employees[id]; !exists {
		return database.ErrNotFound
	}
	delete(r.store.employees, id)
	return nil
}

func (r *InMemoryEmployeeRepo) FindByDepartmentID(deptID string) ([]*model.Employee, error) {
	r.store.empMu.RLock()
	defer r.store.empMu.RUnlock()

	var result []*model.Employee
	for _, emp := range r.store.employees {
		if emp.DepartmentID == deptID {
			result = append(result, emp)
		}
	}
	return result, nil
}