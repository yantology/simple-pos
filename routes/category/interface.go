package category

import "github.com/yantology/simple-pos/pkg/customerror"

// Repository defines the data access methods for categories
type Repository interface {
	GetAllCategoriesByUserID(userID int) ([]Category, *customerror.CustomError) // Changed userID to int
	// GetCategoryByID now requires userID for authorization
	GetCategoryByID(id int, userID int) (*Category, *customerror.CustomError) // Changed id and userID to int
	// GetCategoryByName now requires userID for authorization
	GetCategoryByName(name string, userID int) (*Category, *customerror.CustomError) // Changed userID to int
	// CreateCategory now requires userID
	CreateCategory(category *CreateCategory, userID int) (*Category, *customerror.CustomError) // Changed userID to int
	// UpdateCategory now requires userID for authorization
	UpdateCategory(id int, userID int, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) // Changed id and userID to int
	// DeleteCategory now requires userID for authorization
	DeleteCategory(id int, userID int) *customerror.CustomError // Changed id and userID to int
}
