package category

import "github.com/yantology/simple-pos/pkg/customerror"

// Repository defines the data access methods for categories
type Repository interface {
	GetAllCategoriesByUserID(userID string) ([]Category, *customerror.CustomError) // Add method to get all categories by user ID
	// GetCategoryByID now requires userID for authorization
	GetCategoryByID(id string, userID string) (*Category, *customerror.CustomError)
	// GetCategoryByName now requires userID for authorization
	GetCategoryByName(name string, userID string) (*Category, *customerror.CustomError)
	// CreateCategory now requires userID
	CreateCategory(category *CreateCategory, userID string) (*Category, *customerror.CustomError) // Use CreateCategory from models.go, add userID
	// UpdateCategory now requires userID for authorization
	UpdateCategory(id string, userID string, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) // Use UpdateCategoryRequest from models.go
	// DeleteCategory now requires userID for authorization
	DeleteCategory(id string, userID string) *customerror.CustomError
}
