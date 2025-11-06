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

// List xử lý GET /employees
func (h *Handler) List(c *gin.Context) {
	employees, err := h.empRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
		return
	}

	// Return empty array instead of null
	if employees == nil {
		employees = []*models.Employee{}
	}

	c.JSON(http.StatusOK, employees)
}

// Update xử lý PUT /employees/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	// Check if employee exists
	existingEmp, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Bind and validate request
	var req models.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify new department exists
	_, err = h.deptRepo.FindByID(req.DepartmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	// Update employee
	existingEmp.Name = req.Name
	existingEmp.Email = req.Email
	existingEmp.DepartmentID = req.DepartmentID

	if err := h.empRepo.Update(existingEmp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, existingEmp)
}

// Delete xử lý DELETE /employees/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Check if employee exists
	_, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Delete employee
	if err := h.empRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.Status(http.StatusNoContent)
}