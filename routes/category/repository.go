package category

import "github.com/yantology/simple-ecommerce/pkg/customerror"

// CategoryRepository implements the Repository interface
type CategoryRepository struct {
	postgres Repository
}

// NewCategoryRepository creates a new repository instance
func NewCategoryRepository(postgres Repository) *CategoryRepository {
	return &CategoryRepository{postgres: postgres}
}

// GetCategoryByID retrieves a category by its ID
func (r *CategoryRepository) GetCategoryByID(id int) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByID(id)
}

// GetCategoryByName retrieves a category by its name
func (r *CategoryRepository) GetCategoryByName(name string) (*Category, *customerror.CustomError) {
	return r.postgres.GetCategoryByName(name)
}

// GetCategoriesByUserID retrieves all categories created by a specific user
func (r *CategoryRepository) GetCategoriesByUserID(userID int) ([]Category, *customerror.CustomError) {
	return r.postgres.GetCategoriesByUserID(userID)
}

// CreateCategory creates a new category
func (r *CategoryRepository) CreateCategory(category *CreateCategoryRequest) (*Category, *customerror.CustomError) {
	return r.postgres.CreateCategory(category)
}

// UpdateCategory updates an existing category
func (r *CategoryRepository) UpdateCategory(id int, category *UpdateCategoryRequest) (*Category, *customerror.CustomError) {
	return r.postgres.UpdateCategory(id, category)
}

// DeleteCategory deletes a category by ID
func (r *CategoryRepository) DeleteCategory(id int) *customerror.CustomError {
	return r.postgres.DeleteCategory(id)
}
