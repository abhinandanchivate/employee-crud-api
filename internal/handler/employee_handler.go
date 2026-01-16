package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/dto"
	"github.com/abhinandanchivate/employee-crud-api/internal/service"
	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
	"github.com/abhinandanchivate/employee-crud-api/pkg/response"
)

type EmployeeHandler struct {
	service service.EmployeeService
	logger  *logger.Logger
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
		logger:  logger.New(),
	}
}

// CreateEmployee handles POST /api/v1/employees
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req dto.CreateEmployeeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request body", map[string]interface{}{
			"error": err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	employee, err := h.service.CreateEmployee(req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "email already exists":
			status = http.StatusConflict
		default:
			if err.Error()[:12] == "validation " {
				status = http.StatusBadRequest
			}
		}

		h.logger.Error("Failed to create employee", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		response.Error(c, status, "Failed to create employee", err.Error())
		return
	}

	h.logger.Info("Employee created via API", map[string]interface{}{
		"id":    employee.ID,
		"email": employee.Email,
	})
	response.Success(c, http.StatusCreated, "Employee created successfully", dto.ToEmployeeResponse(employee))
}

// GetEmployee handles GET /api/v1/employees/:id
func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid employee ID in request", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		response.Error(c, http.StatusBadRequest, "Invalid employee ID", err.Error())
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "employee not found" {
			status = http.StatusNotFound
		}

		h.logger.Error("Failed to get employee", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		response.Error(c, status, "Failed to get employee", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Employee retrieved successfully", dto.ToEmployeeResponse(employee))
}

// ListEmployees handles GET /api/v1/employees
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	var req dto.ListEmployeesRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", map[string]interface{}{
			"error": err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	employees, total, err := h.service.ListEmployees(req)
	if err != nil {
		h.logger.Error("Failed to list employees", map[string]interface{}{
			"error": err.Error(),
		})
		response.Error(c, http.StatusInternalServerError, "Failed to list employees", err.Error())
		return
	}

	responseData := dto.ToEmployeesListResponse(employees, total, req.Page, req.PageSize)
	response.Success(c, http.StatusOK, "Employees listed successfully", responseData)
}

// UpdateEmployee handles PUT /api/v1/employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid employee ID in request", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		response.Error(c, http.StatusBadRequest, "Invalid employee ID", err.Error())
		return
	}

	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind request body", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	employee, err := h.service.UpdateEmployee(uint(id), req)
	if err != nil {
		status := http.StatusInternalServerError
		switch err.Error() {
		case "employee not found":
			status = http.StatusNotFound
		case "email already exists":
			status = http.StatusConflict
		default:
			if err.Error()[:12] == "validation " {
				status = http.StatusBadRequest
			}
		}

		h.logger.Error("Failed to update employee", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		response.Error(c, status, "Failed to update employee", err.Error())
		return
	}

	h.logger.Info("Employee updated via API", map[string]interface{}{
		"id": employee.ID,
	})
	response.Success(c, http.StatusOK, "Employee updated successfully", dto.ToEmployeeResponse(employee))
}

// DeleteEmployee handles DELETE /api/v1/employees/:id
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid employee ID in request", map[string]interface{}{
			"error": err.Error(),
			"id":    c.Param("id"),
		})
		response.Error(c, http.StatusBadRequest, "Invalid employee ID", err.Error())
		return
	}

	err = h.service.DeleteEmployee(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "employee not found" {
			status = http.StatusNotFound
		}

		h.logger.Error("Failed to delete employee", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		response.Error(c, status, "Failed to delete employee", err.Error())
		return
	}

	h.logger.Info("Employee deleted via API", map[string]interface{}{
		"id": id,
	})
	response.Success(c, http.StatusOK, "Employee deleted successfully", nil)
}

// SearchEmployees handles GET /api/v1/employees/search
func (h *EmployeeHandler) SearchEmployees(c *gin.Context) {
	var req dto.SearchEmployeeRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Failed to bind query parameters", map[string]interface{}{
			"error": err.Error(),
		})
		response.Error(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	employees, total, err := h.service.SearchEmployees(req)
	if err != nil {
		h.logger.Error("Failed to search employees", map[string]interface{}{
			"error": err.Error(),
			"query": req.Query,
		})
		response.Error(c, http.StatusInternalServerError, "Failed to search employees", err.Error())
		return
	}

	responseData := dto.ToEmployeesListResponse(employees, total, req.Page, req.PageSize)
	response.Success(c, http.StatusOK, "Search completed successfully", responseData)
}

// GetEmployeeByEmail handles GET /api/v1/employees/email
func (h *EmployeeHandler) GetEmployeeByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.Error(c, http.StatusBadRequest, "Email is required", "Query parameter 'email' is required")
		return
	}

	employee, err := h.service.GetEmployeeByEmail(email)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "employee not found" {
			status = http.StatusNotFound
		}

		h.logger.Error("Failed to get employee by email", map[string]interface{}{
			"error": err.Error(),
			"email": email,
		})
		response.Error(c, status, "Failed to get employee", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Employee retrieved successfully", dto.ToEmployeeResponse(employee))
}
