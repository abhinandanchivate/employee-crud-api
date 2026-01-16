package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/dto"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/internal/service"
)

// Mock repository
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockRepo) FindByID(id uint) (*models.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockRepo) FindByEmail(email string) (*models.Employee, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockRepo) FindAll(page, pageSize int, active *bool, sortBy, order string) ([]models.Employee, int64, error) {
	args := m.Called(page, pageSize, active, sortBy, order)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepo) Update(employee *models.Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockRepo) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) Search(query string, page, pageSize int) ([]models.Employee, int64, error) {
	args := m.Called(query, page, pageSize)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

// Test Suite
type EmployeeServiceTestSuite struct {
	suite.Suite
	mockRepo *MockRepo
	service  service.EmployeeService
}

func TestEmployeeServiceTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeServiceTestSuite))
}

func (suite *EmployeeServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepo)
	// This will fail initially - service.NewEmployeeService doesn't exist
	suite.service = service.NewEmployeeService(suite.mockRepo)
}

// RED PHASE: Create failing tests for each endpoint

func (suite *EmployeeServiceTestSuite) TestCreateEmployee_Success() {
	// Test data
	req := dto.CreateEmployeeRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		HireDate:  time.Now(),
	}

	// Setup expectations
	suite.mockRepo.On("FindByEmail", req.Email).Return(nil, nil)
	suite.mockRepo.On("Create", mock.AnythingOfType("*models.Employee")).Return(nil)

	// Execute - This will fail initially
	employee, err := suite.service.CreateEmployee(req)

	// Verify
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), employee)
	assert.Equal(suite.T(), "John", employee.FirstName)
	assert.Equal(suite.T(), "john.doe@example.com", employee.Email)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeServiceTestSuite) TestCreateEmployee_DuplicateEmail() {
	req := dto.CreateEmployeeRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		HireDate:  time.Now(),
	}

	existingEmployee := &models.Employee{
		ID:    1,
		Email: "john.doe@example.com",
	}

	// Setup expectation
	suite.mockRepo.On("FindByEmail", req.Email).Return(existingEmployee, nil)

	// Execute - Will fail initially
	employee, err := suite.service.CreateEmployee(req)

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), employee)
	assert.Equal(suite.T(), "email already exists", err.Error())

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeServiceTestSuite) TestGetEmployeeByID_Success() {
	employeeID := uint(1)
	expectedEmployee := &models.Employee{
		ID:    employeeID,
		Email: "john.doe@example.com",
	}

	// Setup expectation
	suite.mockRepo.On("FindByID", employeeID).Return(expectedEmployee, nil)

	// Execute - Will fail initially
	employee, err := suite.service.GetEmployeeByID(employeeID)

	// Verify
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), employee)
	assert.Equal(suite.T(), employeeID, employee.ID)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeServiceTestSuite) TestGetEmployeeByID_NotFound() {
	employeeID := uint(999)

	// Setup expectation
	suite.mockRepo.On("FindByID", employeeID).Return(nil, nil)

	// Execute - Will fail initially
	employee, err := suite.service.GetEmployeeByID(employeeID)

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), employee)
	assert.Equal(suite.T(), "employee not found", err.Error())

	suite.mockRepo.AssertExpectations(suite.T())
}

// Add tests for other service methods...
