package auth

import "github.com/yantology/retail-pro-be/pkg/customerror"

// AuthDBInterface defines the interface for authentication database operations
type AuthDBInterface interface {
	// CheckExistingEmail checks if an email already exists in the database
	CheckExistingEmail(email string) *customerror.CustomError

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
