package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yantology/simple-pos/config"
	jwtPkg "github.com/yantology/simple-pos/pkg/jwt"
)

// AuthMiddleware is a struct for authentication middleware
type AuthMiddleware struct {
	jwtService  jwtPkg.JWTService
	tokenConfig *config.TokenConfig
}

// NewAuthMiddleware creates a new instance of authentication middleware
func NewAuthMiddleware(jwtService jwtPkg.JWTService, tokenConfig *config.TokenConfig) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		tokenConfig: tokenConfig,
	}
}

// AuthRequired validates JWT token from cookies or authorization header
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		var err error

		// Try to get token from cookie first
		token, err = c.Cookie(m.tokenConfig.AccessTokenName)

		// If token not in cookie, try Authorization header
		if err != nil || token == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Tidak ada token autentikasi",
				})
				c.Abort()
				return
			}

			// Extract bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Format token tidak valid",
				})
				c.Abort()
				return
			}

			token = parts[1]
		}

		// Validate token
		claims, err := m.jwtService.ValidateAccessTokenClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid atau kadaluarsa",
			})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// ExtractUserClaims extracts user claims from the context
func ExtractUserClaims(c *gin.Context) *UserClaims {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")

	return &UserClaims{
		UserID: userID.(string),
		Email:  email.(string),
	}
}

// UserClaims represents user information extracted from token
type UserClaims struct {
	UserID string
	Email  string
}
