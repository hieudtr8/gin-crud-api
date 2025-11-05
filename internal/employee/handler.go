package employee

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	// Employee handler cần cả 2 repo
	empRepo  database.EmployeeRepository
	deptRepo database.DepartmentRepository
}

func NewHandler(empRepo database.EmployeeRepository, deptRepo database.DepartmentRepository) *Handler {
	return &Handler{empRepo: empRepo, deptRepo: deptRepo}
}

// Create xử lý POST /employees
func (h *Handler) Create(c *gin.Context) {
	var req models.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logic "business" (kiểm tra xem Department có tồn tại không)
	// Đây là lý do chúng ta cần 'deptRepo'
	_, err := h.deptRepo.FindByID(req.DepartmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	emp := &models.Employee{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		DepartmentID: req.DepartmentID,
	}

	if err := h.empRepo.Save(emp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, emp)
}

// Get xử lý GET /employees/:id
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// (Tương tự, ông implement List, Update, Delete...)