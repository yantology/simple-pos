package category

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// CategoryRepository implements the Repository interface
type CategoryRepository struct {
	postgres Repository
}

// NewCategoryRepository creates a new repository instance
func NewCategoryRepository(postgres Repository) Repository {
	return &CategoryRepository{postgres: postgres}
}

// GetAllCategoriesByUserID retrieves all categories for a specific user
func (r *CategoryRepository) GetAllCategoriesByUserID(userID string) ([]Category, *customerror.CustomError) {
	return r.postgres.GetAllCategoriesByUserID(userID)
}

// GetCategoryByID retrieves a category by its ID and user ID
func (r *CategoryRepository) GetCategoryByID(id string, userID string) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByID(id, userID)
}

// GetCategoryByName retrieves a category by its name and user ID
func (r *CategoryRepository) GetCategoryByName(name string, userID string) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByName(name, userID)
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(category *CreateCategory, userID string) (*Category, *customerror.CustomError) {
	return r.postgres.CreateCategory(category, userID)
}

// UpdateCategory updates an existing category, passing userID for authorization
func (r *CategoryRepository) UpdateCategory(id string, userID string, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) {
	return r.postgres.UpdateCategory(id, userID, category)
}

// DeleteCategory deletes a category by ID, passing userID for authorization
func (r *CategoryRepository) DeleteCategory(id string, userID string) *customerror.CustomError {
	return r.postgres.DeleteCategory(id, userID)
}
