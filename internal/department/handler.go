package department

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler 'giữ' (holds) repository interface.
// Đây là 'Dependency Injection'.
type Handler struct {
	repo    database.DepartmentRepository
	empRepo database.EmployeeRepository // Added for cascade delete
}

// NewHandler là 'constructor'
func NewHandler(repo database.DepartmentRepository, empRepo database.EmployeeRepository) *Handler {
	return &Handler{
		repo:    repo,
		empRepo: empRepo,
	}
}

// Create xử lý POST /departments
func (h *Handler) Create(c *gin.Context) {
	var req models.CreateDepartmentRequest
	
    // 'ShouldBindJSON' là cách 'Gin' validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept := &models.Department{
		ID:   uuid.NewString(), // Tạo ID ngẫu nhiên
		Name: req.Name,
	}

	if err := h.repo.Save(dept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, dept)
}

// Get xử lý GET /departments/:id
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	dept, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, dept)
}

// List xử lý GET /departments
func (h *Handler) List(c *gin.Context) {
	departments, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}

	// Return empty array instead of null
	if departments == nil {
		departments = []*models.Department{}
	}

	c.JSON(http.StatusOK, departments)
}

// Update xử lý PUT /departments/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	// Check if department exists
	existingDept, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Bind and validate request
	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update department
	existingDept.Name = req.Name

	if err := h.repo.Update(existingDept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, existingDept)
}

// Delete xử lý DELETE /departments/:id with cascade delete
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Check if department exists
	_, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Find all employees in this department for cascade delete
	employees, err := h.empRepo.FindByDepartmentID(id)
	if err != nil && err != models.ErrNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department employees"})
		return
	}

	// Delete all employees in the department (cascade delete)
	for _, emp := range employees {
		if err := h.empRepo.Delete(emp.ID); err != nil {
			// Log error but continue with other deletions
			// In production, this should be a transaction
			continue
		}
	}

	// Delete the department
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.Status(http.StatusNoContent)
}