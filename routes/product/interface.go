package product

import "github.com/yantology/simple-pos/pkg/customerror"

// Repository defines the interface for product data operations
type Repository interface {
	Create(productData *CreateProduct, userID int) (*Product, *customerror.CustomError) // Changed userID to int
	GetAll() ([]*Product, *customerror.CustomError)
	Update(id int, userID int, product *UpdateProduct) (*Product, *customerror.CustomError) // Changed id and userID to int
	Delete(id int, userID int) *customerror.CustomError                                     // Changed id and userID to int
	GetByCategoryID(categoryID int) ([]*Product, *customerror.CustomError)
}
