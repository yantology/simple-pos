package auth

import "github.com/yantology/simple-pos/pkg/customerror"

// AuthDBInterface defines the interface for authentication database operations
type AuthDBInterface interface {
	CheckIsNotExistingEmail(email string) *customerror.CustomError

	// CheckExistingEmail checks if an email already exists in the database
	CheckIsExistingEmail(email string) *customerror.CustomError

	// SaveActivationToken saves a new activation token in the database
	SaveActivationToken(req *ActivationTokenRequest) *customerror.CustomError

	// ValidateActivationToken validates if a token exists and is not expired
	GetActivationToken(req *GetActivationTokenRequest) (string, *customerror.CustomError)

	// CreateUser creates a new user in the database
	CreateUser(req *CreateUserRequest) *customerror.CustomError

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(email string) (*User, *customerror.CustomError)

	// UpdateUserPassword updates a user's password
	UpdateUserPassword(req *UpdatePasswordRequest) *customerror.CustomError
}
