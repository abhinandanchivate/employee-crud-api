.PHONY: test unit integration build run clean deps

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

# Variables
APP_NAME=employee-crud-api
TEST_COVERAGE=coverage.out

# TDD Workflow Commands
tdd-unit: ## Run TDD cycle for unit tests (Red-Green-Refactor)
	@echo "${YELLOW}=== TDD: RED PHASE (Write failing tests) ===${NC}"
	@echo "${RED}Expected: Tests should fail${NC}"
	@echo ""
	@echo "${YELLOW}=== Running tests... ===${NC}"
	go test ./tests/unit/... -v || true
	@echo ""
	@echo "${GREEN}=== TDD: GREEN PHASE (Make tests pass) ===${NC}"
	@echo "Implement the minimum code to make tests pass"
	@echo ""
	@echo "${YELLOW}=== TDD: REFACTOR PHASE ===${NC}"
	@echo "Refactor code while keeping tests green"

tdd-integration: ## Run TDD cycle for integration tests
	@echo "${YELLOW}Running integration TDD cycle...${NC}"
	go test ./tests/integration/... -v

# Test commands
test: unit integration ## Run all tests

unit: ## Run unit tests
	@echo "${YELLOW}Running unit tests...${NC}"
	go test ./tests/unit/... -v -cover

integration: ## Run integration tests
	@echo "${YELLOW}Running integration tests...${NC}"
	@docker-compose up -d mysql-test
	@sleep 5
	go test ./tests/integration/... -v -cover
	@docker-compose down

test-coverage: ## Generate test coverage report
	@echo "${YELLOW}Generating coverage report...${NC}"
	go test ./... -coverprofile=${TEST_COVERAGE}
	go tool cover -html=${TEST_COVERAGE} -o coverage.html
	@echo "${GREEN}Coverage report generated: coverage.html${NC}"

# Development commands
build: ## Build the application
	@echo "${YELLOW}Building ${APP_NAME}...${NC}"
	go build -o bin/${APP_NAME} ./cmd/api

run: ## Run the application
	@echo "${YELLOW}Starting ${APP_NAME}...${NC}"
	go run ./cmd/api

clean: ## Clean build artifacts
	@echo "${YELLOW}Cleaning up...${NC}"
	rm -rf bin/ coverage.out coverage.html

deps: ## Install dependencies
	@echo "${YELLOW}Installing dependencies...${NC}"
	go mod download
	go mod tidy

# Database commands
db-migrate: ## Run database migrations
	@echo "${YELLOW}Running migrations...${NC}"
	# Add migration commands here

db-seed: ## Seed database with test data
	@echo "${YELLOW}Seeding database...${NC}"
	# Add seed commands here

# Help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${YELLOW}%-20s${NC} %s\n", $$1, $$2}'

.DEFAULT_GOAL := help