package category

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// Repository defines the data access methods for categories
type Repository interface {
	GetCategoryByID(id int) (*Category, *customerror.CustomError)
	GetCategoryByName(name string) (*Category, *customerror.CustomError)
	GetCategoriesByUserID(userID int) ([]Category, *customerror.CustomError)
	CreateCategory(category *CreateCategoryRequest) (*Category, *customerror.CustomError)
	UpdateCategory(id int, category *UpdateCategoryRequest) (*Category, *customerror.CustomError)
	DeleteCategory(id int) *customerror.CustomError
}

// Service defines the business logic operations for categories that don't require database access
type Service interface {
	ValidateCategoryID(id int) *customerror.CustomError
	ValidateCategoryName(name string) *customerror.CustomError
	ValidateUserID(userID int) *customerror.CustomError
	ValidateCreateCategoryRequest(category *CreateCategoryRequest) *customerror.CustomError
	ValidateUpdateCategoryRequest(id int, category *UpdateCategoryRequest) *customerror.CustomError
	CheckOwnership(categoryUserID, requestUserID int) *customerror.CustomError
}
