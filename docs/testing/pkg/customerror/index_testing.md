# Customer Error Package Tests

This document describes the test cases for the `customerror` package in the retail-pro-be application.

## Test Overview

The tests for the `customerror` package verify the error handling and custom error creation functionality, with a focus on PostgreSQL-specific error handling. The tests use table-driven tests with the `github.com/stretchr/testify/assert` package.

## Test Files

- `pkg/customerror/index_test.go`: Contains all the tests for the customerror package functions and methods

## Test Suites

### 1. TestNewCustomError

This test suite validates the functionality of the `NewCustomError` function and the associated methods (`Message()`, `Code()`, and `Original()`).

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Basic Error | Tests creation and methods with a normal error | `original`: errors.New("original error")<br>`message`: "test error"<br>`httpCode`: 400 | `Message()`: "test error"<br>`Code()`: 400<br>`Original()`: "original error" |
| Nil Original Error | Tests creation and methods with nil error | `original`: nil<br>`message`: "test error with nil original"<br>`httpCode`: 500 | `Message()`: "test error with nil original"<br>`Code()`: 500<br>`Original()`: "" |

### 2. TestNewPostgresError

This test suite validates the `NewPostgresError` function which provides specialized handling for PostgreSQL errors.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Nil Error | Tests handling of nil error | `err`: nil | Returns nil |
| Unique Violation | Tests PostgreSQL unique violation | `err`: pq.Error{Code: "23505"} | `Message()`: "Record already exists"<br>`Code()`: 409 |
| Foreign Key Violation | Tests foreign key constraint violation | `err`: pq.Error{Code: "23503"} | `Message()`: "Foreign key violation"<br>`Code()`: 400 |
| String Data Too Long | Tests string truncation error | `err`: pq.Error{Code: "22001"} | `Message()`: "String data is too long"<br>`Code()`: 400 |
| Unhandled Postgres Error | Tests unhandled PostgreSQL error | `err`: pq.Error{Code: "42P01"} | `Message()`: "Database error"<br>`Code()`: 500 |
| Record Not Found | Tests SQL no rows error | `err`: sql.ErrNoRows | `Message()`: "Record not found"<br>`Code()`: 404 |
| Generic Error | Tests non-PostgreSQL errors | `err`: generic error | `Message()`: "Database error"<br>`Code()`: 500 |

## Running the Tests

To run these tests, use the following command from the project root directory:

```bash
go test -v ./pkg/customerror
```

For coverage report:
```bash
go test -coverprofile=coverage.out ./pkg/customerror
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage

The tests cover:

1. CustomError Type
   - Creation with various error types
   - Message retrieval
   - HTTP status code retrieval
   - Original error message retrieval
   - Handling of nil errors

2. PostgreSQL Error Handling
   - Common PostgreSQL error codes (23505, 23503, 22001)
   - SQL-specific errors (ErrNoRows)
   - Unhandled PostgreSQL errors
   - Generic Go errors
   - Nil error handling

3. Edge Cases
   - Nil original errors
   - Empty messages
   - Various HTTP status codes