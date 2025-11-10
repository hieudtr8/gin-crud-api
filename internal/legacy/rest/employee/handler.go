package employee

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/legacy"
	"gin-crud-api/internal/graph/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	empRepo  database.EmployeeRepository
	deptRepo database.DepartmentRepository // Required for department validation
}

func NewHandler(empRepo database.EmployeeRepository, deptRepo database.DepartmentRepository) *Handler {
	return &Handler{empRepo: empRepo, deptRepo: deptRepo}
}

// Create handles POST /employees
func (h *Handler) Create(c *gin.Context) {
	var req legacy.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate department exists
	_, err := h.deptRepo.FindByID(req.DepartmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return
	}

	emp := &model.Employee{
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

// Get handles GET /employees/:id
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	emp, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, emp)
}

// List handles GET /employees
func (h *Handler) List(c *gin.Context) {
	employees, err := h.empRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve employees"})
		return
	}

	if employees == nil {
		employees = []*model.Employee{}
	}

	c.JSON(http.StatusOK, employees)
}

// Update handles PUT /employees/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	existingEmp, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var req legacy.UpdateEmployeeRequest
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

	existingEmp.Name = req.Name
	existingEmp.Email = req.Email
	existingEmp.DepartmentID = req.DepartmentID

	if err := h.empRepo.Update(existingEmp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update employee"})
		return
	}

	c.JSON(http.StatusOK, existingEmp)
}

// Delete handles DELETE /employees/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.empRepo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	if err := h.empRepo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete employee"})
		return
	}

	c.Status(http.StatusNoContent)
}