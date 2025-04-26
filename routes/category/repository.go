package category

import "github.com/yantology/simple-pos/pkg/customerror"

// CategoryRepository implements the Repository interface
type CategoryRepository struct {
	postgres Repository
}

// NewCategoryRepository creates a new repository instance
func NewCategoryRepository(postgres Repository) Repository {
	return &CategoryRepository{postgres: postgres}
}

// GetAllCategoriesByUserID retrieves all categories for a specific user
func (r *CategoryRepository) GetAllCategoriesByUserID(userID int) ([]Category, *customerror.CustomError) {
	return r.postgres.GetAllCategoriesByUserID(userID)
}

// GetCategoryByID retrieves a category by its ID and user ID
func (r *CategoryRepository) GetCategoryByID(id int, userID int) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByID(id, userID)
}

// GetCategoryByName retrieves a category by its name and user ID
func (r *CategoryRepository) GetCategoryByName(name string, userID int) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByName(name, userID)
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(category *CreateCategory, userID int) (*Category, *customerror.CustomError) {
	return r.postgres.CreateCategory(category, userID)
}

// UpdateCategory updates an existing category, passing userID for authorization
func (r *CategoryRepository) UpdateCategory(id int, userID int, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) {
	return r.postgres.UpdateCategory(id, userID, category)
}

// DeleteCategory deletes a category by ID, passing userID for authorization
func (r *CategoryRepository) DeleteCategory(id int, userID int) *customerror.CustomError {
	return r.postgres.DeleteCategory(id, userID)
}
