package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/simple-pos/config"
)

func TestInitJWTConfig(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *config.JWTConfig
		shouldError    bool
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"JWT_ACCESS_SECRET":           "test-access-secret",
				"JWT_REFRESH_SECRET":          "test-refresh-secret",
				"JWT_ACCESS_DURATION_MINUTES": "30",
				"JWT_REFRESH_DURATION_DAYS":   "14",
				"JWT_ISSUER":                  "test-issuer",
			},
			expectedConfig: &config.JWTConfig{
				AccessSecret:    "test-access-secret",
				RefreshSecret:   "test-refresh-secret",
				AccessDuration:  30 * time.Minute,
				RefreshDuration: 14 * 24 * time.Hour,
				Issuer:          "test-issuer",
			},
			shouldError: false,
		},
		{
			name: "missing access secret",
			envVars: map[string]string{
				"JWT_REFRESH_SECRET": "test-refresh-secret",
			},
			expectedConfig: nil,
			shouldError:    true,
		},
		{
			name: "missing refresh secret",
			envVars: map[string]string{
				"JWT_ACCESS_SECRET": "test-access-secret",
			},
			expectedConfig: nil,
			shouldError:    true,
		},
		{
			name: "default durations",
			envVars: map[string]string{
				"JWT_ACCESS_SECRET":  "test-access-secret",
				"JWT_REFRESH_SECRET": "test-refresh-secret",
			},
			expectedConfig: &config.JWTConfig{
				AccessSecret:    "test-access-secret",
				RefreshSecret:   "test-refresh-secret",
				AccessDuration:  15 * time.Minute,   // Default access duration
				RefreshDuration: 7 * 24 * time.Hour, // Default refresh duration
				Issuer:          "retail-pro",       // Default issuer
			},
			shouldError: false,
		},
		{
			name: "invalid duration values",
			envVars: map[string]string{
				"JWT_ACCESS_SECRET":           "test-access-secret",
				"JWT_REFRESH_SECRET":          "test-refresh-secret",
				"JWT_ACCESS_DURATION_MINUTES": "invalid",
				"JWT_REFRESH_DURATION_DAYS":   "invalid",
			},
			expectedConfig: &config.JWTConfig{
				AccessSecret:    "test-access-secret",
				RefreshSecret:   "test-refresh-secret",
				AccessDuration:  15 * time.Minute,   // Default access duration
				RefreshDuration: 7 * 24 * time.Hour, // Default refresh duration
				Issuer:          "retail-pro",       // Default issuer
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables for test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Run test
			config, err := config.InitJWTConfig()

			if tt.shouldError {
				assert.NotNil(t, err, "Expected an error but got none")
				assert.Nil(t, config, "Expected nil config when error occurs")
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, config, "Expected non-nil config")
				assert.Equal(t, tt.expectedConfig.AccessSecret, config.AccessSecret)
				assert.Equal(t, tt.expectedConfig.RefreshSecret, config.RefreshSecret)
				assert.Equal(t, tt.expectedConfig.AccessDuration, config.AccessDuration)
				assert.Equal(t, tt.expectedConfig.RefreshDuration, config.RefreshDuration)
				assert.Equal(t, tt.expectedConfig.Issuer, config.Issuer)
			}
		})
	}
}
