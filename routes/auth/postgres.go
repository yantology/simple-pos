package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/yantology/simple-pos/pkg/customerror"
)

// Verify interface implementation
var _ AuthDBInterface = (*authPostgres)(nil)

type authPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) AuthDBInterface {
	log.Println("[AuthPostgres] NewAuthPostgres: Initializing")
	return &authPostgres{db: db}
}

func (ap *authPostgres) CheckIsNotExistingEmail(email string) *customerror.CustomError {
	log.Printf("[AuthPostgres] CheckIsNotExistingEmail: Checking email: %s\n", email)
	// For empty email, return an error
	if email == "" {
		log.Println("[AuthPostgres] CheckIsNotExistingEmail: Email is empty")
		return customerror.NewCustomError(nil, "email is required", http.StatusBadRequest)
	}

	// Check if email exists with a simple count query (more consistent approach)
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	log.Printf("[AuthPostgres] CheckIsNotExistingEmail: Executing query: %s with email: %s\n", query, email)
	err := ap.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Printf("[AuthPostgres] CheckIsNotExistingEmail: Error querying email count for %s: %v\n", email, err)
		return customerror.NewPostgresError(err)
	}

	// If count > 0, email exists
	if count > 0 {
		log.Printf("[AuthPostgres] CheckIsNotExistingEmail: Email %s already exists (count: %d)\n", email, count)
		return customerror.NewCustomError(nil, "email already exists", http.StatusConflict)
	}

	log.Printf("[AuthPostgres] CheckIsNotExistingEmail: Email %s does not exist (count: %d)\n", email, count)
	// Email doesn't exist, so it's available
	return nil
}

func (ap *authPostgres) CheckIsExistingEmail(email string) *customerror.CustomError {
	log.Printf("[AuthPostgres] CheckIsExistingEmail: Checking email: %s\n", email)
	var exists int // Using int for Scan, though only 1 or sql.ErrNoRows is expected
	query := "SELECT 1 FROM users WHERE email = $1"
	log.Printf("[AuthPostgres] CheckIsExistingEmail: Executing query: %s with email: %s\n", query, email)
	err := ap.db.QueryRow(query, email).Scan(&exists)

	if err == nil {
		log.Printf("[AuthPostgres] CheckIsExistingEmail: Email %s exists\n", email)
		return nil // Email exists
	}
	if err == sql.ErrNoRows {
		log.Printf("[AuthPostgres] CheckIsExistingEmail: Email %s not found\n", email)
		return customerror.NewCustomError(nil, "email not found", http.StatusNotFound) // Email doesn't exist
	}

	log.Printf("[AuthPostgres] CheckIsExistingEmail: Error checking email %s: %v\n", email, err)
	return customerror.NewPostgresError(err) // Database error
}

func (ap *authPostgres) SaveActivationToken(req *ActivationTokenRequest) *customerror.CustomError {
	log.Printf("[AuthPostgres] SaveActivationToken: Saving token for email: %s, type: %s\n", req.Email, req.TokenType)
	query := `INSERT INTO activation_tokens (email, token_hash, type, expires_at) 
			  VALUES ($1, $2, $3, NOW() + ($4 || ' minutes')::interval)
			  ON CONFLICT (email, type) DO UPDATE
			  SET token_hash = EXCLUDED.token_hash,
				  expires_at = EXCLUDED.expires_at` // Corrected to use EXCLUDED for conflict update

	log.Printf("[AuthPostgres] SaveActivationToken: Executing query for email: %s\n", req.Email)
	_, err := ap.db.Exec(query, req.Email, req.ActivationCode, req.TokenType, req.ExpiryMinutes)
	if err != nil {
		log.Printf("[AuthPostgres] SaveActivationToken: Error saving activation token for %s: %v\n", req.Email, err)
		return customerror.NewPostgresError(err)
	}
	log.Printf("[AuthPostgres] SaveActivationToken: Successfully saved token for email: %s\n", req.Email)
	return nil
}

func (ap *authPostgres) GetActivationToken(req *GetActivationTokenRequest) (string, *customerror.CustomError) {
	log.Printf("[AuthPostgres] GetActivationToken: Getting token for email: %s, type: %s\n", req.Email, req.TokenType)
	var storedHash string
	query := `SELECT token_hash FROM activation_tokens 
			  WHERE email = $1 AND type = $2 AND expires_at > NOW()`

	log.Printf("[AuthPostgres] GetActivationToken: Executing query for email: %s, type: %s\n", req.Email, req.TokenType)
	err := ap.db.QueryRow(query, req.Email, req.TokenType).Scan(&storedHash)
	if err == sql.ErrNoRows {
		log.Printf("[AuthPostgres] GetActivationToken: Token not found or expired for email: %s, type: %s\n", req.Email, req.TokenType)
		return "", customerror.NewCustomError(err, "token not found or expired", http.StatusNotFound)
	}
	if err != nil {
		log.Printf("[AuthPostgres] GetActivationToken: Error getting token for email: %s, type: %s: %v\n", req.Email, req.TokenType, err)
		return "", customerror.NewPostgresError(err)
	}

	log.Printf("[AuthPostgres] GetActivationToken: Successfully retrieved token for email: %s, type: %s\n", req.Email, req.TokenType)
	return storedHash, nil
}

