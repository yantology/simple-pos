# Package Usage Documentation

This document provides information about the utility packages in the `pkg` directory and how to use them effectively in your application.

## Table of Contents

1. [CustomError Package]
2. [JWT Package]
3. [ResendUtils Package]
4. [DTO Package]

## CustomError Package

The `customerror` package provides standardized error handling for HTTP applications, with special handling for PostgreSQL database errors.

### Key Features

- Custom error types with HTTP status codes
- Specialized handling for common PostgreSQL error codes
- Consistent error response structure

### Usage

```go
import "github.com/yantology/simple-pos/pkg/customerror"

// Creating a custom error
err := customerror.NewCustomError(originalError, "User-friendly message", http.StatusBadRequest)

// For database operations, convert PostgreSQL errors
if dbErr != nil {
    return customerror.NewPostgresError(dbErr)
}

// Using the error in an HTTP handler
if err != nil {
    customErr, ok := err.(*customerror.CustomError)
    if ok {
        return c.JSON(customErr.Code(), map[string]string{
            "message": customErr.Message(),
        })
    }
    return c.JSON(http.StatusInternalServerError, map[string]string{
        "message": "Internal server error",
    })
}
```

### PostgreSQL Error Handling

The package automatically maps common PostgreSQL error codes to appropriate HTTP status codes:

| PostgreSQL Error | Code | HTTP Status | Message |
|-----------------|------|-------------|---------|
| unique_violation | 23505 | 409 Conflict | Record already exists |
| foreign_key_violation | 23503 | 400 Bad Request | Foreign key violation |
| string_data_right_truncation | 22001 | 400 Bad Request | String data is too long |
| sql.ErrNoRows | - | 404 Not Found | Original error message |
| Other database errors | - | 500 Internal Server Error | Database error |

## JWT Package

The `jwt` package provides JWT (JSON Web Token) authentication functionality for securing APIs.

### Key Features

- Generate and validate access and refresh tokens
- Configurable token expiration times
- Type-safe token claims

### Usage

```go
import (
    "time"
    "github.com/yantology/simple-pos/pkg/jwt"
)

// Initialize JWT service with custom parameters
jwtService := jwt.NewJWTService(
    "your-access-secret",       // Access token secret
    "your-refresh-secret",      // Refresh token secret
    30 * time.Minute,          // Access token duration
    7 * 24 * time.Hour,        // Refresh token duration
    "your-app-name"            // Token issuer
)

// Or use with defaults
jwtService := jwt.NewJWTService("", "", 0, 0, "")

// Generate tokens
accessToken, err := jwtService.GenerateAccesToken("user123", "user@example.com")
refreshToken, err := jwtService.GenerateRefreshToken("user123", "user@example.com")

// Validate tokens
accessClaims, err := jwtService.ValidateAccessTokenClaims(accessToken)
if err != nil {
    // Handle invalid token error
}

// Access the claims
userID := accessClaims.UserID
email := accessClaims.Email
tokenType := accessClaims.TypeToken
```

### Default Configuration

| Parameter | Default Value |
|-----------|---------------|
| Access token duration | 15 minutes |
| Refresh token duration | 7 days |
| Access secret | "access" |
| Refresh secret | "refresh" |
| Issuer | "retail-pro" |

## ResendUtils Package

The `resendutils` package provides email functionality using the [Resend](https://resend.com) service.

### Key Features

- Simple email sending interface
- HTML email support
- Error handling with CustomError integration

### Usage

```go
import "github.com/yantology/simple-pos/pkg/resendutils"

// Initialize the Resend utility
resend := resendutils.NewResendUtils(
    "re_123your_api_key",    // Resend API key
    "yourdomain.com"         // From domain
)

// Send an email
html := "<h1>Your Activation Code</h1><p>Your code is: 123456</p>"
subject := "Account Activation"
recipients := []string{"user@example.com"}

err := resend.Send(html, subject, recipients)
if err != nil {
    // Handle error
    fmt.Println(err.Message())
}
```

### Error Handling

If email sending fails, the method returns a CustomError with:
- HTTP status code: 500 (Internal Server Error)
- Error message: "Failed to send email"
- Original error: The underlying error from the Resend API

## DTO Package

The `dto` package provides standardized structures for API responses, ensuring consistency across your application.

### Key Features

- Generic data response structure for any type
- Simple message response for status or error messages
- Designed for use in HTTP handlers and services

### Usage

```go
import "github.com/yantology/simple-pos/pkg/dto"

// Returning a data response in a handler
data := YourDataType{/* ... */}
c.JSON(http.StatusOK, dto.DataResponse[YourDataType]{
    Data:    data,
    Message: "Operation completed successfully",
})

// Returning a message response
c.JSON(http.StatusOK, dto.MessageResponse{
    Message: "Operation completed successfully",
})
```

### Structures

```go
// DataResponse represents a generic data response
// @Description Generic data response model
type DataResponse[T any] struct {
    Data    T      `json:"data"`
    Message string `json:"message" example:"Operation completed successfully"`
}

// MessageResponse represents a generic message response
// @Description Generic message response model
type MessageResponse struct {
    Message string `json:"message" example:"Operation completed successfully"`
}
```

## Integration Example

Here's an example demonstrating how these packages can work together:

```go
func SendActivationEmail(userID, email, code string) error {
    // Initialize JWT service
    jwtService := jwt.NewJWTService("secret", "refresh-secret", 0, 0, "my-app")
    
    // Initialize Resend service
    resendUtils := resendutils.NewResendUtils("re_api_key", "example.com")
    
    // Generate activation token
    token, err := jwtService.GenerateAccesToken(userID, email)
    if err != nil {
        return customerror.NewCustomError(err, "Failed to generate token", http.StatusInternalServerError)
    }
    
    // Create email content
    html := fmt.Sprintf("<h1>Your Activation Code</h1><p>Code: %s</p>", code)
    
    // Send email
    if err := resendUtils.Send(html, "Activation Code", []string{email}); err != nil {
        return err // Already a customerror.CustomError
    }
    
    return nil
}
```
