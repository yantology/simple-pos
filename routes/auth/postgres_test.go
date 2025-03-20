package auth_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/yantology/retail-pro-be/routes/auth"
)

func TestCheckExistingEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	tests := []struct {
		name      string
		email     string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:  "email exists",
			email: "test@example.com",
			mockSetup: func() {
				mock.ExpectQuery("SELECT 1 FROM users WHERE email = \\$1").
					WithArgs("test@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name:  "email not exists",
			email: "notfound@example.com",
			mockSetup: func() {
				mock.ExpectQuery("SELECT 1 FROM users WHERE email = \\$1").
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ap.CheckExistingEmail(tt.email)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestSaveActivationToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	tests := []struct {
		name      string
		req       *auth.ActivationTokenRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful save",
			req: &auth.ActivationTokenRequest{
				Email:         "test@example.com",
				TokenHash:     "hashed_token_123",
				TokenType:     "activation",
				ExpiryMinutes: 30,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO activation_tokens").
					WithArgs("test@example.com", "hashed_token_123", "activation", 30).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "db error",
			req: &auth.ActivationTokenRequest{
				Email:         "test@example.com",
				TokenHash:     "hashed_token_123",
				TokenType:     "activation",
				ExpiryMinutes: 30,
			},
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO activation_tokens").
					WithArgs("test@example.com", "hashed_token_123", "activation", 30).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ap.SaveActivationToken(tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestValidateActivationToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	tests := []struct {
		name      string
		req       *auth.ActivationTokenRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "valid token",
			req: &auth.ActivationTokenRequest{
				Email:     "test@example.com",
				TokenHash: "hashed_token_123",
				TokenType: "activation",
			},
			mockSetup: func() {
				mock.ExpectQuery("SELECT id FROM activation_tokens").
					WithArgs("test@example.com", "hashed_token_123", "activation").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			wantErr: false,
		},
		{
			name: "invalid token",
			req: &auth.ActivationTokenRequest{
				Email:     "test@example.com",
				TokenHash: "invalid_hash",
				TokenType: "activation",
			},
			mockSetup: func() {
				mock.ExpectQuery("SELECT id FROM activation_tokens").
					WithArgs("test@example.com", "invalid_hash", "activation").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ap.ValidateActivationToken(tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCreateauth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	tests := []struct {
		name      string
		req       *auth.CreateUserRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful create",
			req: &auth.CreateUserRequest{
				Email:        "test@example.com",
				Fullname:     "Test auth.User",
				PasswordHash: "hashedpassword",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", "Test auth.User", "hashedpassword").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("DELETE FROM activation_tokens").
					WithArgs("test@example.com").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "transaction error",
			req: &auth.CreateUserRequest{
				Email:        "test@example.com",
				Fullname:     "Test auth.User",
				PasswordHash: "hashedpassword",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test@example.com", "Test auth.User", "hashedpassword").
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ap.CreateUser(tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	now := time.Now()
	tests := []struct {
		name      string
		email     string
		mockSetup func()
		want      *auth.User
		wantErr   bool
	}{
		{
			name:  "user exists",
			email: "test@example.com",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "fullname", "password_hash", "created_at", "updated_at"}).
					AddRow("1", "test@example.com", "Test auth.User", "hashedpassword", now, now)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &auth.User{
				ID:           "1",
				Email:        "test@example.com",
				Fullname:     "Test auth.User",
				PasswordHash: "hashedpassword",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "notfound@example.com",
			mockSetup: func() {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE").
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := ap.GetUserByEmail(tt.email)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()

	ap := auth.NewAuthPostgres(db)

	tests := []struct {
		name      string
		req       *auth.UpdatePasswordRequest
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "successful update",
			req: &auth.UpdatePasswordRequest{
				Email:           "test@example.com",
				NewPasswordHash: "newhashedpassword",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users").
					WithArgs("newhashedpassword", "test@example.com").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("DELETE FROM activation_tokens").
					WithArgs("test@example.com").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &auth.UpdatePasswordRequest{
				Email:           "notfound@example.com",
				NewPasswordHash: "newhashedpassword",
			},
			mockSetup: func() {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE users").
					WithArgs("newhashedpassword", "notfound@example.com").
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := ap.UpdateUserPassword(tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
