package models

import "fmt"

// Department là model chính
type Department struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Employee là model chính
type Employee struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	DepartmentID string `json:"department_id"`
}

// --- DTOs (Data Transfer Objects) cho Requests ---
// Dùng DTOs là 'best practice' để tách rời
// API (cách client gửi) và Model (cách lưu trong DB).
// Nó cũng giúp validation dễ dàng.

// CreateDepartmentRequest DTO
type CreateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateDepartmentRequest DTO
type UpdateDepartmentRequest struct {
	Name string `json:"name" binding:"required"`
}

// CreateEmployeeRequest DTO
type CreateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID string `json:"department_id" binding:"required"`
}

// UpdateEmployeeRequest DTO
type UpdateEmployeeRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID string `json:"department_id" binding:"required"`
}

var ErrNotFound = fmt.Errorf("record not found")