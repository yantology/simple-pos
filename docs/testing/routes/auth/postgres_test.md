# Auth Postgres Tests

This document describes the test cases for the `auth` package's postgres implementation in the retail-pro-be application.

## Test Overview

These tests cover the PostgreSQL database operations for user authentication and management. The tests use `go-sqlmock` to mock database interactions and `testify/assert` for assertions.

## Test Files

- `postgres_test.go`: Contains all database operation tests for the auth package

## Test Suites

### 1. TestCheckExistingEmail

Tests the email existence checking functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| email exists | Tests when email exists in database | email: "test@example.com" | No error |
| email not exists | Tests when email does not exist | email: "notfound@example.com" | Error: sql.ErrNoRows |

### 2. TestSaveActivationToken

Tests saving activation tokens to the database.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| successful save | Tests successful token save | ActivationTokenRequest with email, token_hash, type="activation", expiry=30 | No error |
| db error | Tests database error handling | ActivationTokenRequest with email, token_hash, type="activation", expiry=30 | Error: sql.ErrConnDone |

### 3. TestValidateActivationToken

Tests validation of activation tokens.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid token | Tests validation of valid token | ActivationTokenRequest with email, token_hash, type="activation" | No error |
| invalid token | Tests validation of invalid token | ActivationTokenRequest with email, token_hash, type="activation" | Error: sql.ErrNoRows |

### 4. TestCreateUser

Tests user creation functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| successful create | Tests successful user creation | CreateUserRequest with valid data | No error |
| transaction error | Tests handling of transaction errors | CreateUserRequest with valid data | Error: sql.ErrConnDone |

### 5. TestGetUserByEmail

Tests retrieving user information by email.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| user exists | Tests retrieval of existing user | email: "test@example.com" | User object, No error |
| user not found | Tests retrieval of non-existent user | email: "notfound@example.com" | nil, Error: sql.ErrNoRows |

### 6. TestUpdateUserPassword

Tests password update functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| successful update | Tests successful password update | UpdatePasswordRequest with valid data | No error |
| user not found | Tests update for non-existent user | UpdatePasswordRequest with invalid email | Error: "user not found" |

## Running the Tests

```bash
go test -v ./routes/auth
```

## Test Coverage

The tests cover the following features:

1. User Management
   - Email existence checking
   - User creation with secure password hashing
   - User retrieval
   - Password updates with secure hashing

2. Token Management
   - Secure token hash storage
   - Token validation with hash comparison
   - Token cleanup after use

3. Error Handling
   - Database connection errors
   - Transaction errors
   - Not found errors
   - Invalid data errors

## Implementation Details

1. Database Mocking
   - Uses `go-sqlmock` for database interaction simulation
   - Mocks both successful and error scenarios
   - Validates SQL queries and parameters for security

2. Testing Strategy
   - Table-driven tests for comprehensive coverage
   - Transaction testing with commit and rollback scenarios
   - Security validation for token and password handling
   - Clean separation of test cases

3. Assertions
   - Uses `testify/assert` package
   - Focuses on Nil/NotNil and Equal assertions
   - Validates both successful operations and error cases