package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/golang-starter-template/config"
)

func TestCorsConfig(t *testing.T) {
	tests := []struct {
		name       string
		envOrigins string
		expected   []string
	}{
		{
			name:       "default origins",
			envOrigins: "",
			expected:   []string{"*"},
		},
		{
			name:       "single custom origin",
			envOrigins: "http://localhost:3000",
			expected:   []string{"http://localhost:3000"},
		},
		{
			name:       "multiple custom origins",
			envOrigins: "http://localhost:3000,http://example.com",
			expected:   []string{"http://localhost:3000", "http://example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables if needed
			if tt.envOrigins != "" {
				os.Setenv("CORS_ALLOW_ORIGINS", tt.envOrigins)
			}

			// Test the handler creation
			handler := config.CorsConfig()
			assert.NotNil(t, handler, "CORS handler should not be nil")
		})
	}
}

func TestGetEnvAsSlice(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "empty environment variable",
			envKey:       "TEST_SLICE",
			envValue:     "",
			defaultValue: []string{"default"},
			expected:     []string{"default"},
		},
		{
			name:         "single value",
			envKey:       "TEST_SLICE",
			envValue:     "value1",
			defaultValue: []string{"default"},
			expected:     []string{"value1"},
		},
		{
			name:         "multiple values",
			envKey:       "TEST_SLICE",
			envValue:     "value1,value2,value3",
			defaultValue: []string{"default"},
			expected:     []string{"value1", "value2", "value3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variable if needed
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
			}

			// Test the helper function
			result := config.GetEnvAsSlice(tt.envKey, tt.defaultValue)
			assert.Equal(t, tt.expected, result, "Slice values should match expected")
		})
	}
}
