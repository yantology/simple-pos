package config

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

type JWTConfig struct {
	AccessSecret    string
	RefreshSecret   string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
	Issuer          string
}

func InitJWTConfig() (*JWTConfig, *customerror.CustomError) {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		log.Println("JWT access secret is not set")
		return nil, customerror.NewCustomError(nil, "JWT access secret is not set", http.StatusUnauthorized)
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		log.Println("JWT refresh secret is not set")
		return nil, customerror.NewCustomError(nil, "JWT refresh secret is not set", http.StatusUnauthorized)
	}

	// Get access duration from env or use default (15 minutes)
	accessDurationMinutes := 15
	if envDuration := os.Getenv("JWT_ACCESS_DURATION_MINUTES"); envDuration != "" {
		if parsed, err := strconv.Atoi(envDuration); err == nil {
			accessDurationMinutes = parsed
		}
	}

	// Get refresh duration from env or use default (7 days)
	refreshDurationDays := 7
	if envDuration := os.Getenv("JWT_REFRESH_DURATION_DAYS"); envDuration != "" {
		if parsed, err := strconv.Atoi(envDuration); err == nil {
			refreshDurationDays = parsed
		}
	}

	// Get issuer from env or use default
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "retail-pro"
	}

	return &JWTConfig{
		AccessSecret:    accessSecret,
		RefreshSecret:   refreshSecret,
		AccessDuration:  time.Duration(accessDurationMinutes) * time.Minute,
		RefreshDuration: time.Duration(refreshDurationDays) * 24 * time.Hour,
		Issuer:          issuer,
	}, nil
}
