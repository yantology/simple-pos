package product

import (
	"errors"
)

// Error definitions
var (
	ErrInvalidPrice = errors.New("price must be greater than zero")
)

// Repository defines the interface for product data operations
type Repository interface {
	Create(product Product) (Product, error)
	GetByID(id int) (Product, error)
	GetAll() ([]Product, error)
	Update(id int, product UpdateProductRequest) (Product, error)
	Delete(id int) error
	GetByUserID(userID int) ([]Product, error)
	GetByCategoryID(categoryID int) ([]Product, error)
}

// Service defines the interface for product business logic
// Service should only contain methods that don't require database access
type Service interface {
	ValidateProductInput(productRequest CreateProductRequest) error
	PrepareProductForCreation(productRequest CreateProductRequest) Product
	ValidateUpdateRequest(productRequest UpdateProductRequest) error
	FormatProductsResponse(products []Product) ProductListResponse
	FormatProductResponse(product Product) ProductResponse
}
