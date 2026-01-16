package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/internal/repository"
)

// Mock repository for testing
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) FindByID(id uint) (*models.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindByEmail(email string) (*models.Employee, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) FindAll(page, pageSize int, active *bool, sortBy, order string) ([]models.Employee, int64, error) {
	args := m.Called(page, pageSize, active, sortBy, order)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

func (m *MockEmployeeRepository) Update(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Search(query string, page, pageSize int) ([]models.Employee, int64, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

// RED PHASE: Define what the repository should do
func TestEmployeeRepositoryInterface(t *testing.T) {
	t.Run("Repository should implement all required methods", func(t *testing.T) {
		// This test verifies the interface contract
		var repo repository.EmployeeRepository = &MockEmployeeRepository{}
		assert.NotNil(t, repo, "Repository should implement the interface")

		// Test that we can call all interface methods
		// These calls will fail initially
		assert.Implements(t, (*repository.EmployeeRepository)(nil), repo)
	})
}
