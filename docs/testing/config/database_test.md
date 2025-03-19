# Database Configuration Tests

This document describes the test cases for the database configuration in the retail-pro-be application.

## Test Overview

The tests verify that the database configuration is properly loaded from environment variables and that database connections can be established correctly.

## Test Files

- `config/database_test.go`: Contains tests for database configuration and connection

## Test Suites

### 1. TestInitDatabaseConfig

Tests the initialization of database configuration with different environment variable scenarios.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| with default values | Tests configuration with no environment variables set | No env vars | Default MySQL configuration |
| with custom values | Tests configuration with all environment variables set | Full set of DB env vars | Custom PostgreSQL configuration |

### 2. TestConnectDatabase

Tests the database connection functionality with different drivers.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| mysql connection | Tests MySQL connection string formation | MySQL config | Valid connection |
| postgres connection | Tests PostgreSQL connection string formation | PostgreSQL config | Valid connection |

## Running the Tests

```bash
go test -v ./config
```

## Test Coverage

1. Database Configuration
   - Default values
   - Environment variable overrides
   - MySQL configuration
   - PostgreSQL configuration

2. Database Connection
   - MySQL connection string formation
   - PostgreSQL connection string formation
   - Connection error handling
