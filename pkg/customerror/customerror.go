package customerror

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

type CustomError struct {
	httpCode int
	message  string
	original error
}

// NewCustomError creates a new custom error
func NewCustomError(original error, message string, httpCode int) *CustomError {
	return &CustomError{
		httpCode: httpCode,
		message:  message,
		original: original,
	}
}

// Error implements the error interface
func (ce *CustomError) Message() string {
	return ce.message
}

// Helper method to extract original database error message
func (ce *CustomError) Original() string {
	if ce.original != nil {
		return ce.original.Error()
	}
	return ""
}

// Helper method to extract original database error code
func (ce *CustomError) Code() int {
	// Even if Original is nil, we should return the HTTPCode
	return ce.httpCode
}

// NewPostgresError creates a custom error from PostgreSQL errors
func NewPostgresError(err error) *CustomError {
	if err == nil {
		return nil
	}

	// Handle specific PostgreSQL errors
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505": // unique_violation
			return NewCustomError(err, "Record already exists", http.StatusConflict)
		case "23503": // foreign_key_violation
			return NewCustomError(err, "Foreign key violation", http.StatusBadRequest)
		case "22001": // string_data_right_truncation
			return NewCustomError(err, "String data is too long", http.StatusBadRequest)
		}
	}

	// Generic database error
	if errors.Is(err, sql.ErrNoRows) {
		return NewCustomError(err, err.Error(), http.StatusNotFound)
	}

	// Default error handling
	return NewCustomError(err, "Database error", http.StatusInternalServerError)
}
