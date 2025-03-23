package config

import (
	"os"
	"strconv"
	"time"
)

// accessTokenName := getEnv("ACCESS_TOKEN_COOKIE_NAME", "access_token")
// 	refreshTokenName := getEnv("REFRESH_TOKEN_COOKIE_NAME", "refresh_token")
// 	cookiePath := getEnv("COOKIE_PATH", "/")
// 	cookieDomain := getEnv("COOKIE_DOMAIN", "")
// 	secureCookie := getEnvBool("COOKIE_SECURE", true)

// 	// Access token duration (in seconds)
// 	accessTokenExpiry := 900 // 15 minutes
// 	// Refresh token duration (longer than access token)
// 	refreshTokenExpiry := 86400 // 24 hours

type TokenConfig struct {
	AccessTokenName    string
	RefreshTokenName   string
	CookiePath         string
	CookieDomain       string
	secureCookie       bool
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func InitTokenConfig() *TokenConfig {
	accessTokenName := os.Getenv("ACCESS_TOKEN_COOKIE_NAME")
	if accessTokenName == "" {
		accessTokenName = "access_token"
	}

	refreshTokenName := os.Getenv("REFRESH_TOKEN_COOKIE_NAME")
	if refreshTokenName == "" {
		refreshTokenName = "refresh_token"
	}

	cookiePath := os.Getenv("COOKIE_PATH")
	if cookiePath == "" {
		cookiePath = "/"
	}

	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	if cookieDomain == "" {
		cookieDomain = ""
	}

	secureCookie := os.Getenv("COOKIE_SECURE")
	if secureCookie == "" {
		secureCookie = "true"
	}
	secureCookieBool := true
	if secureCookie == "false" {
		secureCookieBool = false
	}

	accessTokenExpiry := os.Getenv("ACCESS_TOKEN_EXPIRY_minutes")
	accessTokenExpiryMinutes := 15 // Default: 15 minutes
	if accessTokenExpiry != "" {
		if minutes, err := strconv.Atoi(accessTokenExpiry); err == nil {
			accessTokenExpiryMinutes = minutes
		}
	}
	accessTokenExpiryDuration := time.Duration(accessTokenExpiryMinutes) * time.Minute

	refreshTokenExpiry := os.Getenv("REFRESH_TOKEN_EXPIRY_hours")
	refreshTokenExpiryHours := 24 // Default: 24 hours
	if refreshTokenExpiry != "" {
		if hours, err := strconv.Atoi(refreshTokenExpiry); err == nil {
			refreshTokenExpiryHours = hours
		}
	}
	refreshTokenExpiryDuration := time.Duration(refreshTokenExpiryHours) * time.Hour

	return &TokenConfig{
		AccessTokenName:    accessTokenName,
		RefreshTokenName:   refreshTokenName,
		CookiePath:         cookiePath,
		CookieDomain:       cookieDomain,
		secureCookie:       secureCookieBool,
		accessTokenExpiry:  accessTokenExpiryDuration,
		refreshTokenExpiry: refreshTokenExpiryDuration,
	}
}
