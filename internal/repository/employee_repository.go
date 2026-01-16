package repository

import (
	"errors"

	"gorm.io/gorm"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/pkg/logger"
)

type EmployeeRepository interface {
	Create(employee *models.Employee) error
	FindByID(id uint) (*models.Employee, error)
	FindByEmail(email string) (*models.Employee, error)
	FindAll(page, pageSize int, active *bool, sortBy, order string) ([]models.Employee, int64, error)
	Update(employee *models.Employee) error
	Delete(id uint) error
	Search(query string, page, pageSize int) ([]models.Employee, int64, error)
}

type employeeRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db:     db,
		logger: logger.New(),
	}
}

func (r *employeeRepository) Create(employee *models.Employee) error {
	if err := r.db.Create(employee).Error; err != nil {
		r.logger.Error("Failed to create employee", map[string]interface{}{
			"error": err.Error(),
			"email": employee.Email,
		})
		return err
	}

	r.logger.Info("Employee created successfully", map[string]interface{}{
		"id":    employee.ID,
		"email": employee.Email,
	})
	return nil
}

func (r *employeeRepository) FindByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Debug("Employee not found", map[string]interface{}{
				"id": id,
			})
			return nil, nil
		}
		r.logger.Error("Failed to find employee by ID", map[string]interface{}{
			"error": err.Error(),
			"id":    id,
		})
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindByEmail(email string) (*models.Employee, error) {
	var employee models.Employee
	err := r.db.Where("email = ?", email).First(&employee).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error("Failed to find employee by email", map[string]interface{}{
			"error": err.Error(),
			"email": email,
		})
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindAll(page, pageSize int, active *bool, sortBy, order string) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.Model(&models.Employee{})

	if active != nil {
		query = query.Where("is_active = ?", *active)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count employees", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Apply sorting
	if sortBy != "" {
		if order == "" {
			order = "asc"
		}
		query = query.Order(sortBy + " " + order)
	} else {
		query = query.Order("created_at DESC")
	}

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		r.logger.Error("Failed to find employees", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, err
	}

	r.logger.Debug("Retrieved employees list", map[string]interface{}{
		"count": len(employees),
		"page":  page,
		"size":  pageSize,
	})

	return employees, total, nil
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	if err := r.db.Save(employee).Error; err != nil {
		r.logger.Error("Failed to update employee", map[string]interface{}{
			"error": err.Error(),
			"id":    employee.ID,
		})
		return err
	}

	r.logger.Info("Employee updated successfully", map[string]interface{}{
		"id": employee.ID,
	})
	return nil
}

func (r *employeeRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Employee{}, id)
	if result.Error != nil {
		r.logger.Error("Failed to delete employee", map[string]interface{}{
			"error": result.Error.Error(),
			"id":    id,
		})
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("No employee found to delete", map[string]interface{}{
			"id": id,
		})
		return gorm.ErrRecordNotFound
	}

	r.logger.Info("Employee deleted successfully", map[string]interface{}{
		"id": id,
	})
	return nil
}

func (r *employeeRepository) Search(query string, page, pageSize int) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := r.db.Model(&models.Employee{}).
		Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR position LIKE ? OR department LIKE ?",
			searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count search results", map[string]interface{}{
			"error": err.Error(),
			"query": query,
		})
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := dbQuery.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&employees).Error; err != nil {
		r.logger.Error("Failed to search employees", map[string]interface{}{
			"error": err.Error(),
			"query": query,
		})
		return nil, 0, err
	}

	r.logger.Debug("Search completed", map[string]interface{}{
		"query": query,
		"count": len(employees),
		"total": total,
	})

	return employees, total, nil
}
