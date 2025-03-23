package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yantology/retail-pro-be/pkg/customerror"
	jwtPkg "github.com/yantology/retail-pro-be/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	// Email validation and token generation
	ValidateEmail(email string) *customerror.CustomError
	GenerateActivationToken() (string, *customerror.CustomError)
	// User registration and password management
	ValidateRegistrationInput(req RegistrationRequest) *customerror.CustomError
	HashString(input string) (string, *customerror.CustomError)
	VerifyHash(hashedString, input string) *customerror.CustomError
	ValidatePasswordInput(password, passwordConfirmation string) *customerror.CustomError

	// Token operations
	GenerateTokenPair(req TokenPairRequest) (TokenPairResponse, *customerror.CustomError)
	ValidateRefreshTokenClaims(token string) (*jwtPkg.TokenClaims, *customerror.CustomError)
}

type authService struct {
	jwtService jwtPkg.JWTService
}

// NewAuthService creates a new instance of the AuthService
func NewAuthService(jwtService jwtPkg.JWTService) AuthService {
	// Compile email regex once during initialization
	return &authService{
		jwtService: jwtService,
	}
}

// ValidateEmail checks if the email format is valid
func (s *authService) ValidateEmail(email string) *customerror.CustomError {
	parsedEmail, err := mail.ParseAddress(email)
	if err != nil {
		return customerror.NewCustomError(err, "Format email tidak valid", http.StatusBadRequest)
	}
	if parsedEmail.Address == "" {
		return customerror.NewCustomError(nil, "Email tidak boleh kosong", http.StatusBadRequest)
	}
	return nil
}

// GenerateActivationToken generates a 6-digit token for activation
func (s *authService) GenerateActivationToken() (string, *customerror.CustomError) {
	// Generate a 6-digit numeric token
	token := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	return token, nil
}

// ValidateRegistrationInput validates user registration input
func (s *authService) ValidateRegistrationInput(req RegistrationRequest) *customerror.CustomError {
	// Validate email
	if err := s.ValidateEmail(req.Email); err != nil {
		return err
	}

	// Validate username length
	if len(req.Username) > 30 {
		return customerror.NewCustomError(nil, "Username maksimal 30 karakter", http.StatusBadRequest)
	}

	// Validate password match
	if req.Password != req.PasswordConfirmation {
		return customerror.NewCustomError(nil, "Password tidak cocok dengan konfirmasi", http.StatusBadRequest)
	}

	return nil
}

// HashString securely hashes a string using bcrypt
func (s *authService) HashString(input string) (string, *customerror.CustomError) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", customerror.NewCustomError(err, "Gagal mengenkripsi string", http.StatusInternalServerError)
	}
	return string(hashedString), nil
}

// VerifyHash verifies if the provided input matches the stored hash
func (s *authService) VerifyHash(hashedString, input string) *customerror.CustomError {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(input))
	if err != nil {
		return customerror.NewCustomError(err, "Hash tidak cocok", http.StatusUnauthorized)
	}
	return nil
}

// GenerateTokenPair generates an access token and refresh token pair
func (s *authService) GenerateTokenPair(req TokenPairRequest) (TokenPairResponse, *customerror.CustomError) {
	accessToken, err := s.jwtService.GenerateAccesToken(req.UserID, req.Email)
	if err != nil {
		return TokenPairResponse{}, customerror.NewCustomError(err, "Gagal membuat access token", http.StatusInternalServerError)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(req.UserID, req.Email)
	if err != nil {
		return TokenPairResponse{}, customerror.NewCustomError(err, "Gagal membuat refresh token", http.StatusInternalServerError)
	}

	return TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidatePasswordInput validates password reset input
func (s *authService) ValidatePasswordInput(password, passwordConfirmation string) *customerror.CustomError {
	if password != passwordConfirmation {
		return customerror.NewCustomError(nil, "Password baru tidak cocok dengan konfirmasi", http.StatusBadRequest)
	}
	return nil
}

// ValidateTokenClaims validates and extracts claims from a JWT token
func (s *authService) ValidateRefreshTokenClaims(token string) (*jwtPkg.TokenClaims, *customerror.CustomError) {
	claims, err := s.jwtService.ValidateRefreshTokenClaims(token)
	if err != nil {
		var message string
		var statusCode int

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "Token sudah kadaluarsa"
			} else {
				log.Println("Token tidak valid:", err)
				message = "Token tidak valid"
			}
			statusCode = http.StatusUnauthorized
		} else {
			message = "Gagal memvalidasi token"
			statusCode = http.StatusInternalServerError
		}

		return nil, customerror.NewCustomError(err, message, statusCode)
	}
	return claims, nil
}
