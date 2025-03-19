package jwt_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/retail-pro-be/pkg/jwt"
)

func TestNewJWTService(t *testing.T) {
	tests := []struct {
		name            string
		accessSecret    string
		refreshSecret   string
		accessDuration  time.Duration
		refreshDuration time.Duration
		issuer          string
	}{
		{
			name:            "with default values",
			accessSecret:    "",
			refreshSecret:   "",
			accessDuration:  0,
			refreshDuration: 0,
			issuer:          "",
		},
		{
			name:            "with custom values",
			accessSecret:    "custom-access",
			refreshSecret:   "custom-refresh",
			accessDuration:  30 * time.Minute,
			refreshDuration: 14 * 24 * time.Hour,
			issuer:          "custom-issuer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := jwt.NewJWTService(tt.accessSecret, tt.refreshSecret, tt.accessDuration, tt.refreshDuration, tt.issuer)
			assert.NotNil(t, service, "Service should not be nil")
		})
	}
}

func TestJWTService_GenerateAccesToken(t *testing.T) {
	service := jwt.NewJWTService("test-secret", "refresh-secret", time.Hour, time.Hour*24, "test-issuer")

	tests := []struct {
		name      string
		userID    string
		email     string
		userType  string
		wantError bool
	}{
		{
			name:      "valid token generation",
			userID:    "123",
			email:     "test@example.com",
			userType:  "user",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateAccesToken(tt.userID, tt.email, tt.userType)
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Validate the generated token
				claims, err := service.ValidateTokenClaims(token)
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, claims.UserID)
				assert.Equal(t, tt.email, claims.Email)
				assert.Equal(t, tt.userType, claims.Type)
			}
		})
	}
}

func TestJWTService_GenerateRefreshToken(t *testing.T) {
	service := jwt.NewJWTService("access-secret", "test-refresh", time.Hour, time.Hour*24, "test-issuer")

	tests := []struct {
		name      string
		userID    string
		email     string
		userType  string
		wantError bool
	}{
		{
			name:      "valid refresh token generation",
			userID:    "123",
			email:     "test@example.com",
			userType:  "user",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenerateRefreshToken(tt.userID, tt.email, tt.userType)
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestJWTService_ValidateTokenClaims(t *testing.T) {
	service := jwt.NewJWTService("test-secret", "refresh-secret", time.Hour, time.Hour*24, "test-issuer")

	tests := []struct {
		name       string
		setupToken func() string
		wantError  bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := service.GenerateAccesToken("123", "test@example.com", "user")
				return token
			},
			wantError: false,
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid.token.here"
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := service.ValidateTokenClaims(token)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.NotEmpty(t, claims.UserID)
				assert.NotEmpty(t, claims.Email)
				assert.NotEmpty(t, claims.Type)
			}
		})
	}
}

func TestJWTService_GetTokenClaims(t *testing.T) {
	service := jwt.NewJWTService("test-secret", "refresh-secret", time.Hour, time.Hour*24, "test-issuer")

	tests := []struct {
		name       string
		setupToken func() string
		wantError  bool
	}{
		{
			name: "valid token",
			setupToken: func() string {
				token, _ := service.GenerateAccesToken("123", "test@example.com", "user")
				return token
			},
			wantError: false,
		},
		{
			name: "invalid token",
			setupToken: func() string {
				return "invalid.token.here"
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.setupToken()
			claims, err := service.GetTokenClaims(token)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.NotEmpty(t, claims.UserID)
				assert.NotEmpty(t, claims.Email)
				assert.NotEmpty(t, claims.Type)
			}
		})
	}
}
