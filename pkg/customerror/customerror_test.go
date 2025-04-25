package customerror_test

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/yantology/simple-pos/pkg/customerror"
)

func TestNewCustomError(t *testing.T) {
	tests := []struct {
		name         string
		original     error
		message      string
		httpCode     int
		wantMsg      string
		wantCode     int
		wantOriginal string
	}{
		{
			name:         "basic error",
			original:     errors.New("original error"),
			message:      "test error",
			httpCode:     http.StatusBadRequest,
			wantMsg:      "test error",
			wantCode:     http.StatusBadRequest,
			wantOriginal: "original error",
		},
		{
			name:         "nil original error",
			original:     nil,
			message:      "test error with nil original",
			httpCode:     http.StatusInternalServerError,
			wantMsg:      "test error with nil original",
			wantCode:     http.StatusInternalServerError,
			wantOriginal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := customerror.NewCustomError(tt.original, tt.message, tt.httpCode)

			assert.Equal(t, tt.wantMsg, err.Message(), "Message should match")
			assert.Equal(t, tt.wantCode, err.Code(), "Code should match")
			assert.Equal(t, tt.wantOriginal, err.Original(), "Original message should match")
		})
	}
}

func TestNewPostgresError(t *testing.T) {
	createPqError := func(code string) error {
		return &pq.Error{Code: pq.ErrorCode(code)}
	}

	tests := []struct {
		name     string
		err      error
		wantMsg  string
		wantCode int
	}{
		{
			name:     "nil error",
			err:      nil,
			wantMsg:  "",
			wantCode: 0,
		},
		{
			name:     "unique violation",
			err:      createPqError("23505"),
			wantMsg:  "Record already exists",
			wantCode: http.StatusConflict,
		},
		{
			name:     "foreign key violation",
			err:      createPqError("23503"),
			wantMsg:  "Foreign key violation",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "string data too long",
			err:      createPqError("22001"),
			wantMsg:  "String data is too long",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "unhandled postgres error",
			err:      createPqError("42P01"), // relation does not exist
			wantMsg:  createPqError("42P01").Error(),
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "record not found",
			err:      sql.ErrNoRows,
			wantMsg:  "Record not found",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "generic error",
			err:      errors.New("some generic error"),
			wantMsg:  "Dsome generic error",
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customErr := customerror.NewPostgresError(tt.err)

			if tt.err == nil {
				assert.Nil(t, customErr, "Should return nil for nil error")
				return
			}

			assert.NotNil(t, customErr, "Should return a custom error")
			assert.Equal(t, tt.wantMsg, customErr.Message(), "Message should match")
			assert.Equal(t, tt.wantCode, customErr.Code(), "Code should match")
			assert.NotEmpty(t, customErr.Original(), "Original error should be preserved")
		})
	}
}
