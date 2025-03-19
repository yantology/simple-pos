package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

// InitDatabaseConfig initializes database configuration from environment variables
func InitDatabaseConfig() *DBConfig {
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

	if env := os.Getenv("DB_PASSWORD"); env != "" {
		log.Println("DB_PASSWORD => ", "[MASKED]") // For security, don't log the actual password
		dbConfig.Password = env
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

	if dbConfig.Driver == "mysql" {
		dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		db, errConnection = openFunc(dbConfig.Driver, dsnMysql)
	} else if dbConfig.Driver == "postgres" {
		dsnPgSql := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
		db, errConnection = openFunc(dbConfig.Driver, dsnPgSql)
	} else {
		fmt.Printf("invalid database driver: %s\n", dbConfig.Driver)
		os.Exit(1)
	}

	if errConnection != nil {
		fmt.Printf("failed to connect to database: %v", errConnection)
		os.Exit(1)
	}

	return db
}
