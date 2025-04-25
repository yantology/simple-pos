package auth

import (
	"database/sql"
	"fmt"
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
	return &authPostgres{db: db}
}

func (ap *authPostgres) CheckIsNotExistingEmail(email string) *customerror.CustomError {
	// For empty email, return an error
	if email == "" {
		return customerror.NewCustomError(nil, "email is required", http.StatusBadRequest)
	}

	// Check if email exists with a simple count query (more consistent approach)
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1"
	err := ap.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	// If count > 0, email exists
	if count > 0 {
		return customerror.NewCustomError(nil, "email already exists", http.StatusConflict)
	}

	// Email doesn't exist, so it's available
	return nil
}

func (ap *authPostgres) CheckIsExistingEmail(email string) *customerror.CustomError {
	var exists int
	err := ap.db.QueryRow("SELECT 1 FROM users WHERE email = $1", email).Scan(&exists)

	if err == nil {
		return nil // Email exists
	}
	if err == sql.ErrNoRows {
		return customerror.NewCustomError(nil, "email not found", http.StatusNotFound) // Email doesn't exist
	}

	return customerror.NewPostgresError(err) // Database error
}

func (ap *authPostgres) SaveActivationToken(req *ActivationTokenRequest) *customerror.CustomError {
	query := `INSERT INTO activation_tokens (email, token_hash, type, expires_at) 
			  VALUES ($1, $2, $3, NOW() + ($4 || ' minutes')::interval)
			  ON CONFLICT (email, type) DO UPDATE
			  SET token_hash = $2,
				  expires_at = NOW() + ($4 || ' minutes')::interval`

	_, err := ap.db.Exec(query, req.Email, req.ActivationCode, req.TokenType, req.ExpiryMinutes)
	if err != nil {
		log.Println("Error saving activation token:", err)
		return customerror.NewPostgresError(err)
	}
	return nil
}

func (ap *authPostgres) GetActivationToken(req *GetActivationTokenRequest) (string, *customerror.CustomError) {
	var storedHash string
	query := `SELECT token_hash FROM activation_tokens 
			  WHERE email = $1 AND type = $2 AND expires_at > NOW()`

	err := ap.db.QueryRow(query, req.Email, req.TokenType).Scan(&storedHash)
	if err == sql.ErrNoRows {
		return "", customerror.NewCustomError(err, "token not found or expired", http.StatusNotFound)
	}
	if err != nil {
		return "", customerror.NewPostgresError(err)
	}

	return storedHash, nil
}

func (ap *authPostgres) CreateUser(req *CreateUserRequest) *customerror.CustomError {
	tx, err := ap.db.Begin()
	if err != nil {
		return customerror.NewPostgresError(err)
	}
	defer tx.Rollback()

	// Insert new user
	_, err = tx.Exec(`INSERT INTO users (email, fullname, password_hash) VALUES ($1, $2, $3)`,
		req.Email, req.Fullname, req.PasswordHash)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	// Delete activation token
	_, err = tx.Exec(`DELETE FROM activation_tokens WHERE email = $1`, req.Email)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		return customerror.NewPostgresError(err)
	}
	return nil
}

func (ap *authPostgres) GetUserByEmail(email string) (*User, *customerror.CustomError) {
	user := &User{}
	err := ap.db.QueryRow(`
		SELECT id, email, fullname, password_hash, created_at, updated_at 
		FROM users WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.Fullname, &user.PasswordHash,
		&user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		fmt.Println("User not found:", email)
		return nil, customerror.NewCustomError(err, "user not found", http.StatusNotFound)
	}
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return nil, customerror.NewPostgresError(err)
	}
	return user, nil
}

func (ap *authPostgres) UpdateUserPassword(req *UpdatePasswordRequest) *customerror.CustomError {
	tx, err := ap.db.Begin()
	if err != nil {
		return customerror.NewPostgresError(err)
	}
	defer tx.Rollback()

	// Update password
	result, err := tx.Exec(`
		UPDATE users 
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE email = $2`,
		req.NewPasswordHash, req.Email)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	// Check if user exists
	rows, err := result.RowsAffected()
	if err != nil {
		return customerror.NewPostgresError(err)
	}
	if rows == 0 {
		return customerror.NewCustomError(nil, "user not found", http.StatusNotFound)
	}

	// Delete activation tokens
	_, err = tx.Exec(`DELETE FROM activation_tokens WHERE email = $1`, req.Email)
	if err != nil {
		return customerror.NewPostgresError(err)
	}

	if err = tx.Commit(); err != nil {
		return customerror.NewPostgresError(err)
	}
	return nil
}
