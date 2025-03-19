# JWT Configuration Tests

This document describes the test cases for the JWT configuration in the retail-pro-be application.

## Test Overview

These tests verify that the JWT configuration properly handles various environment variable configurations and validates required fields.

## Test Files

- `config/jwt_test.go`: Contains tests for JWT configuration initialization

## Test Suites

### 1. TestInitJWTConfig

Tests the initialization of JWT configuration with different scenarios.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid configuration | Tests with all env vars set | JWT_ACCESS_SECRET="test-access-secret"<br>JWT_REFRESH_SECRET="test-refresh-secret"<br>JWT_ACCESS_DURATION_MINUTES="30"<br>JWT_REFRESH_DURATION_DAYS="14"<br>JWT_ISSUER="test-issuer" | Valid config with specified values |
| missing access secret | Tests without access secret | JWT_REFRESH_SECRET="test-refresh-secret" | Error, nil config |
| missing refresh secret | Tests without refresh secret | JWT_ACCESS_SECRET="test-access-secret" | Error, nil config |
| default durations | Tests with only required secrets | JWT_ACCESS_SECRET="test-access-secret"<br>JWT_REFRESH_SECRET="test-refresh-secret" | Config with default durations and issuer |
| invalid duration values | Tests with invalid duration values | JWT_ACCESS_SECRET="test-access-secret"<br>JWT_REFRESH_SECRET="test-refresh-secret"<br>JWT_ACCESS_DURATION_MINUTES="invalid"<br>JWT_REFRESH_DURATION_DAYS="invalid" | Config with default durations |

## Running the Tests

```bash
go test -v ./config -run "TestJWT.*"
```

## Test Coverage

1. Required Configuration
   - Access secret validation
   - Refresh secret validation
   - Error handling for missing required fields

2. Optional Configuration
   - Access token duration
   - Refresh token duration
   - Issuer name
   - Default values handling

3. Duration Parsing
   - Valid duration values
   - Invalid duration values
   - Default duration fallback

4. Error Handling
   - Missing required fields
   - Invalid configuration
   - Error message validation