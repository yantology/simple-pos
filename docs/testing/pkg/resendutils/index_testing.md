# ResendUtils Package Tests

This document describes the test cases for the `resendutils` package in the retail-pro-be application.

## Test Overview

The ResendUtils package provides email sending functionality using the Resend API. These tests verify the initialization of the ResendUtils client and email sending functionality.

## Test Files

- `pkg/resendutils/resendutils_test.go`: Contains tests for the ResendUtils implementation

## Test Suites

### 1. TestNewResendUtils

Tests the creation of new ResendUtils instances.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid initialization | Tests creation with valid parameters | apiKey: "test-api-key"<br>domain: "test.com"<br>name: "Test Sender" | Non-nil ResendUtils instance |

### 2. TestResendUtils_Send

Tests the email sending functionality.

| Test Case | Description | Input | Expected Output |
|-----------|-------------|-------|----------------|
| valid email data | Tests sending with valid parameters | html: "<h1>Test Email</h1>"<br>subject: "Test Subject"<br>to: ["test@example.com"] | No error |
| empty recipient | Tests sending with empty recipient list | html: "<h1>Test Email</h1>"<br>subject: "Test Subject"<br>to: [] | Error |

## Running the Tests

```bash
go test -v ./pkg/resendutils
```

## Test Coverage

1. ResendUtils Creation
   - Initialization with valid parameters
   - Parameter validation

2. Email Sending
   - Valid email sending
   - Error handling for invalid inputs
   - Recipient validation