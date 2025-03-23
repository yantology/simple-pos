package jwt

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenClaims represents the claims in a JWT token.
type TokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TypeToken string `json:"type_token"`
	jwt.StandardClaims
}

type JWTService interface {
	GenerateAccesToken(userID, email string) (string, error)
	GenerateRefreshToken(userID, email string) (string, error)
	ValidateAccessTokenClaims(token string) (*TokenClaims, error)
	ValidateRefreshTokenClaims(token string) (*TokenClaims, error)
}

// Berbagai konstanta dan error yang sering digunakan
const (
	// DefaultAccessTokenExpiration adalah durasi default untuk access token (15 menit)
	DefaultAccessDuration = 15 * time.Minute
	// DefaultRefreshTokenExpiration adalah durasi default untuk refresh token (7 hari)
	DefaultRefreshDuration = 7 * 24 * time.Hour
	// TokenTypeAccess adalah type untuk access token
	DefaultaccessSecret = "access"
	// TokenTypeRefresh adalah type untuk refresh token
	DefaultrefreshSecret = "refresh"
)

// NewJWTService creates a new JWT service with the provided secret keys and expiration times.

const (
	// Add default issuer
	DefaultIssuer = "retail-pro"
)

type jwtService struct {
	accessSecret   string
	refreshSecret  string
	accessDuration time.Duration
	refresDuration time.Duration
	issuer         string
}

func NewJWTService(accessSecret, refreshSecret string, accessDuration, refresDuration time.Duration, issuer string) JWTService {
	if accessSecret == "" {
		accessSecret = DefaultaccessSecret
	}

	if refreshSecret == "" {
		refreshSecret = DefaultrefreshSecret
	}

	if accessDuration == 0 {
		accessDuration = DefaultAccessDuration
	}
	if refresDuration == 0 {
		refresDuration = DefaultRefreshDuration
	}
	if issuer == "" {
		issuer = DefaultIssuer
	}
	return &jwtService{
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessDuration: accessDuration,
		refresDuration: refresDuration,
		issuer:         issuer,
	}
}

func (j *jwtService) GenerateAccesToken(userID, email string) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		Email:     email,
		TypeToken: "access",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.accessDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

func (j *jwtService) GenerateRefreshToken(userID, email string) (string, error) {
	claims := TokenClaims{
		UserID:    userID,
		Email:     email,
		TypeToken: "refresh",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.refresDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    j.issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *jwtService) ValidateAccessTokenClaims(token string) (*TokenClaims, error) {
	return j.GetTokenClaims(token, j.accessSecret)
}

func (j *jwtService) ValidateRefreshTokenClaims(token string) (*TokenClaims, error) {
	return j.GetTokenClaims(token, j.refreshSecret)
}

func (j *jwtService) GetTokenClaims(token string, secret string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
