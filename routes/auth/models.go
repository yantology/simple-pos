package auth

import "time"

// TokenPairRequest represents the input parameters for generating token pairs
type TokenPairRequest struct {
	UserID string
	Email  string
}

// TokenPairResponse represents the response structure for token pairs
type TokenPairResponse struct {
	AccessToken  string
	RefreshToken string
}

// RegistrationRequest represents the input parameters for user registration
type RegistrationRequest struct {
	Email                string
	Username             string
	Password             string
	PasswordConfirmation string
}

// User represents the user data structure from database
type User struct {
	ID           string
	Email        string
	Fullname     string
	PasswordHash string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

// ActivationTokenRequest represents input for token activation operations
type ActivationTokenRequest struct {
	Email          string
	TokenType      string
	ActivationCode string
	ExpiryMinutes  int
}

type GetActivationTokenRequest struct {
	Email     string
	TokenType string
}

// CreateUserRequest represents input for creating a new user
type CreateUserRequest struct {
	Email        string
	Fullname     string
	PasswordHash string
}

// UpdatePasswordRequest represents input for updating user password
type UpdatePasswordRequest struct {
	Email           string
	NewPasswordHash string
}
