package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/simple-pos/config"
)

func TestInitResendConfig(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		wantError bool
	}{
		{
			name: "with valid configuration",
			envVars: map[string]string{
				"RESEND_API_KEY": "test-api-key",
				"RESEND_DOMAIN":  "test.com",
				"RESEND_NAME":    "Test Sender",
			},
			wantError: false,
		},
		{
			name:      "with missing API key",
			envVars:   map[string]string{},
			wantError: true,
		},
		{
			name: "with missing domain",
			envVars: map[string]string{
				"RESEND_API_KEY": "test-api-key",
			},
			wantError: true,
		},
		{
			name: "with missing name",
			envVars: map[string]string{
				"RESEND_API_KEY": "test-api-key",
				"RESEND_DOMAIN":  "test.com",
			},
			wantError: true,
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
			config, err := config.InitResendConfig()

			if tt.wantError {
				assert.NotNil(t, err, "Expected an error but got none")
				assert.Nil(t, config, "Expected nil config when error occurs")
			} else {
				assert.Nil(t, err, "Unexpected error")
				assert.NotNil(t, config, "Expected non-nil config")
			}
		})
	}
}
