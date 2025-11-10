package router

import (
	"gin-crud-api/internal/legacy/rest/department"
	"gin-crud-api/internal/legacy/rest/employee"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup initializes the router with handlers
func Setup(
	deptHandler *department.Handler,
	empHandler *employee.Handler,
) *gin.Engine {

	r := gin.Default()

	// Middleware
	r.Use(gin.Recovery()) // Recover from panics
	r.Use(gin.Logger())    // Log requests

	// API versioning
	v1 := r.Group("/api/v1")
	{
		// Department Routes
		deptGroup := v1.Group("/departments")
		{
			deptGroup.POST("", deptHandler.Create)
			deptGroup.GET("/:id", deptHandler.Get)
			deptGroup.GET("", deptHandler.List)
			deptGroup.PUT("/:id", deptHandler.Update)
			deptGroup.DELETE("/:id", deptHandler.Delete)
		}

		// Employee Routes
		empGroup := v1.Group("/employees")
		{
			empGroup.POST("", empHandler.Create)
			empGroup.GET("/:id", empHandler.Get)
			empGroup.GET("", empHandler.List)
			empGroup.PUT("/:id", empHandler.Update)
			empGroup.DELETE("/:id", empHandler.Delete)
		}
	}
    
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return r
}