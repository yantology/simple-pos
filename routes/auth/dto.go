package auth

// Request DTOs

// TokenRequest represents a request for an activation token
type TokenRequest struct {
	Email string `json:"email" binding:"required"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email                string `json:"email" binding:"required"`
	Fullname             string `json:"fullname" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	ActivationCode       string `json:"activation_code" binding:"required"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ForgetPasswordRequest represents a password reset request
type ForgetPasswordRequest struct {
	Email                   string `json:"email" binding:"required"`
	NewPassword             string `json:"new_password" binding:"required"`
	NewPasswordConfirmation string `json:"new_password_confirmation" binding:"required"`
	ActivationCode          string `json:"activation_code" binding:"required"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// LoginResponseData represents login response data
type jwtResponseData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// LoginResponse represents a login response
type DataResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
