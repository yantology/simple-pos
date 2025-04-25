package product

import "github.com/yantology/simple-pos/pkg/customerror"

// Repository defines the interface for product data operations
type Repository interface {
	Create(product *CreateProduct, userID string) (*Product, *customerror.CustomError) // Use CreateProduct from models.go
	GetAll() ([]*Product, *customerror.CustomError)
	Update(id, userID string, product *UpdateProduct) (*Product, *customerror.CustomError) // Use UpdateProduct from models.go
	Delete(id, userID string) *customerror.CustomError                                     // Add userID
	GetByUserID(userID string) ([]*Product, *customerror.CustomError)                      // userID is already string
	GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError)
}
