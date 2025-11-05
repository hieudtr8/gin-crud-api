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
	repo database.DepartmentRepository
	// Chúng ta có thể 'inject' EmployeeRepo nếu cần
}

// NewHandler là 'constructor'
func NewHandler(repo database.DepartmentRepository) *Handler {
	return &Handler{repo: repo}
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

// (Tương tự, ông implement List, Update, Delete...)
//
// func (h *Handler) List(c *gin.Context) { ... }
// func (h *Handler) Update(c *gin.Context) { ... }
// func (h *Handler) Delete(c *gin.Context) { ... }