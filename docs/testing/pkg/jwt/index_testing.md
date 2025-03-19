# JWT Package Tests

This document describes the test cases for the `jwt` package in the retail-pro-be application.

## Test Overview

The JWT package provides JSON Web Token functionality for authentication and authorization. These tests verify the token generation, validation, and claims extraction functionality.

## Test Files

- `jwt_test.go`: Contains all test cases for the JWT service implementation
- `jwt_mock.go`: Contains the mock implementation for the JWTService interface

## Test Suites

### 1. TestNewJWTService

Tests the creation of new JWT service instances.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| with default values | Creates service with default configuration | Empty values | Valid service instance with default values |
| with custom values | Creates service with custom configuration | Custom secrets, durations, issuer | Valid service instance with custom values |

### 2. TestJWTService_GenerateAccesToken

Tests the generation of access tokens.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid token generation | Generates access token with user details | userID, email, userType | Valid JWT token string |

### 3. TestJWTService_GenerateRefreshToken

Tests the generation of refresh tokens.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid refresh token generation | Generates refresh token with user details | userID, email, userType | Valid JWT token string |

### 4. TestJWTService_ValidateTokenClaims

Tests token validation functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid token | Validates a correctly formed token | Valid JWT token | Token claims, no error |
| invalid token | Attempts to validate an invalid token | Invalid token string | Error |

### 5. TestJWTService_GetTokenClaims

Tests extraction of claims from tokens.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid token | Extracts claims from valid token | Valid JWT token | Token claims, no error |
| invalid token | Attempts to extract claims from invalid token | Invalid token string | Error |

## Running the Tests

```bash
go test -v ./pkg/jwt
```

## Test Coverage

1. Token Generation
   - Access token generation
   - Refresh token generation
   - Token expiration
   - Token issuer

2. Token Validation
   - Valid token validation
   - Invalid token handling
   - Token claims extraction

3. Service Configuration
   - Default configuration
   - Custom configuration
   - Parameter validation

## Mock Usage

The mock implementation (`MockJWTService`) can be used in other packages' tests where JWT functionality is required. Example usage:

```go
mockJWT := &mock.MockJWTService{}
mockJWT.On("GenerateAccesToken", "123", "test@example.com", "user").Return("token", nil)
```