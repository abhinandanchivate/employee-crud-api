package service

import (
	"errors"
	"fmt"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/dto"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/internal/repository"
	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
)

type EmployeeService interface {
	CreateEmployee(req dto.CreateEmployeeRequest) (*models.Employee, error)
	GetEmployeeByID(id uint) (*models.Employee, error)
	GetEmployeeByEmail(email string) (*models.Employee, error)
	ListEmployees(req dto.ListEmployeesRequest) ([]models.Employee, int64, error)
	UpdateEmployee(id uint, req dto.UpdateEmployeeRequest) (*models.Employee, error)
	DeleteEmployee(id uint) error
	SearchEmployees(req dto.SearchEmployeeRequest) ([]models.Employee, int64, error)
}

type employeeService struct {
	repo   repository.EmployeeRepository
	logger *logger.Logger
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{
		repo:   repo,
		logger: logger.New(),
	}
}

func (s *employeeService) CreateEmployee(req dto.CreateEmployeeRequest) (*models.Employee, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid create employee request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if email already exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		s.logger.Error("Error checking email existence", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if existing != nil {
		s.logger.Warn("Email already exists", map[string]interface{}{
			"email": req.Email,
		})
		return nil, errors.New("email already exists")
	}

	// Create employee
	employee := &models.Employee{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Position:    req.Position,
		Department:  req.Department,
		Salary:      req.Salary,
		HireDate:    req.HireDate,
		IsActive:    true,
	}

	if err := s.repo.Create(employee); err != nil {
		s.logger.Error("Failed to create employee", map[string]interface{}{
			"error": err.Error(),
			"email": req.Email,
		})
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	s.logger.Info("Employee created successfully", map[string]interface{}{
		"id":    employee.ID,
		"email": employee.Email,
	})

	return employee, nil
}

func (s *employeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	if id == 0 {
		return nil, errors.New("invalid employee ID")
	}

	employee, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.Error("Error fetching employee by ID", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}
	if employee == nil {
		s.logger.Debug("Employee not found", map[string]interface{}{
			"id": id,
		})
		return nil, errors.New("employee not found")
	}

	s.logger.Debug("Employee fetched by ID", map[string]interface{}{
		"id":    employee.ID,
		"email": employee.Email,
	})

	return employee, nil
}

func (s *employeeService) GetEmployeeByEmail(email string) (*models.Employee, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	employee, err := s.repo.FindByEmail(email)
	if err != nil {
		s.logger.Error("Error fetching employee by email", map[string]interface{}{
			"error": err.Error(),
			"email": email,
		})
		return nil, fmt.Errorf("failed to fetch employee: %w", err)
	}
	if employee == nil {
		s.logger.Debug("Employee not found by email", map[string]interface{}{
			"email": email,
		})
		return nil, errors.New("employee not found")
	}

	return employee, nil
}

func (s *employeeService) ListEmployees(req dto.ListEmployeesRequest) ([]models.Employee, int64, error) {
	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid list employees request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	employees, total, err := s.repo.FindAll(req.Page, req.PageSize, req.Active, req.SortBy, req.Order)
	if err != nil {
		s.logger.Error("Error listing employees", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	s.logger.Debug("Employees listed", map[string]interface{}{
		"count": len(employees),
		"page":  req.Page,
		"size":  req.PageSize,
	})

	return employees, total, nil
}

func (s *employeeService) UpdateEmployee(id uint, req dto.UpdateEmployeeRequest) (*models.Employee, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid update employee request", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing employee
	employee, err := s.GetEmployeeByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being updated and if it already exists
	if req.Email != nil && *req.Email != employee.Email {
		existing, err := s.repo.FindByEmail(*req.Email)
		if err != nil {
			s.logger.Error("Error checking email during update", map[string]interface{}{
				"error": err.Error(),
				"email": *req.Email,
			})
			return nil, fmt.Errorf("error checking email: %w", err)
		}
		if existing != nil {
			s.logger.Warn("Email already exists during update", map[string]interface{}{
				"email": *req.Email,
			})
			return nil, errors.New("email already exists")
		}
		employee.Email = *req.Email
	}

	// Update fields if provided
	if req.FirstName != nil {
		employee.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		employee.LastName = *req.LastName
	}
	if req.PhoneNumber != nil {
		employee.PhoneNumber = *req.PhoneNumber
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.Department != nil {
		employee.Department = *req.Department
	}
	if req.Salary != nil {
		employee.Salary = *req.Salary
	}
	if req.HireDate != nil {
		employee.HireDate = *req.HireDate
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}

	if err := s.repo.Update(employee); err != nil {
		s.logger.Error("Error updating employee", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	s.logger.Info("Employee updated successfully", map[string]interface{}{
		"id": employee.ID,
	})

	return employee, nil
}

func (s *employeeService) DeleteEmployee(id uint) error {
	if id == 0 {
		return errors.New("invalid employee ID")
	}

	// Check if employee exists
	_, err := s.GetEmployeeByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Error deleting employee", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	s.logger.Info("Employee deleted successfully", map[string]interface{}{
		"id": id,
	})

	return nil
}

func (s *employeeService) SearchEmployees(req dto.SearchEmployeeRequest) ([]models.Employee, int64, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		s.logger.Error("Invalid search employees request", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, fmt.Errorf("validation failed: %w", err)
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	employees, total, err := s.repo.Search(req.Query, req.Page, req.PageSize)
	if err != nil {
		s.logger.Error("Error searching employees", map[string]interface{}{
			"error": err.Error(),
			"query": req.Query,
		})
		return nil, 0, fmt.Errorf("failed to search employees: %w", err)
	}

	s.logger.Debug("Search completed", map[string]interface{}{
		"query": req.Query,
		"count": len(employees),
		"total": total,
	})

	return employees, total, nil
}
