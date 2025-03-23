# Auth Repository Tests

This document describes the test cases for the `repository.go` file in the `auth` package of the retail-pro-be application.

## Test Overview

The tests verify the functionality of the `AuthRepository` which handles authentication-related data access operations. The repository acts as a wrapper around database operations and provides an interface for the service layer.

## Test Files

- `repository_test.go`: Contains tests for the `AuthRepository` struct methods

## Test Suites

### 1. TestRepositoryCheckExistingEmail

Tests the `CheckExistingEmail` method of the `AuthRepository` struct which is used to verify if an email already exists in the system.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - Email doesn't exist | Tests when an email doesn't exist in the database | email: "new@example.com" | No error |
| Failure case - Email already exists | Tests when an email already exists in the database | email: "existing@example.com" | CustomError with 400 status code |

### 2. TestRepositorySaveActivationToken

Tests the `SaveActivationToken` method of the `AuthRepository` struct which saves activation tokens to the database.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - Token saved | Tests successful token saving | ActivationTokenRequest with valid data | No error |
| Failure case - Database error | Tests when the database returns an error | ActivationTokenRequest with valid data | CustomError with 500 status code |

### 3. TestRepositoryValidateActivationToken

Tests the `ValidateActivationToken` method of the `AuthRepository` struct which validates if a token exists and is valid.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - Valid token | Tests when a token is valid | ActivationTokenRequest with valid token | No error |
| Failure case - Invalid token | Tests when a token is invalid or expired | ActivationTokenRequest with invalid token | CustomError with 400 status code |

### 4. TestRepositoryCreateUser

Tests the `CreateUser` method of the `AuthRepository` struct which creates a new user in the database.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - User created | Tests successful user creation | CreateUserRequest with valid data | No error |
| Failure case - Database error | Tests when the database returns an error during user creation | CreateUserRequest with valid data | CustomError with 500 status code |

### 5. TestRepositoryGetUserByEmail

Tests the `GetUserByEmail` method of the `AuthRepository` struct which retrieves a user by their email address.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - User found | Tests when a user is found with the provided email | email: "user@example.com" | User object, No error |
| Failure case - User not found | Tests when no user exists with the provided email | email: "nonexistent@example.com" | nil, CustomError with 404 status code |

### 6. TestRepositoryUpdateUserPassword

Tests the `UpdateUserPassword` method of the `AuthRepository` struct which updates a user's password.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case - Password updated | Tests successful password update | UpdatePasswordRequest with valid data | No error |
| Failure case - User not found | Tests when the user doesn't exist | UpdatePasswordRequest with nonexistent email | CustomError with 404 status code |

## Running the Tests

```bash
go test -v ./routes/auth
```

Or for specific tests:

```bash
go test -v -run TestRepositoryCheckExistingEmail ./routes/auth
```

## Test Coverage

The tests cover the following functionality:

1. Email Validation
   - Checking if an email already exists

2. Token Management
   - Saving activation tokens
   - Validating activation tokens

3. User Management
   - Creating new users
   - Retrieving users by email
   - Updating user passwords

## Testing Approach

The tests use the `github.com/stretchr/testify/mock` package to mock the `AuthDBInterface`. This allows testing the repository layer in isolation without depending on an actual database. The tests verify that:

1. The repository correctly passes calls to the database layer
2. The repository correctly returns errors from the database layer
3. The repository correctly returns results from the database layer

All tests use table-driven testing to minimize code duplication and cover multiple test cases efficiently.