func (ap *authPostgres) CreateUser(req *CreateUserRequest) *customerror.CustomError {
	log.Printf("[AuthPostgres] CreateUser: Attempting to create user with email: %s\n", req.Email)
	tx, err := ap.db.Begin()
	if err != nil {
		log.Printf("[AuthPostgres] CreateUser: Error beginning transaction for %s: %v\n", req.Email, err)
		return customerror.NewPostgresError(err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[AuthPostgres] CreateUser: Recovered in CreateUser, rolling back transaction for %s: %v\n", req.Email, r)
			tx.Rollback()
		}
	}()

	// Insert new user
	log.Printf("[AuthPostgres] CreateUser: Inserting user %s into users table\n", req.Email)
	_, err = tx.Exec(`INSERT INTO users (email, fullname, password_hash) VALUES ($1, $2, $3)`,
		req.Email, req.Fullname, req.PasswordHash)
	if err != nil {
		log.Printf("[AuthPostgres] CreateUser: Error inserting user %s: %v\n", req.Email, err)
		tx.Rollback()
		return customerror.NewPostgresError(err)
	}

	// Delete activation token
	log.Printf("[AuthPostgres] CreateUser: Deleting activation token for user %s\n", req.Email)
	_, err = tx.Exec(`DELETE FROM activation_tokens WHERE email = $1 AND type = 'registration'`, req.Email) // Specify type for safety
	if err != nil {
		log.Printf("[AuthPostgres] CreateUser: Error deleting activation token for %s: %v\n", req.Email, err)
		tx.Rollback()
		return customerror.NewPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[AuthPostgres] CreateUser: Error committing transaction for %s: %v\n", req.Email, err)
		return customerror.NewPostgresError(err)
	}
	log.Printf("[AuthPostgres] CreateUser: Successfully created user %s and deleted activation token\n", req.Email)
	return nil
}

func (ap *authPostgres) GetUserByEmail(email string) (*User, *customerror.CustomError) {
	log.Printf("[AuthPostgres] GetUserByEmail: Attempting to get user by email: %s\n", email)
	user := &User{}
	query := `
		SELECT id, email, fullname, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`
	log.Printf("[AuthPostgres] GetUserByEmail: Executing query for email: %s\n", email)
	err := ap.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Fullname, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		log.Printf("[AuthPostgres] GetUserByEmail: User not found with email: %s\n", email)
		return nil, customerror.NewCustomError(err, "user not found", http.StatusNotFound)
	}
	if err != nil {
		log.Printf("[AuthPostgres] GetUserByEmail: Error retrieving user with email %s: %v\n", email, err)
		return nil, customerror.NewPostgresError(err)
	}
	log.Printf("[AuthPostgres] GetUserByEmail: Successfully retrieved user ID %d for email: %s\n", user.ID, email)
	return user, nil
}

func (ap *authPostgres) UpdateUserPassword(req *UpdatePasswordRequest) *customerror.CustomError {
	log.Printf("[AuthPostgres] UpdateUserPassword: Attempting to update password for email: %s\n", req.Email)
	tx, err := ap.db.Begin()
	if err != nil {
		log.Printf("[AuthPostgres] UpdateUserPassword: Error beginning transaction for %s: %v\n", req.Email, err)
		return customerror.NewPostgresError(err)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[AuthPostgres] UpdateUserPassword: Recovered in UpdateUserPassword, rolling back transaction for %s: %v\n", req.Email, r)
			tx.Rollback()
		}
	}()

	// Update password
	log.Printf("[AuthPostgres] UpdateUserPassword: Updating password_hash for user %s\n", req.Email)
	result, err := tx.Exec(`
		UPDATE users 
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE email = $2`,
		req.NewPasswordHash, req.Email)
	if err != nil {
		log.Printf("[AuthPostgres] UpdateUserPassword: Error updating password for %s: %v\n", req.Email, err)
		tx.Rollback()
		return customerror.NewPostgresError(err)
	}

	// Check if user exists
	rows, err := result.RowsAffected()
	if err != nil {
		log.Printf("[AuthPostgres] UpdateUserPassword: Error getting rows affected for %s: %v\n", req.Email, err)
		tx.Rollback()
		return customerror.NewPostgresError(err)
	}
	if rows == 0 {
		log.Printf("[AuthPostgres] UpdateUserPassword: User not found with email %s during password update\n", req.Email)
		tx.Rollback() // Rollback because user not found means no update happened
		return customerror.NewCustomError(nil, "user not found", http.StatusNotFound)
	}

	// Delete activation tokens (specifically for forget-password)
	log.Printf("[AuthPostgres] UpdateUserPassword: Deleting 'forget-password' activation tokens for user %s\n", req.Email)
	_, err = tx.Exec(`DELETE FROM activation_tokens WHERE email = $1 AND type = 'forget-password'`, req.Email)
	if err != nil {
		log.Printf("[AuthPostgres] UpdateUserPassword: Error deleting 'forget-password' activation tokens for %s: %v\n", req.Email, err)
		tx.Rollback()
		return customerror.NewPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[AuthPostgres] UpdateUserPassword: Error committing transaction for %s: %v\n", req.Email, err)
		return customerror.NewPostgresError(err)
	}
	log.Printf("[AuthPostgres] UpdateUserPassword: Successfully updated password for user %s and deleted 'forget-password' token\n", req.Email)
	return nil
}
