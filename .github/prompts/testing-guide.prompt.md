# Testing Guidelines for Retail Pro Backend

This document provides guidelines on how to create unit tests and their corresponding documentation for the Retail Pro Backend project.

## important


you just can Equal,noEqual,ni,Not nil

sometime if need you can using contain

all using testify

## Creating Tests

### General Guidelines

1. Use table-driven tests where appropriate
2. Use the `github.com/stretchr/testify/assert` package for assertions
3. Use the `github.com/stretchr/testify/mock` package for mocking interfaces
4. Name test functions with the prefix `Test` followed by the function name being tested
5. Place test files in the same package directory with a `_test.go` suffix
6. Test both success and failure cases

### Creating a New Test File

```go
package yourpackage_test // Note: use yourpackage_test for black-box testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/golang-starter-template/yourpackage" // Import the package being tested
)

func TestYourFunction(t *testing.T) {
	// Test cases table
	tests := []struct {
		name        string
		input       YourInputType
		expected    YourOutputType
		shouldError bool
	}{
		{
			name:        "success case",
			input:       YourInputType{...},
			expected:    YourOutputType{...},
			shouldError: false,
		},
		{
			name:        "failure case",
			input:       YourInputType{...},
			expected:    YourOutputType{...},
			shouldError: true,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := yourpackage.YourFunction(tt.input)
			
			if tt.shouldError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.expected, result, "Result should match expected value")
			}
		})
	}
}
```

### Using Mocks with testify/mock

The `github.com/stretchr/testify/mock` package provides a powerful and flexible way to mock interfaces in your tests.

```go
package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yantology/golang-starter-template/mock"
	"github.com/yantology/golang-starter-template/repository"
)

// Create a mock implementation of your interface
type MockDB struct {
	mock.Mock
}

// Implement the interface methods
func (m *MockDB) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	// Record the call and get the return values from the mock
	args = append([]interface{}{ctx, query}, args...)
	returnArgs := m.Called(args...)
	
	// Return the mocked values
	return returnArgs.Get(0).([]map[string]interface{}), returnArgs.Error(1)
}

func TestRepository(t *testing.T) {
	// Create a new mock
	mockDB := new(MockDB)
	
	// Set expectations
	expectedResult := []map[string]interface{}{{"id": 1, "name": "test"}}
	mockDB.On("Query", mock.Anything, "SELECT * FROM users WHERE id = ?", 1).
		Return(expectedResult, nil)
	
	// Use the mock in your test
	repo := repository.NewRepository(mockDB)
	result, err := repo.GetUserByID(context.Background(), 1)
	
	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
	
	// Assert that the expectations were met
	mockDB.AssertExpectations(t)
}
```

### Running the Tests

To run all tests:
```bash
go test ./...
```

To run tests in a specific package:
```bash
go test ./pkg/packagename
```

To run with verbose output:
```bash
go test -v ./...
```

To run with coverage:
```bash
go test -cover ./...
```

To generate HTML coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Creating Test Documentation

### Documentation Structure

For each test file, create a corresponding markdown documentation file in the `docs/testing/` directory, following the same package path structure.

Example:
- Test file: `pkg/customerror/index_test.go`
- Doc file(make name file same but markdown): `docs/testing/pkg/customerror/index_test.md`

### Documentation Template

```markdown
# [Package Name] Tests

This document describes the test cases for the `[package]` package in the retail-pro-be application.

## Test Overview

Brief description of what functionality is being tested.

## Test Files

List of test files and brief descriptions.

## Test Suites

### 1. [TestFunctionName]

Description of the test suite.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Case 1 | Description | Input values | Expected output |
| Case 2 | Description | Input values | Expected output |

### 2. [TestFunctionName2]

...additional test suites...

## Running the Tests

```bash
go test -v ./[path/to/package]
```

## Test Coverage

List of what's covered in the tests:

1. Feature A
   - Sub-feature 1
   - Sub-feature 2

2. Feature B
   - Sub-feature 1
   - Sub-feature 2
```

## Example

See the existing test documentation for the customerror package as an example:
`docs/testing/pkg/customerror/README.md`

## Best Practices

1. **Keep Tests Independent**: Each test should be able to run independently of others.
2. **Test Edge Cases**: Test boundary conditions and edge cases.
3. **Keep Tests Fast**: Tests should run quickly.
4. **Document Clearly**: Document the purpose of each test case clearly in the documentation.
5. **Use Table-Driven Tests**: Use table-driven tests to reduce code duplication.
6. **Mock External Dependencies**: Use the mock package to isolate your tests from external dependencies.

## Workflow

1. Identify the functionality to test
2. Create the test file with appropriate test cases
3. Run the tests and fix any issues
4. Create documentation in the docs/testing directory
5. Submit changes with both test code and documentation

## Testing Specific Components

### Database Tests

For database tests, use the `github.com/stretchr/testify/mock` package to mock database interfaces rather than connecting to a real database. This makes your tests faster and more reliable.

```go
// Example of mocking a database connection
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	args = append([]interface{}{query}, args...)
	returnArgs := m.Called(args...)
	return returnArgs.Get(0).(sql.Result), returnArgs.Error(1)
}

func TestDatabaseRepository(t *testing.T) {
	mockDB := new(MockDB)
	
	// Set up mock expectations
	mockResult := sqlmock.NewResult(1, 1)  // id=1, affected=1
	mockDB.On("Exec", "INSERT INTO users (name) VALUES (?)", "test").Return(mockResult, nil)
	
	// Test your repository with the mock DB
	repo := repository.NewUserRepository(mockDB)
	id, err := repo.CreateUser(context.Background(), "test")
	
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockDB.AssertExpectations(t)
}
```

### API Tests

For API tests, use the `httptest` package to simulate HTTP requests.

### Middleware Tests

Test middleware by creating mock contexts and requests.

### Config Tests

Test configuration loading from different sources.