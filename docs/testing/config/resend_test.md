# Resend Configuration Tests

This document describes the test cases for the resend configuration in the retail-pro-be application.

## Test Overview

These tests verify that the resend configuration is properly loaded from environment variables and that appropriate errors are returned when required configuration is missing.

## Test Files

- `config/resend_test.go`: Contains tests for resend configuration initialization

## Test Suites

### 1. TestInitResendConfig

Tests the initialization of resend configuration with different environment variable scenarios.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| with valid configuration | Tests config with all env vars set | RESEND_API_KEY="test-api-key"<br>RESEND_DOMAIN="test.com"<br>RESEND_NAME="Test Sender" | Valid config, no error |
| with missing API key | Tests config with no env vars | No env vars | Error, nil config |
| with missing domain | Tests config with only API key | RESEND_API_KEY="test-api-key" | Error, nil config |
| with missing name | Tests config with missing name | RESEND_API_KEY="test-api-key"<br>RESEND_DOMAIN="test.com" | Error, nil config |

## Running the Tests

```bash
go test -v ./config
```

## Test Coverage

1. Configuration Loading
   - Valid configuration with all parameters
   - Missing required parameters
   - Environment variable handling

2. Error Handling
   - Missing API key validation
   - Missing domain validation
   - Missing name validation
   - Proper error messages
   - Error code validation