package middleware

import (
	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
	"github.com/abhinandanchivate/employee-crud-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := logger.New()
				log.Error("Panic recovered", map[string]interface{}{
					"error": err,
					"path":  c.Request.URL.Path,
				})

				response.InternalServerError(c, "Internal server error", "panic recovered")
				c.Abort()
			}
		}()
		c.Next()
	}
}
