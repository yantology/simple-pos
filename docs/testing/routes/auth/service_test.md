# Auth Service Tests

This document describes the test cases for the auth service in the retail-pro-be application.

## Test Overview

These tests cover the authentication service functionality, including email validation, password management, token generation and validation.

## Test Files

- `service_test.go`: Contains all service-layer tests for the auth package

## Test Suites

### 1. TestNewAuthService

Tests the creation of a new auth service instance.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Create new service | Tests service instantiation | JWT service mock | Non-nil service instance |

### 2. TestValidateEmail

Tests email validation functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Valid email | Tests valid email format | "test@example.com" | No error |
| Empty email | Tests empty email | "" | Error with code 400 |
| Invalid format | Tests invalid email format | "invalid-email" | Error with code 400 |
| Missing domain | Tests email without domain | "test@" | Error with code 400 |

### 3. TestGenerateActivationToken

Tests generation of activation tokens.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Generate token | Tests token generation | None | 6-digit token, no error |

### 4. TestValidateRegistrationInput

Tests registration input validation.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Valid input | Tests valid registration data | Valid RegistrationRequest | No error |
| Invalid email | Tests invalid email | Invalid email in request | Error with code 400 |
| Long username | Tests too long username | >30 chars username | Error with code 400 |
| Password mismatch | Tests mismatched passwords | Different passwords | Error with code 400 |

### 5. TestHashAndVerifyString

Tests the flexible string hashing and verification functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Hash and verify password | Tests password hashing/verification | "testpassword123" | Hash verification succeeds |
| Hash and verify token | Tests token hashing/verification | "123456" | Hash verification succeeds |
| Empty string | Tests empty string handling | "" | Hash verification succeeds |
| Wrong input verification | Tests verification with incorrect input | Wrong string | Error with code 401 |

### 6. TestGenerateTokenPair

Tests JWT token pair generation.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case | Tests successful token generation | Valid TokenPairRequest | Access and refresh tokens |

### 7. TestValidateTokenClaims

Tests JWT token claims validation.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Valid token | Tests valid token claims | Valid JWT token | Valid claims object |

### 8. TestGenerateAccessToken

Tests access token generation.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| Success case | Tests successful token generation | User ID, email, type | Valid access token |

## Running the Tests

```bash
go test -v ./routes/auth
```

## Test Coverage

The tests cover:

1. String Security
   - Generic string hashing (passwords/tokens)
   - Hash verification
   - Security error handling

2. User Management
   - Email validation
   - Registration input validation
   - Password validation

3. Token Management
   - Activation token generation
   - JWT token pair generation
   - Token claims validation

## Implementation Details

1. Mocking
   - Uses testify/mock for JWT service mocking
   - Simulates JWT operations without real tokens

2. Test Structure
   - Table-driven tests
   - Clear test cases with descriptive names
   - Proper error checking

3. Assertions
   - Uses testify/assert package
   - Checks for nil/not-nil conditions
   - Validates error messages and codes