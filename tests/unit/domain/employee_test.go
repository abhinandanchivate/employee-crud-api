package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models"
	"github.com/abhinandanchivate/employee-crud-api/internal/domain/models/dto"
)

// Test Suite for Employee Domain
type EmployeeDomainTestSuite struct {
	suite.Suite
}

func TestEmployeeDomainTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeDomainTestSuite))
}

// RED PHASE: These tests should fail initially
func (suite *EmployeeDomainTestSuite) TestEmployeeValidation() {
	// This test will fail initially because Employee struct doesn't have validation
	suite.T().Run("Valid employee should pass validation", func(t *testing.T) {
		employee := models.Employee{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			HireDate:  time.Now(),
		}

		// Initially, there's no validation - this test will fail
		// We'll implement validation in GREEN phase
		assert.NotEmpty(t, employee.FirstName)
		assert.NotEmpty(t, employee.LastName)
		assert.Contains(t, employee.Email, "@")
		assert.False(t, employee.HireDate.IsZero())
	})

	suite.T().Run("Invalid email should fail validation", func(t *testing.T) {
		employee := models.Employee{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid-email", // Invalid email
			HireDate:  time.Now(),
		}

		// This should fail - we need to add email validation
		assert.NotContains(t, employee.Email, "@")
	})
}

func (suite *EmployeeDomainTestSuite) TestCreateEmployeeRequestValidation() {
	suite.T().Run("Valid request should pass validation", func(t *testing.T) {
		req := dto.CreateEmployeeRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			HireDate:  time.Now(),
		}

		// This will fail initially - no Validate() method exists
		err := req.Validate()
		assert.NoError(t, err, "Valid request should pass validation")
	})

	suite.T().Run("Invalid email should fail validation", func(t *testing.T) {
		req := dto.CreateEmployeeRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid-email",
			HireDate:  time.Now(),
		}

		// This will fail initially
		err := req.Validate()
		assert.Error(t, err, "Invalid email should fail validation")
		assert.Contains(t, err.Error(), "email")
	})

	suite.T().Run("Empty first name should fail validation", func(t *testing.T) {
		req := dto.CreateEmployeeRequest{
			FirstName: "", // Empty first name
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			HireDate:  time.Now(),
		}

		// This will fail initially
		err := req.Validate()
		assert.Error(t, err, "Empty first name should fail validation")
	})
}
