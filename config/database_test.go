package config_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yantology/golang-starter-template/config"
)

func TestInitDatabaseConfig(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected *config.DBConfig
	}{
		{
			name:    "with default values",
			envVars: map[string]string{},
			expected: &config.DBConfig{
				Host:     "127.0.0.1",
				Port:     "3306",
				Name:     "gin_gonic",
				User:     "root",
				Password: "",
				Driver:   "mysql",
			},
		},
		{
			name: "with custom values",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_NAME":     "testdb",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_DRIVER":   "postgres",
			},
			expected: &config.DBConfig{
				Host:     "localhost",
				Port:     "5432",
				Name:     "testdb",
				User:     "testuser",
				Password: "testpass",
				Driver:   "postgres",
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
			config := config.InitDatabaseConfig()

			// Assert results
			assert.Equal(t, tt.expected.Host, config.Host)
			assert.Equal(t, tt.expected.Port, config.Port)
			assert.Equal(t, tt.expected.Name, config.Name)
			assert.Equal(t, tt.expected.User, config.User)
			assert.Equal(t, tt.expected.Password, config.Password)
			assert.Equal(t, tt.expected.Driver, config.Driver)
		})
	}
}

func TestConnectDatabase(t *testing.T) {
	mockDB := &sql.DB{}
	mockOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}

	tests := []struct {
		name    string
		config  *config.DBConfig
		wantErr bool
	}{
		{
			name: "mysql connection",
			config: &config.DBConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				Name:     "testdb",
				User:     "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "postgres connection",
			config: &config.DBConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     "5432",
				Name:     "testdb",
				User:     "user",
				Password: "pass",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := config.ConnectDatabase(tt.config, mockOpen)
			assert.NotNil(t, db)
			assert.Equal(t, mockDB, db)
		})
	}
}
