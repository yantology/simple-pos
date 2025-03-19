# App Configuration Tests

This document describes the test cases for the app configuration in the retail-pro-be application.

## Test Overview

The tests verify that the application configuration is properly loaded from environment variables with appropriate fallback to default values.

## Test Files

- `config/app_test.go`: Contains tests for application configuration initialization

## Test Suites

### 1. TestInitAppConfig

Tests the initialization of application configuration with different environment variable scenarios.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| with default values | Tests configuration with no environment variables set | No env vars | Port: "3000", PublicRoute: "/public", PublicAssetsDir: "./public" |
| with custom port | Tests configuration with custom port set | APP_PORT="8080" | Port: "8080", PublicRoute: "/public", PublicAssetsDir: "./public" |

## Running the Tests

```bash
go test -v ./config
```

## Test Coverage

1. Default Configuration
   - Default port value
   - Default public route
   - Default public assets directory

2. Environment Variable Override
   - Custom port override
   - Environment variable parsing
