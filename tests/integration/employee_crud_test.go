package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/abhinandanchivate/employee-crud-api/internal/config"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models/dto"
	"github.com/abhinandanchivate/employee-crud-api/pkg/database"
)

type EmployeeCRUDIntegrationTestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *database.Database
}

func TestEmployeeCRUDIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeCRUDIntegrationTestSuite))
}

func (suite *EmployeeCRUDIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	config.Load()

	// Initialize test database
	db, err := database.Init()
	if err != nil {
		suite.T().Fatalf("Failed to initialize database: %v", err)
	}
	suite.db = db

	// Create tables
	suite.db.AutoMigrate(&models.Employee{})

	// Create test server
	suite.server = httptest.NewServer(suite.setupRouter())
}

func (suite *EmployeeCRUDIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()

	// Cleanup database
	suite.db.Exec("DROP TABLE IF EXISTS employees")
	database.Close()
}

func (suite *EmployeeCRUDIntegrationTestSuite) SetupTest() {
	// Clear data before each test
	suite.db.Exec("DELETE FROM employees")
}

func (suite *EmployeeCRUDIntegrationTestSuite) setupRouter() http.Handler {
	// This will be implemented in GREEN phase
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented yet"))
	})
}

// RED PHASE: Integration tests that will fail initially

func (suite *EmployeeCRUDIntegrationTestSuite) TestCreateAndRetrieveEmployee() {
	// Create employee
	createReq := dto.CreateEmployeeRequest{
		FirstName: "Integration",
		LastName:  "Test",
		Email:     "integration.test@example.com",
		HireDate:  time.Now(),
	}

	jsonData, _ := json.Marshal(createReq)
	resp, err := http.Post(suite.server.URL+"/api/v1/employees", "application/json", bytes.NewBuffer(jsonData))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createResponse)
	resp.Body.Close()

	assert.True(suite.T(), createResponse["success"].(bool))

	// Get employee ID from response
	data := createResponse["data"].(map[string]interface{})
	employeeID := int(data["id"].(float64))

	// Retrieve the created employee
	resp, err = http.Get(fmt.Sprintf("%s/api/v1/employees/%d", suite.server.URL, employeeID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var getResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&getResponse)
	resp.Body.Close()

	assert.True(suite.T(), getResponse["success"].(bool))
	assert.Equal(suite.T(), "Integration", getResponse["data"].(map[string]interface{})["first_name"])
}

func (suite *EmployeeCRUDIntegrationTestSuite) TestFullCRUDWorkflow() {
	// Create
	createReq := dto.CreateEmployeeRequest{
		FirstName: "CRUD",
		LastName:  "Test",
		Email:     "crud.test@example.com",
		HireDate:  time.Now(),
	}

	jsonData, _ := json.Marshal(createReq)
	resp, _ := http.Post(suite.server.URL+"/api/v1/employees", "application/json", bytes.NewBuffer(jsonData))

	var createResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createResponse)
	resp.Body.Close()

	data := createResponse["data"].(map[string]interface{})
	employeeID := int(data["id"].(float64))

	// Update
	firstName := "Updated"
	updateReq := dto.UpdateEmployeeRequest{
		FirstName: &firstName,
	}

	jsonData, _ = json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/employees/%d", suite.server.URL, employeeID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Delete
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/employees/%d", suite.server.URL, employeeID), nil)
	resp, err = client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Verify deletion
	resp, err = http.Get(fmt.Sprintf("%s/api/v1/employees/%d", suite.server.URL, employeeID))
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusNotFound, resp.StatusCode)
}
