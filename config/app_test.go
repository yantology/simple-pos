package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/retail-pro-be/config"
)

func TestInitAppConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *config.AppConfig
	}{
		{
			name:    "with default values",
			envVars: map[string]string{},
			expected: &config.AppConfig{
				Port:            "3000",
				PublicRoute:     "/public",
				PublicAssetsDir: "./public",
			},
		},
		{
			name: "with custom port",
			envVars: map[string]string{
				"APP_PORT": "8080",
			},
			expected: &config.AppConfig{
				Port:            "8080",
				PublicRoute:     "/public",
				PublicAssetsDir: "./public",
			},
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
			config := config.InitAppConfig()

			// Assert results
			assert.Equal(t, tt.expected.Port, config.Port)
			assert.Equal(t, tt.expected.PublicRoute, config.PublicRoute)
			assert.Equal(t, tt.expected.PublicAssetsDir, config.PublicAssetsDir)
		})
	}
}
