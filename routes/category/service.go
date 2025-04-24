package category

import (
	"net/http"

	"github.com/yantology/simple-ecommerce/pkg/customerror"
)

// CategoryService implements the Service interface
type CategoryService struct {
	// No repository dependency
}

// NewCategoryService creates a new service instance
func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// ValidateCategoryID validates a category ID
func (s *CategoryService) ValidateCategoryID(id int) *customerror.CustomError {
	if id <= 0 {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Invalid category ID", http.StatusBadRequest)
	}
	return nil
}

// ValidateCategoryName validates a category name
func (s *CategoryService) ValidateCategoryName(name string) *customerror.CustomError {
	if name == "" {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Category name cannot be empty", http.StatusBadRequest)
	}
	return nil
}

// ValidateUserID validates a user ID
func (s *CategoryService) ValidateUserID(userID int) *customerror.CustomError {
	if userID <= 0 {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Invalid user ID", http.StatusBadRequest)
	}
	return nil
}

// ValidateCreateCategoryRequest validates the create category request
func (s *CategoryService) ValidateCreateCategoryRequest(category *CreateCategoryRequest) *customerror.CustomError {
	if category == nil {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Category request cannot be nil", http.StatusBadRequest)
	}

	if category.Name == "" {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Category name cannot be empty", http.StatusBadRequest)
	}

	if category.UserID <= 0 {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Invalid user ID", http.StatusBadRequest)
	}

	return nil
}

// ValidateUpdateCategoryRequest validates the update category request
func (s *CategoryService) ValidateUpdateCategoryRequest(id int, category *UpdateCategoryRequest) *customerror.CustomError {
	if id <= 0 {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Invalid category ID", http.StatusBadRequest)
	}

	if category == nil {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Category request cannot be nil", http.StatusBadRequest)
	}

	if category.Name == "" {
		// Use NewCustomError with http.StatusBadRequest
		return customerror.NewCustomError(nil, "Category name cannot be empty", http.StatusBadRequest)
	}

	return nil
}

// CheckOwnership verifies that the user making the request owns the category
func (s *CategoryService) CheckOwnership(categoryUserID, requestUserID int) *customerror.CustomError {
	if categoryUserID != requestUserID {
		// Use NewCustomError with http.StatusForbidden
		return customerror.NewCustomError(nil, "You don't have permission to modify this category", http.StatusForbidden)
	}
	return nil
}
