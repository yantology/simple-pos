package config

import (
	"os"
	"strconv"
	"time"
)

type TokenConfig struct {
	AccessTokenName    string
	RefreshTokenName   string
	CookiePath         string
	CookieDomain       string
	SecureCookie       bool
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
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
		SecureCookie:       secureCookieBool,
		AccessTokenExpiry:  accessTokenExpiryDuration,
		RefreshTokenExpiry: refreshTokenExpiryDuration,
	}
}
