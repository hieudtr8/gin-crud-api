package department

import (
	"gin-crud-api/internal/database"
	"gin-crud-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo    database.DepartmentRepository
	empRepo database.EmployeeRepository // Required for cascade delete
}

func NewHandler(repo database.DepartmentRepository, empRepo database.EmployeeRepository) *Handler {
	return &Handler{
		repo:    repo,
		empRepo: empRepo,
	}
}

// Create handles POST /departments
func (h *Handler) Create(c *gin.Context) {
	var req models.CreateDepartmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dept := &models.Department{
		ID:   uuid.NewString(),
		Name: req.Name,
	}

	if err := h.repo.Save(dept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department"})
		return
	}

	c.JSON(http.StatusCreated, dept)
}

// Get handles GET /departments/:id
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	dept, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.JSON(http.StatusOK, dept)
}

// List handles GET /departments
func (h *Handler) List(c *gin.Context) {
	departments, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}

	if departments == nil {
		departments = []*models.Department{}
	}

	c.JSON(http.StatusOK, departments)
}

// Update handles PUT /departments/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	existingDept, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingDept.Name = req.Name

	if err := h.repo.Update(existingDept); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update department"})
		return
	}

	c.JSON(http.StatusOK, existingDept)
}

// Delete handles DELETE /departments/:id with cascade delete
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	_, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	// Cascade delete: find and delete all employees in this department
	employees, err := h.empRepo.FindByDepartmentID(id)
	if err != nil && err != models.ErrNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check department employees"})
		return
	}

	for _, emp := range employees {
		if err := h.empRepo.Delete(emp.ID); err != nil {
			// TODO: In production, use database transactions for atomic operations
			continue
		}
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete department"})
		return
	}

	c.Status(http.StatusNoContent)
}