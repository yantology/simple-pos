# CORS Configuration Tests

This document describes the test cases for the CORS configuration in the retail-pro-be application.

## Test Overview

These tests verify that the CORS configuration properly handles different origin scenarios and environment variable parsing.

## Test Files

- `config/cros_test.go`: Contains tests for CORS configuration initialization and helper functions

## Test Suites

### 1. TestCorsConfig

Tests the CORS configuration initialization with different origin scenarios.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| default origins | Tests with no environment variable set | CORS_ALLOW_ORIGINS="" | Single origin "*" |
| single custom origin | Tests with one origin | CORS_ALLOW_ORIGINS="http://localhost:3000" | Single specified origin |
| multiple custom origins | Tests with multiple origins | CORS_ALLOW_ORIGINS="http://localhost:3000,http://example.com" | Two specified origins |

### 2. TestGetEnvAsSlice

Tests the helper function that parses environment variables into string slices.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| empty environment variable | Tests with no env var set | TEST_SLICE="" | Default value |
| single value | Tests with single value | TEST_SLICE="value1" | ["value1"] |
| multiple values | Tests with comma-separated values | TEST_SLICE="value1,value2,value3" | ["value1", "value2", "value3"] |

## Running the Tests

```bash
go test -v ./config -run "TestCors.*"
```

## Test Coverage

1. CORS Configuration
   - Default configuration handling
   - Single origin configuration
   - Multiple origins configuration
   - CORS middleware creation

2. Environment Variable Parsing
   - Empty environment variables
   - Single value parsing
   - Multiple value parsing
   - Default value handling