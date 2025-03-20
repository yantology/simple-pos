package auth

import "github.com/yantology/retail-pro-be/pkg/customerror"

type authRepository struct {
	db AuthDBInterface
}

func NewAuthRepository(db AuthDBInterface) *authRepository {
	return &authRepository{db: db}
}

func (ar *authRepository) CheckExistingEmail(email string) *customerror.CustomError {
	return ar.db.CheckExistingEmail(email)
}

func (ar *authRepository) SaveActivationToken(req *ActivationTokenRequest) *customerror.CustomError {
	return ar.db.SaveActivationToken(req)
}

func (ar *authRepository) ValidateActivationToken(req *ActivationTokenRequest) *customerror.CustomError {
	return ar.db.ValidateActivationToken(req)
}

func (ar *authRepository) CreateUser(req *CreateUserRequest) *customerror.CustomError {
	return ar.db.CreateUser(req)
}

func (ar *authRepository) GetUserByEmail(email string) (*User, *customerror.CustomError) {
	return ar.db.GetUserByEmail(email)
}

func (ar *authRepository) UpdateUserPassword(req *UpdatePasswordRequest) *customerror.CustomError {
	return ar.db.UpdateUserPassword(req)
}
