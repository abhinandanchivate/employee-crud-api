package dto

import (
	"time"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
)

type EmployeeResponse struct {
	ID          uint      `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Position    string    `json:"position,omitempty"`
	Department  string    `json:"department,omitempty"`
	Salary      float64   `json:"salary,omitempty"`
	HireDate    time.Time `json:"hire_date"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeesListResponse struct {
	Employees  []EmployeeResponse `json:"employees"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ToEmployeeResponse(employee *models.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:          employee.ID,
		FirstName:   employee.FirstName,
		LastName:    employee.LastName,
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		Position:    employee.Position,
		Department:  employee.Department,
		Salary:      employee.Salary,
		HireDate:    employee.HireDate,
		IsActive:    employee.IsActive,
		CreatedAt:   employee.CreatedAt,
		UpdatedAt:   employee.UpdatedAt,
	}
}

func ToEmployeesListResponse(employees []models.Employee, total int64, page, pageSize int) EmployeesListResponse {
	employeeResponses := make([]EmployeeResponse, len(employees))
	for i, emp := range employees {
		employeeResponses[i] = ToEmployeeResponse(&emp)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return EmployeesListResponse{
		Employees:  employeeResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
