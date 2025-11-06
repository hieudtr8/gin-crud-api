package router

import (
	"gin-crud-api/internal/department"
	"gin-crud-api/internal/employee"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 'Setup' "tiêm" (inject) các handlers vào
func Setup(
	deptHandler *department.Handler,
	empHandler *employee.Handler,
) *gin.Engine {

	r := gin.Default()
	
    // Dùng 'gin.Recovery()' để "bắt" panic và
    // trả về 500 thay vì 'crash' server
	r.Use(gin.Recovery())
    // 'gin.Logger()' để log request
	r.Use(gin.Logger())

	// Gom nhóm routes (Versioning)
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
    
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

	return r
}