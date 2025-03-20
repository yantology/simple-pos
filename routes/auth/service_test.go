package auth_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yantology/retail-pro-be/pkg/jwt"
	"github.com/yantology/retail-pro-be/routes/auth"
)

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) GenerateAccesToken(userID, email, userType string) (string, error) {
	args := m.Called(userID, email, userType)
	return args.String(0), args.Error(1)
}

func (m *mockJWTService) GenerateRefreshToken(userID, email, userType string) (string, error) {
	args := m.Called(userID, email, userType)
	return args.String(0), args.Error(1)
}

func (m *mockJWTService) ValidateTokenClaims(token string) (*jwt.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenClaims), args.Error(1)
}

func (m *mockJWTService) GetTokenClaims(token string) (*jwt.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenClaims), args.Error(1)
}

func TestNewAuthService(t *testing.T) {
	jwtService := &mockJWTService{}
	service := auth.NewAuthService(jwtService)
	assert.NotNil(t, service)
}

func TestValidateEmail(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	tests := []struct {
		name         string
		email        string
		wantErr      bool
		expectedCode int
		expectedMsg  string
	}{
		{
			name:    "Valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:         "Empty email",
			email:        "",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Format email tidak valid",
		},
		{
			name:         "Invalid email format",
			email:        "invalid-email",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Format email tidak valid",
		},
		{
			name:         "Email without domain",
			email:        "test@",
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Format email tidak valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.ValidateEmail(tt.email)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedCode, err.Code())
				assert.Equal(t, tt.expectedMsg, err.Message())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGenerateActivationToken(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	token, err := authService.GenerateActivationToken()
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestValidateRegistrationInput(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	tests := []struct {
		name         string
		req          auth.RegistrationRequest
		wantErr      bool
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Valid registration input",
			req: auth.RegistrationRequest{
				Email:                "test@example.com",
				Username:             "testuser",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			req: auth.RegistrationRequest{
				Email:                "invalid-email",
				Username:             "testuser",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Format email tidak valid",
		},
		{
			name: "Username too long",
			req: auth.RegistrationRequest{
				Email:                "test@example.com",
				Username:             "thisusernameiswaytoolongandshouldfail",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Username maksimal 30 karakter",
		},
		{
			name: "Password mismatch",
			req: auth.RegistrationRequest{
				Email:                "test@example.com",
				Username:             "testuser",
				Password:             "password123",
				PasswordConfirmation: "password456",
			},
			wantErr:      true,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Password tidak cocok dengan konfirmasi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := authService.ValidateRegistrationInput(tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedCode, err.Code())
				assert.Equal(t, tt.expectedMsg, err.Message())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	password := "testpassword123"

	// Test password hashing
	hashedPassword, err := authService.HashString(password)
	assert.Nil(t, err)
	assert.NotNil(t, hashedPassword)

	// Test password verification
	err = authService.VerifyHash(hashedPassword, password)
	assert.Nil(t, err)

	// Test password verification with wrong password
	err = authService.VerifyHash(hashedPassword, "wrongpassword")
	assert.NotNil(t, err)
	assert.Equal(t, "Email atau password salah", err.Message())
	assert.Equal(t, http.StatusUnauthorized, err.Code())
}

func TestHashAndVerifyString(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	tests := []struct {
		name          string
		input         string
		wantHashErr   bool
		wantVerifyErr bool
	}{
		{
			name:          "Hash and verify password",
			input:         "testpassword123",
			wantHashErr:   false,
			wantVerifyErr: false,
		},
		{
			name:          "Hash and verify token",
			input:         "123456",
			wantHashErr:   false,
			wantVerifyErr: false,
		},
		{
			name:          "Empty string",
			input:         "",
			wantHashErr:   false,
			wantVerifyErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test string hashing
			hashedString, err := authService.HashString(tt.input)
			if tt.wantHashErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotEqual(t, tt.input, hashedString)

				// Test hash verification with correct input
				err = authService.VerifyHash(hashedString, tt.input)
				if tt.wantVerifyErr {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}

				// Test hash verification with wrong input
				err = authService.VerifyHash(hashedString, "wrong"+tt.input)
				assert.NotNil(t, err)
				assert.Equal(t, "Hash tidak cocok", err.Message())
				assert.Equal(t, http.StatusUnauthorized, err.Code())
			}
		})
	}
}

func TestGenerateTokenPair(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	tests := []struct {
		name      string
		req       auth.TokenPairRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "Successful token generation",
			req: auth.TokenPairRequest{
				UserID:   "123",
				Email:    "test@example.com",
				UserType: "user",
			},
			mockSetup: func() {
				jwtService.On("GenerateAccesToken", "123", "test@example.com", "user").
					Return("access-token", nil)
				jwtService.On("GenerateRefreshToken", "123", "test@example.com", "user").
					Return("refresh-token", nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwtService.ExpectedCalls = nil
			jwtService.Calls = nil

			tt.mockSetup()

			response, err := authService.GenerateTokenPair(tt.req)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, "", response.AccessToken)
				assert.Equal(t, "", response.RefreshToken)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, response.AccessToken)
				assert.NotNil(t, response.RefreshToken)
			}
		})
	}
}

func TestValidateTokenClaims(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	validToken := "valid.jwt.token"
	mockClaims := &jwt.TokenClaims{
		UserID: "123",
		Email:  "test@example.com",
		Type:   "user",
	}

	jwtService.On("ValidateTokenClaims", validToken).Return(mockClaims, nil)

	claims, err := authService.ValidateTokenClaims(validToken)
	assert.Nil(t, err)
	assert.Equal(t, mockClaims, claims)
}

func TestGenerateAccessToken(t *testing.T) {
	jwtService := &mockJWTService{}
	authService := auth.NewAuthService(jwtService)

	jwtService.On("GenerateAccesToken", "123", "test@example.com", "user").
		Return("new-access-token", nil)

	token, err := authService.GenerateAccessToken("123", "test@example.com", "user")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}
