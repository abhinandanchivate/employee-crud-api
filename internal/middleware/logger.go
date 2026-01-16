package middleware

import (
	"time"

	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after processing
		timestamp := time.Now()
		latency := timestamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		log := logger.New()
		fields := map[string]interface{}{
			"status":     statusCode,
			"latency":    latency.String(),
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
		}

		if errorMessage != "" {
			fields["error"] = errorMessage
		}

		switch {
		case statusCode >= 500:
			log.Error("Server error", fields)
		case statusCode >= 400:
			log.Warn("Client error", fields)
		default:
			log.Info("Request completed", fields)
		}
	}
}
