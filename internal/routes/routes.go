package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/abhinandanchivate/employee-crud-api/internal/handler"
	"github.com/abhinandanchivate/employee-crud-api/internal/middleware"
	"github.com/abhinandanchivate/employee-crud-api/internal/repository"
	"github.com/abhinandanchivate/employee-crud-api/internal/service"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "employee-crud-api",
		})
	})

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		setupEmployeeRoutes(apiV1, db)
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "route not found",
		})
	})
}

func setupEmployeeRoutes(router *gin.RouterGroup, db *gorm.DB) {
	// Initialize dependencies
	employeeRepo := repository.NewEmployeeRepository(db)
	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Employee routes
	employees := router.Group("/employees")
	{
		employees.POST("", employeeHandler.CreateEmployee)
		employees.GET("", employeeHandler.ListEmployees)
		employees.GET("/search", employeeHandler.SearchEmployees)
		employees.GET("/email", employeeHandler.GetEmployeeByEmail)

		employees.GET("/:id", employeeHandler.GetEmployee)
		employees.PUT("/:id", employeeHandler.UpdateEmployee)
		employees.DELETE("/:id", employeeHandler.DeleteEmployee)
	}
}
