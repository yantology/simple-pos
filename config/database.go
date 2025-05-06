package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
)

// DBConfig holds database configuration
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Driver   string
}

// InitDatabaseConfig initializes database configuration from environment variables or mounted secrets
func InitDatabaseConfig() *DBConfig {
	log.Println("Initializing database configuration...") // Added log
	dbConfig := &DBConfig{
		Host:     "127.0.0.1",
		Port:     "3306",
		Name:     "gin_gonic",
		User:     "root",
		Password: "",
		Driver:   "mysql",
	}

	if env := os.Getenv("DB_HOST"); env != "" {
		log.Println("DB_HOST => ", env)
		dbConfig.Host = env
	}

	if env := os.Getenv("DB_PORT"); env != "" {
		log.Println("DB_PORT => ", env)
		dbConfig.Port = env
	}

	if env := os.Getenv("DB_NAME"); env != "" {
		log.Println("DB_NAME => ", env)
		dbConfig.Name = env
	}

	if env := os.Getenv("DB_USER"); env != "" {
		log.Println("DB_USER => ", env)
		dbConfig.User = env
	}

	// Read DB_PASSWORD: prioritize env var, fallback to secret file
	dbPasswordEnv := os.Getenv("DB_PASSWORD")
	if dbPasswordEnv != "" {
		log.Println("DB_PASSWORD => Read from environment variable [MASKED]")
		dbConfig.Password = dbPasswordEnv
	} else {
		// Try reading from the mounted secret file (Cloud Run)
		secretPath := "/secrets/db_password/db_password"
		if _, err := os.Stat(secretPath); err == nil {
			log.Printf("DB_PASSWORD => Attempting to read from secret file: %s\n", secretPath)
			passwordBytes, err := ioutil.ReadFile(secretPath)
			if err != nil {
				log.Fatalf("Failed to read database password secret file %s: %v\n", secretPath, err)
			}
			dbConfig.Password = strings.TrimSpace(string(passwordBytes))
			log.Println("DB_PASSWORD => Successfully read from secret file [MASKED]")
		} else {
			log.Println("DB_PASSWORD => Not found in environment variable or secret file. Using default (if any) or empty.")
			// Keep the default empty password if neither env var nor file exists
		}
	}

	if env := os.Getenv("DB_DRIVER"); env != "" {
		log.Println("DB_DRIVER => ", env)
		dbConfig.Driver = env
	}

	return dbConfig
}

// Define a type for the Open function
type OpenFunc func(driverName, dataSourceName string) (*sql.DB, error)

// ConnectDatabase establishes a connection to the database
func ConnectDatabase(dbConfig *DBConfig, openFunc OpenFunc) *sql.DB {
	var db *sql.DB
	var errConnection error
	var dsn string // Declare dsn variable

	if dbConfig.Driver == "mysql" {
		// Handle MySQL connection (assuming it needs standard host/port)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		log.Printf("Connecting to MySQL with DSN: %s:%s@tcp(%s:%s)/%s\\n", dbConfig.User, "[MASKED]", dbConfig.Host, dbConfig.Port, dbConfig.Name) // Log DSN without password

	} else if dbConfig.Driver == "postgres" {
		// Check if DB_HOST is a Unix socket path for Cloud SQL
		if strings.HasPrefix(dbConfig.Host, "/") { // Check if Host starts with '/' (likely a socket path)
			// Use Unix socket DSN format
			dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
				dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name)
			log.Printf("Connecting to PostgreSQL (Cloud SQL Socket) with DSN: host=%s user=%s password=%s dbname=%s sslmode=disable\\n", dbConfig.Host, dbConfig.User, "[MASKED]", dbConfig.Name) // Log DSN without password
		} else {
			// Use standard host/port DSN format
			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", // Removed TimeZone for simplicity, add if needed
				dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)
			log.Printf("Connecting to PostgreSQL (Host/Port) with DSN: host=%s port=%s user=%s password=%s dbname=%s sslmode=disable\\n", dbConfig.Host, dbConfig.Port, dbConfig.User, "[MASKED]", dbConfig.Name) // Log DSN without password
		}
	} else {
		log.Fatalf("invalid database driver: %s\\n", dbConfig.Driver) // Use log.Fatalf for fatal errors
	}

	// Attempt connection
	db, errConnection = openFunc(dbConfig.Driver, dsn)
	if errConnection != nil {
		log.Fatalf("failed to connect to database (%s): %v\\nDSN used: %s\\n", dbConfig.Driver, errConnection, maskPassword(dsn)) // Log error with masked DSN
	}

	log.Printf("Successfully connected to %s database.\\n", dbConfig.Driver)
	return db
}

// Helper function to mask password in DSN for logging
func maskPassword(dsn string) string {
	// Basic masking for key=value format
	re := regexp.MustCompile(`(password|passwd|pwd)=([^ ]+)`)
	maskedDSN := re.ReplaceAllString(dsn, "$1=[MASKED]")

	// Masking for postgres://user:password@host format
	rePostgres := regexp.MustCompile(`:(//[^:]+:)([^@]+)@`)
	maskedDSN = rePostgres.ReplaceAllString(maskedDSN, ":$1[MASKED]@")

	return maskedDSN
}
