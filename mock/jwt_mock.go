package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/yantology/retail-pro-be/pkg/jwt"
)

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateAccesToken(userID, email, userType string) (string, error) {
	args := m.Called(userID, email, userType)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID, email, userType string) (string, error) {
	args := m.Called(userID, email, userType)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateTokenClaims(token string) (*jwt.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenClaims), args.Error(1)
}

func (m *MockJWTService) GetTokenClaims(token string) (*jwt.TokenClaims, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.TokenClaims), args.Error(1)
}
