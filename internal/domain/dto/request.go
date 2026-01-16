package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateEmployeeRequest struct {
	FirstName   string    `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string    `json:"last_name" validate:"required,min=2,max=100"`
	Email       string    `json:"email" validate:"required,email,max=150"`
	PhoneNumber string    `json:"phone_number" validate:"omitempty,max=20"`
	Position    string    `json:"position" validate:"omitempty,max=100"`
	Department  string    `json:"department" validate:"omitempty,max=100"`
	Salary      float64   `json:"salary" validate:"omitempty,min=0,max=9999999.99"`
	HireDate    time.Time `json:"hire_date" validate:"required"`
}

type UpdateEmployeeRequest struct {
	FirstName   *string    `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName    *string    `json:"last_name" validate:"omitempty,min=2,max=100"`
	Email       *string    `json:"email" validate:"omitempty,email,max=150"`
	PhoneNumber *string    `json:"phone_number" validate:"omitempty,max=20"`
	Position    *string    `json:"position" validate:"omitempty,max=100"`
	Department  *string    `json:"department" validate:"omitempty,max=100"`
	Salary      *float64   `json:"salary" validate:"omitempty,min=0,max=9999999.99"`
	HireDate    *time.Time `json:"hire_date"`
	IsActive    *bool      `json:"is_active"`
}

type SearchEmployeeRequest struct {
	Query    string `form:"q" validate:"required,min=1"`
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
}

type ListEmployeesRequest struct {
	Page     int    `form:"page" validate:"omitempty,min=1"`
	PageSize int    `form:"page_size" validate:"omitempty,min=1,max=100"`
	Active   *bool  `form:"active"`
	SortBy   string `form:"sort_by" validate:"omitempty,oneof=first_name last_name email hire_date salary created_at"`
	Order    string `form:"order" validate:"omitempty,oneof=asc desc"`
}

func (r *CreateEmployeeRequest) Validate() error {
	return validate.Struct(r)
}

func (r *UpdateEmployeeRequest) Validate() error {
	return validate.Struct(r)
}

func (r *SearchEmployeeRequest) Validate() error {
	return validate.Struct(r)
}

func (r *ListEmployeesRequest) Validate() error {
	return validate.Struct(r)
}
