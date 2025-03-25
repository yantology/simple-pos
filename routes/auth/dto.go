package auth

// Request DTOs

// TokenRequest represents the request for generating activation tokens
// @Description Token request model
type TokenRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// RegisterRequest represents the registration request
// @Description Registration request model
type RegisterRequest struct {
	Email                string `json:"email" binding:"required,email" example:"user@example.com"`
	Fullname             string `json:"fullname" binding:"required" example:"John Doe"`
	Password             string `json:"password" binding:"required" example:"securePassword123"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required" example:"securePassword123"`
	ActivationCode       string `json:"activation_code" binding:"required" example:"123456"`
}

// LoginRequest represents the login request
// @Description Login request model
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"securePassword123"`
}

// ForgetPasswordRequest represents the password reset request
// @Description Password reset request model
type ForgetPasswordRequest struct {
	Email                   string `json:"email" binding:"required,email" example:"user@example.com"`
	ActivationCode          string `json:"activation_code" binding:"required" example:"123456"`
	NewPassword             string `json:"new_password" binding:"required" example:"newSecurePassword123"`
	NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required" example:"newSecurePassword123"`
}

// RefreshTokenRequest represents the refresh token request
// @Description Refresh token request model
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// MessageResponse represents a generic message response
// @Description Generic message response model
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// JWTResponseData represents the JWT token response data
// @Description JWT token response data model
type JWTResponseData struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// DataResponse represents a generic data response
// @Description Generic data response model
type DataResponse[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message" example:"Operation completed successfully"`
}
