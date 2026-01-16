package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "not found")
}

func InternalServerError(c *gin.Context, message string, err string) {
	Error(c, http.StatusInternalServerError, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, "unauthorized")
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, "forbidden")
}

func Conflict(c *gin.Context, message string, err string) {
	Error(c, http.StatusConflict, message, err)
}
