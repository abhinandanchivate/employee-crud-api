package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/abhinandanchivate/employee-crud-api/internal/config"
	"github.com/abhinandanchivate/employee-crud-api/internal/routes"
	"github.com/abhinandanchivate/employee-crud-api/pkg/database"
	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
)

func main() {
	// Load configuration
	config.Load()
	cfg := config.Get()

	// Initialize logger
	logger := logger.New()
	logger.Info("Starting application", map[string]interface{}{
		"name":    cfg.App.Name,
		"version": cfg.App.Version,
		"env":     cfg.App.Env,
	})

	// Initialize database
	db, err := database.Init()
	if err != nil {
		logger.Fatal("Failed to initialize database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer database.Close()

	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, db.DB)

	// Start server
	addr := ":" + cfg.App.Port
	logger.Info("Server starting", map[string]interface{}{
		"port": cfg.App.Port,
		"addr": addr,
	})

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
