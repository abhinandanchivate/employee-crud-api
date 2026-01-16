package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/dto"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/internal/handler"
)

// Mock Service
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) CreateEmployee(req dto.CreateEmployeeRequest) (*models.Employee, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) GetEmployeeByEmail(email string) (*models.Employee, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) ListEmployees(req dto.ListEmployeesRequest) ([]models.Employee, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

func (m *MockEmployeeService) UpdateEmployee(id uint, req dto.UpdateEmployeeRequest) (*models.Employee, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeService) DeleteEmployee(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEmployeeService) SearchEmployees(req dto.SearchEmployeeRequest) ([]models.Employee, int64, error) {
	args := m.Called(req)
	return args.Get(0).([]models.Employee), args.Get(1).(int64), args.Error(2)
}

// Test Suite
type EmployeeHandlerTestSuite struct {
	suite.Suite
	mockService *MockEmployeeService
	handler     *handler.EmployeeHandler
	router      *gin.Engine
}

func TestEmployeeHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeHandlerTestSuite))
}

func (suite *EmployeeHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.mockService = new(MockEmployeeService)
	// This will fail initially - handler doesn't exist
	suite.handler = handler.NewEmployeeHandler(suite.mockService)
	suite.router = gin.New()

	// Setup routes
	suite.setupRoutes()
}

func (suite *EmployeeHandlerTestSuite) setupRoutes() {
	api := suite.router.Group("/api/v1")
	{
		api.POST("/employees", suite.handler.CreateEmployee)
		api.GET("/employees", suite.handler.ListEmployees)
		api.GET("/employees/:id", suite.handler.GetEmployee)
		api.PUT("/employees/:id", suite.handler.UpdateEmployee)
		api.DELETE("/employees/:id", suite.handler.DeleteEmployee)
		api.GET("/employees/search", suite.handler.SearchEmployees)
		api.GET("/employees/email", suite.handler.GetEmployeeByEmail)
	}
}

// RED PHASE: Failing tests for each endpoint

func (suite *EmployeeHandlerTestSuite) TestCreateEmployee_Success() {
	// Test data
	req := dto.CreateEmployeeRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		HireDate:  time.Now(),
	}

	expectedEmployee := &models.Employee{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		HireDate:  req.HireDate,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Setup mock expectation
	suite.mockService.On("CreateEmployee", req).Return(expectedEmployee, nil)

	// Create request
	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	reqHttp, _ := http.NewRequest("POST", "/api/v1/employees", bytes.NewBuffer(jsonData))
	reqHttp.Header.Set("Content-Type", "application/json")

	// Execute - Will fail initially
	suite.router.ServeHTTP(w, reqHttp)

	// Verify
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(suite.T(), response["success"].(bool))
	assert.Equal(suite.T(), "Employee created successfully", response["message"])

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *EmployeeHandlerTestSuite) TestCreateEmployee_InvalidRequest() {
	// Invalid request - missing required fields
	req := map[string]interface{}{
		"first_name": "John",
		// Missing last_name, email, hire_date
	}

	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	reqHttp, _ := http.NewRequest("POST", "/api/v1/employees", bytes.NewBuffer(jsonData))
	reqHttp.Header.Set("Content-Type", "application/json")

	// Execute - Will fail initially
	suite.router.ServeHTTP(w, reqHttp)

	// Verify
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *EmployeeHandlerTestSuite) TestGetEmployee_Success() {
	employeeID := uint(1)
	expectedEmployee := &models.Employee{
		ID:        employeeID,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		IsActive:  true,
		HireDate:  time.Now(),
	}

	// Setup mock expectation
	suite.mockService.On("GetEmployeeByID", employeeID).Return(expectedEmployee, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/employees/1", nil)

	// Execute - Will fail initially
	suite.router.ServeHTTP(w, req)

	// Verify
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(suite.T(), response["success"].(bool))
	assert.Equal(suite.T(), "Employee retrieved successfully", response["message"])

	suite.mockService.AssertExpectations(suite.T())
}

func (suite *EmployeeHandlerTestSuite) TestGetEmployee_NotFound() {
	employeeID := uint(999)

	// Setup mock expectation
	suite.mockService.On("GetEmployeeByID", employeeID).Return(nil, assert.AnError)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/employees/999", nil)

	// Execute - Will fail initially
	suite.router.ServeHTTP(w, req)

	// Verify
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	suite.mockService.AssertExpectations(suite.T())
}

// Add tests for other endpoints...
