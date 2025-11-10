package legacy

// Legacy REST API Request DTOs
// These are only used by the legacy REST handlers in internal/legacy/rest/
// For the current GraphQL API, use the auto-generated types from internal/graph/model/

// CreateDepartmentRequest DTO for REST API
type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateDepartmentRequest DTO for REST API
type UpdateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateEmployeeRequest DTO for REST API
type CreateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID string `json:"department_id" binding:"required"`
}

// UpdateEmployeeRequest DTO for REST API
type UpdateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID string `json:"department_id" binding:"required"`
}
